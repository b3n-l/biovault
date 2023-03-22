package lib

import (
	"fmt"
	"log"
	"testing"
)

func TestEncryptDecrypt(t *testing.T) {
	// Encrypt an arbitrary string
	passphrase := []byte("MySecretPassphrase")
	plaintext := []byte("Hello, world!")

	// Encrypt the data
	encryptedData, err := EncryptBytes(plaintext, passphrase)
	if err != nil {
		log.Fatal(err)
	}

	// Decrypt the encrypted data
	decryptedData, err := DecryptBytes(encryptedData, passphrase)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(plaintext)
	fmt.Println(decryptedData)
}
