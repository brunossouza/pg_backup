package controllers

import (
	"log"
	"os"
)

// CheckError check error
func CheckError(err error) {
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
}

// FileExists check if file exists
func FileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

// CreateBackupDirectory CreateBackupDirectory
func CreateBackupDirectory(path string) {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		os.MkdirAll(path, os.FileMode(0777))
	}
}
