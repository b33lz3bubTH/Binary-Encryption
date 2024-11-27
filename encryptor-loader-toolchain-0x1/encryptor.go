
package main

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"io"
	"os"
)

// Generate a secure key using SHA-256 hashing
func generateSecureKey(password string) []byte {
	hash := sha256.Sum256([]byte(password))
	return hash[:]
}

// Pad the data using PKCS#7 padding
func pad(data []byte, blockSize int) []byte {
	padding := blockSize - len(data)%blockSize
	padText := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(data, padText...)
}

// Remove padding from the data (for decryption)
func unpad(data []byte) []byte {
	length := len(data)
	padding := int(data[length-1])
	return data[:length-padding]
}

// Encrypt data using AES CBC mode with the given key
func encryptData(data, key []byte) ([]byte, []byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, nil, err
	}

	iv := make([]byte, aes.BlockSize)
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil, nil, err
	}

	// Pad the data to make it a multiple of the block size
	paddedData := pad(data, aes.BlockSize)
	ciphertext := make([]byte, len(paddedData))
	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(ciphertext, paddedData)

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

	// Generate encryption key from the provided password
	key := generateSecureKey(password)

	// Read input file
	data, err := os.ReadFile(inputFile)
	if err != nil {
		fmt.Println("Error reading input file:", err)
		return
	}

	// Encrypt the data
	encryptedData, iv, err := encryptData(data, key)
	if err != nil {
		fmt.Println("Error encrypting data:", err)
		return
	}

	// Write the IV and encrypted data to the output file
	output, err := os.Create(outputFile)
	if err != nil {
		fmt.Println("Error creating output file:", err)
		return
	}
	defer output.Close()

	output.Write(iv)
	output.Write(encryptedData)

	fmt.Printf("Encrypted binary created as %s\n", outputFile)
}
