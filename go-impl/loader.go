package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha256"
	"fmt"
	"os"
	"golang.org/x/sys/unix"
)

// generateKey generates a 32-byte AES-256 key from a password.
func generateKey(password string) []byte {
	hash := sha256.Sum256([]byte(password))
	return hash[:]
}

// decryptData decrypts data using AES-256-CBC.
func decryptData(encryptedData, key, iv []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	decrypted := make([]byte, len(encryptedData))
	mode := cipher.NewCBCDecrypter(block, iv)
	mode.CryptBlocks(decrypted, encryptedData)

	return decrypted, nil
}

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: loader <encrypted_binary> <password>")
		return
	}

	encryptedFile := os.Args[1]
	password := os.Args[2]

	// Read encrypted data from file
	encryptedData, err := os.ReadFile(encryptedFile)
	if err != nil {
		fmt.Println("Error reading encrypted file:", err)
		return
	}

	// Separate IV and encrypted content
	iv := encryptedData[:aes.BlockSize]
	ciphertext := encryptedData[aes.BlockSize:]

	key := generateKey(password)
	decryptedData, err := decryptData(ciphertext, key, iv)
	if err != nil {
		fmt.Println("Decryption failed:", err)
		return
	}

	// Use MemfdCreate to create an anonymous executable file in memory
	fd, err := unix.MemfdCreate("decrypted_binary", 0)
	if err != nil {
		fmt.Println("Failed to create memfd:", err)
		return
	}

	// Write decrypted binary into the memory file descriptor
	if _, err := unix.Write(fd, decryptedData); err != nil {
		fmt.Println("Failed to write to memfd:", err)
		unix.Close(fd)
		return
	}

	// Execute the binary in memory
	if err := unix.Exec(fmt.Sprintf("/proc/self/fd/%d", fd), []string{fmt.Sprintf("/proc/self/fd/%d", fd)}, os.Environ()); err != nil {
		fmt.Println("Execution failed:", err)
		unix.Close(fd)
		return
	}
}
