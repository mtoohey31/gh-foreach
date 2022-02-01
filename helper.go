package main

import (
	"io/ioutil"
	"log"
	"os"
	"path"
)

func getCacheDir() string {
	dir, err := os.UserCacheDir()
	if err != nil {
		log.Fatalln(err)
	}
	return path.Join(dir, "gh-foreach")
}

func containsString(arr []string, item string) bool {
	for _, v := range arr {
		if v == item {
			return true
		}
	}
	return false
}

func createTmpDir() string {
	tmpDir, err := ioutil.TempDir("/tmp", "gh-foreach")
	if err != nil {
		log.Fatalln(err)
	}
	return tmpDir
}
