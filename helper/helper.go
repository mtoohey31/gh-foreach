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
