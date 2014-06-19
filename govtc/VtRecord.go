package govtc

import (
	"time"
)

type VtRecord struct {
	Id              int64
	Md5             string
	Sha256          string
	Positives       int64
	Total           int64
	PermaLink       string
	ResponseCode	int
	Scans           string
	ScanDate        time.Time
	UpdateDate      time.Time
}
