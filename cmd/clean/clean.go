package clean

import (
	"github.com/mtoohey31/gh-foreach/helper"
	"os"

	"github.com/spf13/cobra"
)

func NewCleanCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "clean",
		Short: "remove cache",
		Long:  `Delete the cached repositories.`,
		Run: func(cmd *cobra.Command, args []string) {
			os.Remove(helper.GetCacheDir())
		},
	}

	return cmd
}
