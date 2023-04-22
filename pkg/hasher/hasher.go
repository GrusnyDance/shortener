package hasher

import (
	"crypto/sha256"
	"encoding/binary"
)

var hashSymbols = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz_"
var hashSymbolsLen uint64 = 63

func Apply(input string) string {
	output := []byte(hashSymbols[0:10])
	// Hash the input data using SHA-256
	hashStr := sha256.Sum256([]byte(input))
	hash := binary.BigEndian.Uint64(hashStr[:])

	// Byte shift to divide by 32
	for i := 0; i < 10 && hash > 0; i, hash = i+1, hash>>5 {
		output[i] = hashSymbols[hash%hashSymbolsLen]
	}
	return string(output)
}
