package run

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"regexp"
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
	Interactive  bool
	RegexString  string
	Regex        *regexp.Regexp
	NoConfirm    bool
}

func NewRunCmd() *cobra.Command {
	opts := runOpts{}

	cmd := &cobra.Command{
		Use:   "run",
		Short: "execute a command",
		Long:  `Execute a command`,
		Run: func(cmd *cobra.Command, args []string) {
			validateOpts(&opts)

			repos := api.GetRepos(opts.Visibility, opts.Affiliations, opts.Languages, opts.Number, *opts.Regex)

			if !opts.NoConfirm {
				names := make([]string, len(repos))
				for i, repo := range repos {
					names[i] = fmt.Sprintf("%s/%s", repo.Owner.Login, repo.Name)
				}
				fmt.Printf("Found:\n%s\nContinue? [Y/n] ", strings.Join(names, "\n"))
				var userResponse string
				fmt.Scanln(&userResponse)
				if !helper.ContainsString([]string{"", "y", "yes"}, strings.ToLower(userResponse)) {
					log.Fatalf("User cancelled with input %s\n", userResponse)
				}
			}

			tmpDir := helper.CreateTmpDir()

			if opts.Interactive {
				wgs := make([]sync.WaitGroup, len(repos))

				for i, r := range repos {
					wgs[i].Add(1)

					go func(i int, r api.Repo) {
						defer wgs[i].Done()

						repo.CreateCopy(r, tmpDir)
					}(i, r)
				}
				for i, r := range repos {
					wgs[i].Wait()

					cmd := exec.Cmd{Path: opts.Shell, Args: []string{opts.Shell, "-c", strings.Join(args, " ")},
						Dir: r.TmpDir(tmpDir), Stdout: os.Stdout, Stderr: os.Stderr, Stdin: os.Stdin}
					cmd.Run()
				}
			} else {
				var wg sync.WaitGroup
				wg.Add(len(repos))

				for _, r := range repos {
					go func(r api.Repo) {
						defer wg.Done()

						repo.CreateCopy(r, tmpDir)

						cmd := exec.Cmd{Path: opts.Shell, Args: []string{opts.Shell, "-c", strings.Join(args, " ")},
							Dir: r.TmpDir(tmpDir), Stdout: os.Stdout, Stderr: os.Stderr}
						cmd.Run()
					}(r)
				}

				wg.Wait()
			}
		},
	}

	// TODO: add user/organization options
	cmd.Flags().StringVarP(&opts.Visibility, "visibility", "v", "all", "filter by repo visibility, one of all, public, or private")
	cmd.Flags().StringArrayVarP(&opts.Affiliations, "affiliation", "a", []string{"owner"}, "filter by affiliation to repo, comma separated list of owner, collaborator, organization_member")
	cmd.Flags().StringArrayVarP(&opts.Languages, "languages", "l", nil, "filter by repos containing one or more of the comma separated list of languages")
	cmd.Flags().StringVarP(&opts.Shell, "shell", "s", os.Getenv("SHELL"), "shell to run command with")
	cmd.Flags().IntVarP(&opts.Number, "number", "n", 30, "max number of repositories operate on")
	cmd.Flags().BoolVarP(&opts.Interactive, "interactive", "i", false, "run commands sequentially and interactively")
	cmd.Flags().StringVarP(&opts.RegexString, "regex", "r", ".*", "filter via regex match on repo name")
	cmd.Flags().BoolVarP(&opts.NoConfirm, "no-confirm", "y", false, "don't ask for confirmation")

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
	opts.Regex, err = regexp.Compile(opts.RegexString)
	if err != nil {
		log.Fatalln(err)
	}
}
