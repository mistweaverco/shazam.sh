package utils

import (
	"log"
	"os"
	"runtime"
)

var ps = string(os.PathSeparator)

func GetDataDirectory() string {
	dir, err := os.UserConfigDir()
	if err != nil {
		log.Fatal(err)
	}
	fullpath := dir + ps + "kelele" + ps
	mkdirerr := os.MkdirAll(fullpath, os.ModePerm)
	if mkdirerr != nil {
		log.Fatal(mkdirerr)
	}
	return fullpath
}

func GetCurrentWorkingDirectory() string {
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	return dir
}

func GetOperatingSystem() string {
	switch runtime.GOOS {
	case "darwin":
		return "macOS"
	case "linux":
		return "Linux"
	case "windows":
		return "Windows"
	default:
		return "Other"
	}
}