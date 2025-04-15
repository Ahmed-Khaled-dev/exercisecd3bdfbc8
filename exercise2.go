package main

import (
	"crypto/rand"
	"fmt"
	"os"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: ./exercise2 <filePath> (e.g., ./exercise2 randomfile)")
		os.Exit(1)
	}

	path := os.Args[1]

	if err := Shred(path); err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}

	fmt.Println("File shredded successfully!")
}

func Shred(filePath string) error {
	// Check if the file exists and is a regular file (not symlink/directory)
	fileInfo, err := os.Lstat(filePath)
	if err != nil {
		return fmt.Errorf("file not found: %w", err)
	}

	if fileInfo.Mode()&os.ModeSymlink != 0 {
		return fmt.Errorf("cannot shred symlinks")
	}

	if fileInfo.IsDir() {
		return fmt.Errorf("cannot shred directories")
	}

	fileSizeInBytes := fileInfo.Size()

	// Overwrite 3 times
	for i := 0; i < 3; i++ {
		file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_TRUNC, 0)
		if err != nil {
			return fmt.Errorf("error opening file: %w", err)
		}

		overwriteData := make([]byte, fileSizeInBytes)
		// Fill 'overwriteData' with random bytes
		if _, err := rand.Read(overwriteData); err != nil {
			file.Close()
			return fmt.Errorf("error generating random data: %w", err)
		}

		if _, err := file.Write(overwriteData); err != nil {
			file.Close()
			return fmt.Errorf("error writing random data to the file: %w", err)
		}

		// Ensures data is written on the file
		if err := file.Sync(); err != nil {
			file.Close()
			return fmt.Errorf("error syncing data: %w", err)
		}

		file.Close()
	}

	// Delete the file
	return os.Remove(filePath)
}
