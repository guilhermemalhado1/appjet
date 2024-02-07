package services

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
	"io/ioutil"
	"strings"
)

const (
	keySize = 32 // AES-256 key size
)

func EncryptAndSaveToken(token string) error {
	// Generate a random key for encryption
	key := make([]byte, keySize)
	if _, err := rand.Read(key); err != nil {
		return fmt.Errorf("error generating encryption key: %w", err)
	}

	// Encrypt the token
	encryptedToken, err := encrypt(token, key)
	if err != nil {
		return fmt.Errorf("error encrypting token: %w", err)
	}

	// Encode the key to base64 to use it as filename
	keyFilename := base64.URLEncoding.EncodeToString(key)

	// Save the encrypted token to a file with the encoded key as filename
	err = ioutil.WriteFile(keyFilename+".security", encryptedToken, 0644)
	if err != nil {
		return fmt.Errorf("error writing encrypted token to file: %w", err)
	}

	return nil
}

// DecryptToken decrypts the token from the encrypted file.
func DecryptToken() (string, error) {
	// Find the .security file in the current directory
	files, err := ioutil.ReadDir(".")
	if err != nil {
		return "", fmt.Errorf("error reading directory: %w", err)
	}

	var encryptedFile string
	for _, file := range files {
		if strings.HasSuffix(file.Name(), ".security") {
			encryptedFile = file.Name()
			break
		}
	}

	if encryptedFile == "" {
		return "", fmt.Errorf("no .security file found in the current directory")
	}

	// Decode the key from the filename
	keyBytes, err := base64.URLEncoding.DecodeString(strings.TrimSuffix(encryptedFile, ".security"))
	if err != nil {
		return "", fmt.Errorf("error decoding key from filename: %w", err)
	}

	// Read the encrypted token from the file
	encryptedToken, err := ioutil.ReadFile(encryptedFile)
	if err != nil {
		return "", fmt.Errorf("error reading encrypted token from file: %w", err)
	}

	// Decrypt the token
	decryptedToken, err := decrypt(encryptedToken, keyBytes)
	if err != nil {
		return "", fmt.Errorf("error decrypting token: %w", err)
	}

	return decryptedToken, nil
}

// encrypt encrypts the plaintext using AES-GCM encryption with the provided key.
func encrypt(plaintext string, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("error creating cipher block: %w", err)
	}

	// Generate a random nonce
	nonce := make([]byte, 12)
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, fmt.Errorf("error generating nonce: %w", err)
	}

	// Create a GCM cipher with the given key
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("error creating GCM cipher: %w", err)
	}

	// Encrypt the plaintext
	ciphertext := aesGCM.Seal(nil, nonce, []byte(plaintext), nil)

	// Prepend the nonce to the ciphertext
	ciphertext = append(nonce, ciphertext...)

	return ciphertext, nil
}

// decrypt decrypts the ciphertext using AES-GCM decryption with the provided key.
func decrypt(ciphertext []byte, key []byte) (string, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", fmt.Errorf("error creating cipher block: %w", err)
	}

	// Extract the nonce from the ciphertext
	nonce := ciphertext[:12]
	ciphertext = ciphertext[12:]

	// Create a GCM cipher with the given key
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("error creating GCM cipher: %w", err)
	}

	// Decrypt the ciphertext
	plaintext, err := aesGCM.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", fmt.Errorf("error decrypting ciphertext: %w", err)
	}

	return string(plaintext), nil
}
