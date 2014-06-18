package govtc

type CacheChecker struct {
	apiKey 			string
	databasePath 	string
}

func NewCacheChecker(apiKey string, databasePath string) *CacheChecker {
	c := new(CacheChecker)
	c.apiKey = apiKey
	c.databasePath = databasePath
	return c
}

func (c CacheChecker)ProcessFile(inputFile string, mode int) {

}

func (c CacheChecker)ProcessHash(hash string, mode int) {

}

