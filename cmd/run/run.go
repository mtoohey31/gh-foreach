package run

import (
	"log"
	"os"
	"os/exec"

	"github.com/mtoohey31/gh-foreach/api"
	"github.com/mtoohey31/gh-foreach/helper"
	"github.com/mtoohey31/gh-foreach/repo"

	"github.com/spf13/cobra"
)

type runOpts struct {
	Visibility   string
	Affiliations []string
	Languages    []string
}

// TODO: handle more than 30 responses, include flag for increasing the default of 30 to 100, I'll have to do something else for values beyond that, because the gh api won't handle more
func NewRunCmd() *cobra.Command {
	opts := runOpts{}

	cmd := &cobra.Command{
		Use:   "run",
		Short: "execute a command",
		Long:  `Execute a command`,
		Run: func(cmd *cobra.Command, args []string) {
			validateOpts(opts)

			repos := api.GetRepos(opts.Visibility, opts.Affiliations, opts.Languages)

			tmpDir := repo.CreateCopies(repos)

			for _, repo := range repos {
				// cmd := exec.Cmd{Path: "/bin/sh", Args: []string{"-c", "ls"}, Dir: repo.TmpDir(tmpDir)}
				cmd := exec.Command(os.Getenv("SHELL"), "-c", "ls")
				cmd.Dir = repo.TmpDir(tmpDir)

				out, _ := cmd.Output()
				log.Println(string(out))
			}
		},
	}

	cmd.Flags().StringVarP(&opts.Visibility, "visibility", "v", "all", "filter by repo visibility, one of all, public, or private, default all")
	cmd.Flags().StringArrayVarP(&opts.Affiliations, "affiliation", "a", []string{"owner"}, "filter by affiliation to repo, comma separated list of owner, collaborator, organization_member, default owner")
	cmd.Flags().StringArrayVarP(&opts.Languages, "languages", "l", nil, "filter by repos containing one or more of the comma separated list of languages, default nil")

	return cmd
}

func validateOpts(opts runOpts) {
	if !helper.ContainsString([]string{"all", "public", "private"}, opts.Visibility) {
		log.Fatalln("Invalid visibility: ", opts.Visibility)
	}
	validAffiliations := []string{"owner", "collaborator", "organization_member"}
	for _, v := range opts.Affiliations {
		if !helper.ContainsString(validAffiliations, v) {
			log.Fatalln("Invalid affiliation: ", v)
		}
	}

}
