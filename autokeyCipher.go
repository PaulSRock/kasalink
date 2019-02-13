package kasalink

import (
	"bytes"
	"encoding/binary"
)

//Kasa uses something called auto key ciphering for communicating with their devices. It's trivial, but does make
//communications non-human readable.
func encrypt(plaintext string) []byte {
	var (
		n                   = len(plaintext)
		buf                 = new(bytes.Buffer)
		err                 error
		ciphertext, payload []byte
		key                 byte
		i                   int
	)

	if err = binary.Write(buf, binary.BigEndian, uint32(n)); err != nil {
		panic(err)
	}
	ciphertext = []byte(buf.Bytes())

	key = byte(0xAB)
	payload = make([]byte, n)
	for i = 0; i < n; i++ {
		payload[i] = plaintext[i] ^ key
		key = payload[i]
	}

	for i := 0; i < len(payload); i++ {
		ciphertext = append(ciphertext, payload[i])
	}

	return ciphertext
}

func decrypt(ciphertext []byte) string {
	var (
		key     = byte(0xAB)
		nextKey byte
		i, n    int
	)
	n = len(ciphertext)
	for i = 0; i < n; i++ {
		nextKey = ciphertext[i]
		ciphertext[i] = ciphertext[i] ^ key
		key = nextKey
	}
	return string(ciphertext)
}
