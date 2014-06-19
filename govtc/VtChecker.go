package govtc

import (
	"fmt"
)

type VtChecker struct {
	Id          	int
	apiKey 			string
	databasePath 	string
	checkLimit		int
	WorkQueue		chan BatchCheck
	WorkerQueue 	chan chan BatchCheck
	QuitChan    	chan bool
}

//
func NewVtChecker(id int,
				  apiKey string,
				  databasePath string,
				  workerQueue chan chan BatchCheck) VtChecker {

	vt := VtChecker {
			Id: 			id,
			WorkerQueue: 	workerQueue,
			WorkQueue: 		make(chan BatchCheck),
			checkLimit: 	4,
			apiKey: 		apiKey,
			databasePath: 	databasePath,
			QuitChan:    	make(chan bool)}

	return vt
}

// This function "starts" the worker by starting a goroutine, that is an infinite "for-select" loop.
func (vt VtChecker) Start() {
	go func() {
		for {
			// Add ourselves into the worker queue.
			vt.WorkerQueue <- vt.WorkQueue

			select {
			case work := <-vt.WorkQueue:
				// Receive a work request.
				fmt.Printf("worker%d: Received work request (%d)\n",vt.Id, work.Hashes)

				//time.Sleep(work.Delay * time.Millisecond)

			case <-vt.QuitChan:
				// We have been asked to stop.
				//fmt.Printf("worker%d stopping\n", w.ID)
				return
			}
		}
	}()
}
