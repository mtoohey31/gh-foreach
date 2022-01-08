package clean

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"
	"regexp"

	"mtoohey.com/gh-foreach/helper"

	"github.com/spf13/cobra"

	"github.com/c2h5oh/datasize"
)

func NewCleanCmd() *cobra.Command {
	var tmpDirs bool

	cmd := &cobra.Command{
		Use:   "clean",
		Short: "remove cache",
		Long:  `Delete the cached repositories.`,
		Run: func(cmd *cobra.Command, args []string) {
			cacheDir := helper.GetCacheDir()
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

			if tmpDirs {
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
		},
	}

	cmd.Flags().BoolVarP(&tmpDirs, "tmp-dirs", "t", false, "clean tmp dirs too")

	return cmd
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
