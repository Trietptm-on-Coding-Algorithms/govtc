package govtc

//
func GetHashTypeFromLength(hash string) int {
	if len(hash) == 32 {
		return HASH_MD5
	} else if len(hash) == 64 {
		return HASH_SHA256
	} else {
		return HASH_UNKNOWN
	}
}
