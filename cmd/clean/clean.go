package clean

import (
	"log"
	"os"
	"path"
	"regexp"

	"github.com/mtoohey31/gh-foreach/helper"

	"github.com/spf13/cobra"
)

func NewCleanCmd() *cobra.Command {
	var tmpDirs bool

	cmd := &cobra.Command{
		Use:   "clean",
		Short: "remove cache",
		Long:  `Delete the cached repositories.`,
		Run: func(cmd *cobra.Command, args []string) {
			err := os.RemoveAll(helper.GetCacheDir())

			if err != nil {
				log.Fatalln(err)
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
						err := os.RemoveAll(path.Join("/tmp", item))

						if err != nil {
							log.Fatalln(err)
						}
					}
				}
			}

			// TODO: display spaced saved
		},
	}

	cmd.Flags().BoolVarP(&tmpDirs, "tmp-dirs", "t", false, "clean tmp dirs too")

	return cmd
}
