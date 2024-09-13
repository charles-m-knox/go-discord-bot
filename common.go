package main

import (
	"log"
	"os"
	"path"
)

// writeOrAppendToFile overwrites the contents of or appends the contents to the
// provided file if it already exists.
func writeOrAppendToFile(b []byte, basePath, fileName string) {
	err := os.MkdirAll(basePath, 0o755)
	if err != nil {
		log.Printf("failed to create basePath: %v", err.Error())
		return
	}

	mu.Lock()

	f, err := os.OpenFile(path.Join(basePath, fileName), os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0o644)
	if err != nil {
		log.Printf("failed to open file %v: %v", fileName, err.Error())
		return
	}

	defer f.Close()

	defer mu.Unlock()

	_, err = f.Write(b)
	if err != nil {
		log.Printf("failed to write to file %v: %v", fileName, err.Error())
		return
	}
}

// writeToFile overwrites the contents of the provided file.
func writeToFile(b []byte, basePath, fileName string) {
	err := os.MkdirAll(basePath, 0o755)
	if err != nil {
		log.Printf("failed to create basePath: %v", err.Error())
		return
	}

	mu.Lock()

	err = os.WriteFile(path.Join(basePath, fileName), b, 0o644)
	if err != nil {
		log.Printf("failed to open file %v: %v", fileName, err.Error())
		return
	}

	defer mu.Unlock()
}
