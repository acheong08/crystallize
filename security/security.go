package security

import (
	_ "github.com/acheong08/crystals-go/dilithium"
	kyber "github.com/acheong08/crystals-go/kyber"

	"encoding/base64"
	"strings"
)

func Generate() ([]byte, []byte) {
	k := kyber.NewKyber1024()
	pk, sk := k.PKEKeyGen(nil)
	return pk, sk
}

func Encrypt(pubkey []byte, message string) string {
	// Split message into 32 byte chunks and append 0s if necessary
	var chunks [][]byte
	for i := 0; i < len(message); i += 32 {
		// If len is less than 32, append 0s
		if len(message[i:]) < 32 {
			chunks = append(chunks, []byte(message[i:]+string(make([]byte, 32-len(message[i:])))))
		} else {
			chunks = append(chunks, []byte(message[i:i+32]))
		}
	}
	// Encrypt each chunk
	k := kyber.NewKyber1024()
	var ciphertext [][]byte
	for _, chunk := range chunks {
		ciphertext = append(ciphertext, k.Encrypt(pubkey, chunk, nil))
	}
	// Convert ciphertext to string
	var out string
	for _, chunk := range ciphertext {
		out += encode(chunk) + ","
	}
	return out
}

func Decrypt(privkey []byte, ciphertext string) string {
	// Parse the ciphertext
	var chunks [][]byte
	for _, chunk := range strings.Split(ciphertext, ",") {
		if chunk != "" {
			chunks = append(chunks, decode(chunk))
		}
	}
	// Decrypt each chunk
	k := kyber.NewKyber1024()
	var plaintext []byte
	for _, chunk := range chunks {
		plaintext = append(plaintext, k.Decrypt(privkey, chunk)...)
	}
	return string(plaintext)
}

func encode(data []byte) string {
	return base64.StdEncoding.EncodeToString(data)
}

func decode(data string) []byte {
	decoded, _ := base64.StdEncoding.DecodeString(data)
	return decoded
}
