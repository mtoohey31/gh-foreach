package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"
	"regexp"

	"github.com/alecthomas/kong"

	"github.com/c2h5oh/datasize"
)

type Clean struct {
	TmpDirs bool `help:"Clean temporary directories too." short:"t"`
}

func (c *Clean) Run(ctx *kong.Context) error {
	cacheDir := GetCacheDir()
	_, err := os.Stat(cacheDir)

	var size int64 = 0

	if !errors.Is(err, os.ErrNotExist) {
		if err != nil {
			log.Fatalln(err)
		}

		size, err = dirSize(cacheDir)

		if err != nil {
			log.Fatalln(err)
		}

		err = os.RemoveAll(cacheDir)

		if err != nil {
			log.Fatalln(err)
		}
	}

	if c.TmpDirs {
		tmp, err := os.Open("/tmp")
		if err != nil {
			log.Fatalln(err)
		}
		dirnames, err := tmp.Readdirnames(0)
		if err != nil {
			log.Fatalln(err)
		}
		tmpRe := regexp.MustCompile("^gh-foreach.*")
		for _, item := range dirnames {
			if tmpRe.MatchString(item) {
				itemPath := path.Join("/tmp", item)

				newSize, err := dirSize(itemPath)

				if err != nil {
					log.Fatalln(err)
				}

				size += newSize

				err = os.RemoveAll(itemPath)

				if err != nil {
					log.Fatalln(err)
				}
			}
		}
	}

	fmt.Printf("Freed %s\n", (datasize.ByteSize(size)).HumanReadable())

	return nil
}

// source: https://stackoverflow.com/a/32482941
func dirSize(path string) (int64, error) {
	var size int64
	err := filepath.Walk(path, func(_ string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			size += info.Size()
		}
		return err
	})
	return size, err
}
