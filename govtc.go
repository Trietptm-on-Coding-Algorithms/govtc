package main

import (
	"bufio"
	"fmt"
	_"log"
	_"os/user"
	"os"
	"strings"
	_"github.com/williballenthin/govt"
	_"github.com/mattn/go-sqlite3"
	"github.com/voxelbrain/goptions"
	"github.com/woanware/govtc/govtc"
	"github.com/woanware/go.util"
	"path"
)

type Options struct {
	Hash   		string        	`goptions:"-h, --hash, description='A single hash (MD5 or SHA256)'"`
	File   		string        	`goptions:"-f, --file, description='File containing hashes'"`
	ApiKey   	string        	`goptions:"-a, --apikey, description='API key for VT REST API'"`
	Database   	string   		`goptions:"-d, --database, description='Path to database directory (defaults to current directory)'"`
	Delimiter   string   		`goptions:"-l, --delimiter, description='The delimiter used for the export. Defaults to \",\"'"`
	Output   	string   		`goptions:"-o, --output, obligatory, description='Output directory (use \".\" for the current dir)'"`
	Mode   		string   		`goptions:"-m, --mode, description='Mode e.g. c = caching, d = database only, l = live'"`
	Help    	goptions.Help 	`goptions:"--help, description='Show help'"`
	configFile 	string
}

const DATABASE_FILE_NAME string  = "govtc.db"
const CONFIG_FILE_NAME string  = "govtc.config"

const (
	MODE_CACHE = 1
	MODE_DB = 2
	MODE_LIVE = 3
)

var (
	options 		*Options
	cacheChecker 	*govtc.CacheChecker
)

func main() {
	options = new(Options)
	goptions.ParseAndFail(options)

	if len(options.Database) == 0 {
		options.Database = goutil.GetApplicationDirectory()
	}

	if len(options.Delimiter) == 0 {
		options.Delimiter = ","
	}

	options.Database = path.Join(options.Database, DATABASE_FILE_NAME)
	options.configFile = path.Join(".", CONFIG_FILE_NAME)

	// The user hasn't supplied an APIKEY so load from the config file
	if len(options.ApiKey) == 0 {
		if goutil.DoesFileExist(options.configFile) == false {
			fmt.Println("The API key has not been set via the command line or config file")
			return
		}

		options.ApiKey = ReadApiKey(options.configFile)

		if len(options.ApiKey) == 0 {
			fmt.Println("The config file does not contain an API key")
			return
		}
	}

	if len(options.Mode) == 0 {
		options.Mode = "c"
	}

	mode := MODE_CACHE
	switch options.Mode {
	case "c":
		mode = MODE_CACHE
	case "d":
		mode = MODE_DB
	case "l":
		mode = MODE_LIVE
	default:
		fmt.Println("Invalid mode e.g. c = caching, d = database only, l = live")
		return
	}

	if len(options.File) == 0 && len(options.Hash) == 0 {
		fmt.Println("Neither the file or hash parameters have been set")
		return
	}

	if len(options.File) > 0 && len(options.Hash) > 0 {
		fmt.Println("Both the file and hash parameters have been set. Choose one or the other")
		return
	}

	if len(strings.TrimSpace(options.File)) > 0 {
		if goutil.DoesFileExist(options.File) == false {
			fmt.Println("The input file does not exist")
			return
		}
	}

	cacheChecker = govtc.NewCacheChecker(options.ApiKey, options.Database)

	if len(strings.TrimSpace(options.File)) > 0 {
		cacheChecker.ProcessFile(options.File, mode)
	} else {
		cacheChecker.ProcessFile(options.Hash, mode)
	}

	/*
	usr, err := user.Current()
	if err != nil {
		log.Fatal( err )
	}
	fmt.Println( usr.HomeDir )*/
}

// Reads the first line from the specified file
func ReadApiKey(configFilePath string) string {
	file, _ := os.Open(configFilePath)
	scanner := bufio.NewScanner(file)
	apiKey := ""
	for scanner.Scan() {
		apiKey = scanner.Text()
		break
	}

	return strings.TrimSpace(apiKey)
}