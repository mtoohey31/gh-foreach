package helper

import (
	"log"
	"os"
	"path"
)

func GetCacheDir() string {
	dir, err := os.UserCacheDir()
	if err != nil {
		log.Fatalln(err)
	}
	return path.Join(dir, "gh-foreach")
}

func ContainsString(arr []string, item string) bool {
	for _, v := range arr {
		if v == item {
			return true
		}
	}
	return false
}
