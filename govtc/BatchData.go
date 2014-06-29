package govtc

import (
	"fmt"
)


type BatchData struct {
	Hashes []string
}

//
func (b *BatchData) Add(hashes []string) {
	for _, hash := range hashes {
		fmt.Println(hash)
		b.Hashes = append(b.Hashes, hash)
	}
}

//
func (b *BatchData) AddHash(hash string) {
	b.Hashes = append(b.Hashes, hash)
}

