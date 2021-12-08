package run

import (
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/mtoohey31/gh-foreach/api"
	"github.com/mtoohey31/gh-foreach/helper"
	"github.com/mtoohey31/gh-foreach/repo"
	"github.com/mtoohey31/which"
	"github.com/spf13/cobra"
)

type runOpts struct {
	Visibility   string
	Affiliations []string
	Languages    []string
	Shell        string
}

// TODO: handle more than 30 responses, include flag for increasing the default of 30 to 100, I'll have to do something else for values beyond that, because the gh api won't handle more
func NewRunCmd() *cobra.Command {
	opts := runOpts{}

	cmd := &cobra.Command{
		Use:   "run",
		Short: "execute a command",
		Long:  `Execute a command`,
		Run: func(cmd *cobra.Command, args []string) {
			validateOpts(&opts)

			repos := api.GetRepos(opts.Visibility, opts.Affiliations, opts.Languages)

			tmpDir := repo.CreateCopies(repos)

			cmds := make([]exec.Cmd, len(repos))

			for i, repo := range repos {
				cmds[i] = exec.Cmd{Path: opts.Shell, Args: []string{opts.Shell, "-c", strings.Join(args, " ")}, Dir: repo.TmpDir(tmpDir), Stdout: os.Stdout, Stderr: os.Stderr}

				err := cmds[i].Start()
				if err != nil {
					log.Println(err)
				}
			}

			for _, cmd := range cmds {
				cmd.Wait()
			}
		},
	}

	cmd.Flags().StringVarP(&opts.Visibility, "visibility", "v", "all", "filter by repo visibility, one of all, public, or private")
	cmd.Flags().StringArrayVarP(&opts.Affiliations, "affiliation", "a", []string{"owner"}, "filter by affiliation to repo, comma separated list of owner, collaborator, organization_member")
	cmd.Flags().StringArrayVarP(&opts.Languages, "languages", "l", nil, "filter by repos containing one or more of the comma separated list of languages")
	cmd.Flags().StringVarP(&opts.Shell, "shell", "s", os.Getenv("SHELL"), "shell to run command with")

	return cmd
}

func validateOpts(opts *runOpts) {
	if !helper.ContainsString([]string{"all", "public", "private"}, opts.Visibility) {
		log.Fatalln("Invalid visibility: ", opts.Visibility)
	}
	validAffiliations := []string{"owner", "collaborator", "organization_member"}
	for _, v := range opts.Affiliations {
		if !helper.ContainsString(validAffiliations, v) {
			log.Fatalln("Invalid affiliation: ", v)
		}
	}
	path, err := which.Which(opts.Shell)
	if err != nil {
		log.Fatalln("Shell not found in path: ", opts.Shell)
	}
	if path != opts.Shell {
		opts.Shell = path
	}
}
