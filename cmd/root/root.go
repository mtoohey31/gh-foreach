package root

import (
	"github.com/spf13/cobra"
	"mtoohey.com/gh-foreach/cmd/clean"
	"mtoohey.com/gh-foreach/cmd/run"
)

type Repo struct {
	Name      string
	URL       string
	Clone_URL string
}

func NewCmdRoot() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "gh-foreach",
		Short: "execute commands across multiple github repositories",
		// Long:  ``,
	}

	cmd.AddCommand(run.NewRunCmd(), clean.NewCleanCmd())

	return cmd
}
