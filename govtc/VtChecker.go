package govtc

import (
	"fmt"
	"strings"
	"github.com/williballenthin/govt"
)

type VtChecker struct {
	Id          	int
	apiKey 			string
	databasePath 	string
	checkLimit		int
	resultsQueue	chan govt.FileReportResults
	workerQueue 	chan *VtChecker
	govtClient 		govt.Client
}

//
func NewVtChecker(id int,
				  apiKey string,
				  databasePath string,
				  workerQueue chan *VtChecker,
				  resultsQueue chan govt.FileReportResults) *VtChecker {

	vt := VtChecker {
			Id: 			id,
			workerQueue: 	workerQueue,
			resultsQueue: 	resultsQueue,
			checkLimit: 	4,
			apiKey: 		apiKey,
			databasePath: 	databasePath,
			govtClient: 	govt.Client{Apikey: apiKey}}

	vt.govtClient.UseDefaultUrl()

	return &vt
}

//
func (vt *VtChecker) Work(bd BatchData) {
	fmt.Printf("Worker %d is working on task\n", vt.Id)
	fmt.Println(strings.Join(bd.Hashes, "#"))
	fmt.Println(len(bd.Hashes))


	frr := make([]govt.FileReport, len(bd.Hashes))

	//md5s := []string {"eeb024f2c81f0d55936fb825d21a91d6", "1F4C43ADFD45381CFDAD1FAFEA16B808"}
	reports, err := vt.govtClient.GetFileReports(bd.Hashes)
	if err != nil {
		fmt.Println("Error requesting report: ", err.Error())

		for _, b := range bd.Hashes {
			var fr govt.FileReport
			fr.Status.ResponseCode = -1

			hashType := GetHashTypeFromLength(b)
			if hashType == HASH_MD5 {
				fr.Md5 = b
			} else {
				fr.Sha256 = b
			}

			frr = append(frr, fr)
		}


		return
	} else {
		fmt.Println("VTHCEKC1")
		//for _, r := range *reports {
		//	fmt.Println(r.ResponseCode)
		//	fmt.Println(r.Md5)
		//	fmt.Println(r.ScanDate)
			//var fr govt.FileReport
			//frr = append(frr, fr)
		//}

		vt.resultsQueue <- *reports
		vt.workerQueue <- vt

//		for _, r := range *reports {
//			fmt.Println(r.ResponseCode)
//			fmt.Println(r.Md5)
//			fmt.Println(r.ScanDate)
//
//			for _, s := range r.Scans {
//				if len(s.Result) > 0 {
//					fmt.Println(s.Version)
//					fmt.Println(s.Result)
//				}
//			}
//		}
	}



	//var fr govt.FileReport
	//fr.Md5 = "1"

	//frr = append(frr, fr)

	//vt.resultsQueue <- frr
	//vt.workerQueue <- vt
}
