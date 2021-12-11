// Encrypt a message behind a passphrase.
package femtioelva

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"io"
	"log"

	"golang.org/x/crypto/nacl/secretbox"
	"golang.org/x/crypto/scrypt"
)

var salt = []byte{0x38, 0xa0, 0x92, 0xac, 0x9d, 0xd1, 0x56, 0x9a}

// generateKey converts a passphrase into a 32-byte key.
// Salt is not random, not safe.
func GenerateKey(passphrase string) [32]byte {
	// Combine passphrase with salt turning it into 32 bytes.
	keyBytes, err := scrypt.Key([]byte(passphrase), salt, 1<<15, 8, 1, 32)
	if err != nil {
		log.Fatal(err)
	}

	var key [32]byte
	copy(key[:], keyBytes)
	return key
}

// Encrypt message using key into base64 string.
func Encrypt(message string, key [32]byte) (string, error) {

	var nonce [24]byte
	if _, err := io.ReadFull(rand.Reader, nonce[:]); err != nil {
		return "", err
	}
	// Put the unique nonce at the beginning of the payload.
	cipher := secretbox.Seal(nonce[:], []byte(message), &nonce, &key)
	return base64.StdEncoding.EncodeToString(cipher), nil
}

// Decrypt payload base64 string using key, return plain text.
func Decrypt(payload string, key [32]byte) (string, error) {

	cipher, err := base64.StdEncoding.DecodeString(payload)
	if err != nil {
		return "", err
	}

	var decryptNonce [24]byte
	copy(decryptNonce[:], cipher[:24])
	plain, ok := secretbox.Open(nil, cipher[24:], &decryptNonce, &key)
	if !ok {
		return "", errors.New("secretbox did not open.")
	}

	return string(plain), nil
}
