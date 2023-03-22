package lib

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"golang.org/x/crypto/pbkdf2"
	"io"
	"log"
)

func EncryptBytes(data []byte, passphrase []byte) ([]byte, error) {
	// Generate a random salt
	salt := make([]byte, 8)
	if _, err := io.ReadFull(rand.Reader, salt); err != nil {
		return nil, err
	}
	log.Println("salt", hex.EncodeToString(salt))
	// Derive a key from the passphrase and salt using PBKDF2
	key := pbkdf2.Key(passphrase, salt, 500000, 32, sha256.New)

	// Generate a new AES-GCM cipher using the derived key
	block, err := aes.NewCipher(key)
	log.Println("key", hex.EncodeToString(key))
	if err != nil {
		return nil, err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	// Generate a new nonce for this encryption
	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}
	log.Println("nonce:", hex.EncodeToString(nonce))

	// Encrypt the data using AES-GCM
	ciphertext := gcm.Seal(nil, nonce, data, nil)
	log.Println("plain ciphertext", hex.EncodeToString(ciphertext))
	// Append the salt and nonce to the ciphertext
	encryptedData := append(salt, nonce...)
	encryptedData = append(encryptedData, ciphertext...)

	return encryptedData, nil
}

func DecryptBytes(data []byte, passphrase []byte) ([]byte, error) {
	// Extract the salt and nonce from the encrypted data
	salt := data[:8]
	nonce := data[8:20]
	ciphertext := data[20:]

	// Derive a key from the passphrase and salt using PBKDF2
	key := pbkdf2.Key(passphrase, salt, 500000, 32, sha256.New)

	// Generate a new AES-GCM cipher using the derived key
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	// Decrypt the ciphertext using AES-GCM
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, err
	}

	return plaintext, nil
}
