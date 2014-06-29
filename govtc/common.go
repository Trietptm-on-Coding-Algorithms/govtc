package govtc

// ##### Constants #####################################################################################################

const (
	MODE_CACHE = 1
	MODE_DB = 2
	MODE_LIVE = 3
)

const (
	HASH_UNKNOWN = 1
	HASH_MD5 = 2
	HASH_SHA256 = 3
)

const DATABASE_FILE_NAME string  = "govtc.db"
const CONFIG_FILE_NAME string  = "govtc.config"
