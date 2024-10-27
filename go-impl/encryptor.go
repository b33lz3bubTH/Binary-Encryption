package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"io"
	"os"
)

// generateKey generates a 32-byte AES-256 key from a password.
func generateSecureKey(password string) []byte {
	hash := sha256.Sum256([]byte(password))
	return hash[:]
}

// encryptData encrypts data using AES-256-CBC.
func encryptData(data, key []byte) ([]byte, []byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, nil, err
	}

	iv := make([]byte, aes.BlockSize)
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil, nil, err
	}

	ciphertext := make([]byte, len(data))
	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(ciphertext, data)

	return ciphertext, iv, nil
}

func main() {
	if len(os.Args) < 4 {
		fmt.Println("Usage: encryptor <input_binary> <output_encrypted> <password>")
		return
	}

	inputFile := os.Args[1]
	outputFile := os.Args[2]
	password := os.Args[3]

	key := generateSecureKey(password)

	data, err := os.ReadFile(inputFile)
	if err != nil {
		fmt.Println("Error reading input file:", err)
		return
	}

	encryptedData, iv, err := encryptData(data, key)
	if err != nil {
		fmt.Println("Error encrypting data:", err)
		return
	}

	output, err := os.Create(outputFile)
	if err != nil {
		fmt.Println("Error creating output file:", err)
		return
	}
	defer output.Close()

	// Write IV and ciphertext to output
	output.Write(iv)
	output.Write(encryptedData)

	fmt.Printf("Encrypted binary created as %s\n", outputFile)
}
