
package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha256"
	"fmt"
	"golang.org/x/sys/unix"
	"os"
)

// generateKey generates a 32-byte AES-256 key from a password.
func generateKey(password string) []byte {
	hash := sha256.Sum256([]byte(password))
	return hash[:]
}

// unpad removes PKCS#7 padding from decrypted data.
func unpad(data []byte) ([]byte, error) {
	length := len(data)
	if length == 0 {
		return nil, fmt.Errorf("decrypted data is empty")
	}
	padding := int(data[length-1])
	if padding > length || padding == 0 {
		return nil, fmt.Errorf("invalid padding size")
	}
	for _, p := range data[length-padding:] {
		if int(p) != padding {
			return nil, fmt.Errorf("invalid padding byte")
		}
	}
	return data[:length-padding], nil
}

// decryptData decrypts data using AES-256-CBC.
func decryptData(encryptedData, key, iv []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	if len(encryptedData)%aes.BlockSize != 0 {
		return nil, fmt.Errorf("encrypted data is not a multiple of the block size")
	}

	decrypted := make([]byte, len(encryptedData))
	mode := cipher.NewCBCDecrypter(block, iv)
	mode.CryptBlocks(decrypted, encryptedData)

	// Remove padding after decryption
	return unpad(decrypted)
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
	if len(encryptedData) < aes.BlockSize {
		fmt.Println("Invalid encrypted data length")
		return
	}
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
	defer unix.Close(fd) // Ensure the file descriptor is closed if not used

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
