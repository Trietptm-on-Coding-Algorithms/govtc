package govtc

import (
	"bufio"
	"database/sql"
	"fmt"
	"os"
	"strings"
	_ "github.com/mattn/go-sqlite3"
)

type BatchCheck struct {
	Hashes []string
}

type CacheChecker struct {
	apiKey 			string
	databasePath 	string
	checkLimit		int
	workQueue		chan BatchCheck
	workerQueue 	chan chan BatchCheck
}

const SQL_SELECT_MD5 string = "SELECT id, md5, sha256, positives, total, permalink, responsecode, scans, scandate, updatedate from vthashes where md5 = ?"
const SQL_SELECT_SHA256 string = "SELECT id, md5, sha256, positives, total, permalink, responsecode, scans, scandate, updatedate from vthashes where sha256 = ?"

func NewCacheChecker(apiKey string, databasePath string) *CacheChecker {
	c := new(CacheChecker)
	c.apiKey = apiKey
	c.databasePath = databasePath
	c.checkLimit = 4
	return c
}

func (c CacheChecker) ProcessFile(hashChannel chan *VtRecord, inputFile string, mode int) {
	c.workQueue = make(chan BatchCheck, 1000)

	c.startDispatcher(2, c.databasePath, c.apiKey)

	// Open the database containing our VT data
	db, err := sql.Open("sqlite3", c.databasePath)
	defer db.Close()
	if err != nil {
		// Error? callback?
		return
	}

	// Prepare the SQL statements
	stmtMd5, err := db.Prepare(SQL_SELECT_MD5)
	if err != nil {
		// Error?
	}
	defer stmtMd5.Close()

	stmtSha256, err := db.Prepare(SQL_SELECT_SHA256)
	if err != nil {
		// Error?
	}
	defer stmtSha256.Close()

	file, _ := os.Open(inputFile)
	defer file.Close()
	scanner := bufio.NewScanner(file)
	line := ""
	tempHashes := make([]string, c.checkLimit)
	for scanner.Scan() {
		line = strings.TrimSpace(scanner.Text())
		if len(line) == 0 {
			continue
		}

		vtRecord, err := isHashInDatabase(stmtMd5, stmtSha256, line)
		if err != nil {
			fmt.Println(err.Error())
			continue
		}

		if mode == MODE_CACHE || mode == MODE_LIVE {
			tempHashes = append(tempHashes, line)
			if len(tempHashes) >= c.checkLimit {
				batchCheck := BatchCheck{}
				//batchCheck := new(BatchCheck)
				batchCheck.Add(tempHashes)
				tempHashes = make([]string, c.checkLimit)

				c.workQueue <- batchCheck
			}
		} else {
			// Hash identified? Callback
			hashChannel <-vtRecord
		}
	}
}

func (c CacheChecker)startDispatcher(workerCount int, databasePath string, apiKey string) {
	// First, initialize the channel we are going to but the workers' work channels into.
	c.workerQueue = make(chan chan BatchCheck, workerCount)

	// Now, create all of our workers.
	for i := 0; i < workerCount; i++ {
		fmt.Println("Starting worker", i+1)
		worker := NewVtChecker(i+1, databasePath, apiKey, c.workerQueue)
		worker.Start()
	}

	go func() {
		for {
			//fmt.Println("LOOPING")
			select {
			case work := <-c.workQueue:
				//fmt.Println("Received work requeust")
				go func() {
					worker := <-c.workerQueue

					//fmt.Println("Dispatching work request")
					worker <- work
				}()
			}
		}
	}()
}

//
func isHashInDatabase(stmtMd5 *sql.Stmt, stmtSha265 *sql.Stmt, hash string) (*VtRecord, error) {
	vtRecord := new(VtRecord)
	if len(hash) == 32 {
		err := stmtMd5.QueryRow(hash).Scan(&vtRecord.Id,
										   &vtRecord.Md5,
			                               &vtRecord.Sha256,
										   &vtRecord.Positives,
			                               &vtRecord.Total,
										   &vtRecord.PermaLink,
			                               &vtRecord.ResponseCode,
										   &vtRecord.Scans,
										   &vtRecord.ScanDate,
										   &vtRecord.UpdateDate)
		if err != nil {
			return vtRecord, err
		}
	} else {
		err := stmtSha265.QueryRow(hash).Scan(&vtRecord.Id,
											  &vtRecord.Md5,
											  &vtRecord.Sha256,
											  &vtRecord.Positives,
											  &vtRecord.Total,
											  &vtRecord.PermaLink,
											  &vtRecord.ResponseCode,
											  &vtRecord.Scans,
											  &vtRecord.ScanDate,
											  &vtRecord.UpdateDate)
		if err != nil {
			return vtRecord, err
		}
	}

	return vtRecord, nil
}

func (c CacheChecker) ProcessHash(hash string, mode int) {

}

//
func (b BatchCheck) Add(hashes []string) {
	for _, hash := range hashes {
		b.Hashes = append(b.Hashes, hash)
	}
}

