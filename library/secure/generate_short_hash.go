package secure

import (
	"crypto/sha256"
	"encoding/binary"
)

func GenerateShortHash(base string, candidates string) string {
	buf := sha256.Sum256([]byte(base))
	builder := ""
	for i := 0; i*4 < len(buf); i++ {
		runeBase := binary.BigEndian.Uint32(buf[i*4 : (i+1)*4])

		builder += string([]rune(candidates)[runeBase%uint32(len(candidates))])
	}

	return builder
}
