package root

import (
	"github.com/mtoohey31/gh-foreach/cmd/clean"
	"github.com/mtoohey31/gh-foreach/cmd/run"
	"github.com/spf13/cobra"
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

	cmd.AddCommand(run.NewRunCmd())
	cmd.AddCommand(clean.NewCleanCmd())

	return cmd
}
