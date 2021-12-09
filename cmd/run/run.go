package run

import (
	"log"
	"os"
	"os/exec"
	"strings"
	"sync"

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
	Number       int
}

func NewRunCmd() *cobra.Command {
	opts := runOpts{}

	cmd := &cobra.Command{
		Use:   "run",
		Short: "execute a command",
		Long:  `Execute a command`,
		Run: func(cmd *cobra.Command, args []string) {
			validateOpts(&opts)

			repos := api.GetRepos(opts.Visibility, opts.Affiliations, opts.Languages, opts.Number)

			tmpDir := helper.CreateTmpDir()

			var wg sync.WaitGroup
			wg.Add(len(repos))

			for i, r := range repos {
				go func(i int, r api.Repo) {
					defer wg.Done()
					repo.CreateCopy(r, tmpDir)

					cmd := exec.Cmd{Path: opts.Shell, Args: []string{opts.Shell, "-c", strings.Join(args, " ")}, Dir: r.TmpDir(tmpDir), Stdout: os.Stdout, Stderr: os.Stderr}
					cmd.Start()
					cmd.Wait()
				}(i, r)
			}

			wg.Wait()
		},
	}

	// TODO: add user/organization options
	cmd.Flags().StringVarP(&opts.Visibility, "visibility", "v", "all", "filter by repo visibility, one of all, public, or private")
	cmd.Flags().StringArrayVarP(&opts.Affiliations, "affiliation", "a", []string{"owner"}, "filter by affiliation to repo, comma separated list of owner, collaborator, organization_member")
	cmd.Flags().StringArrayVarP(&opts.Languages, "languages", "l", nil, "filter by repos containing one or more of the comma separated list of languages")
	cmd.Flags().StringVarP(&opts.Shell, "shell", "s", os.Getenv("SHELL"), "shell to run command with")
	cmd.Flags().IntVarP(&opts.Number, "number", "n", 30, "max number of repositories operate on")

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
	if opts.Number == 0 {
		log.Fatalln("Number must be non-zero.")
	} else if opts.Number > 100 {
		// TODO: support numbers greater than 100 by making multiple API calls
		log.Fatalln("Number must be at most 100.")
	}
}
