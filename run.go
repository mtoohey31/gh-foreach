package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"sync"

	"github.com/alecthomas/kong"
	"mtoohey.com/gh-foreach/api"
	"mtoohey.com/gh-foreach/helper"
	"mtoohey.com/gh-foreach/repo"
)

type Run struct {
	Visibility   string         `short:"v" enum:"all,public,private" default:"all" help:"Filter by repo visibility."`
	Affiliations []string       `short:"a" enum:"owner,collaborator,organization_member" default:"owner" help:"Filter by affiliation to repo."`
	Languages    []string       `short:"l" help:"Filter by repos containing one or more of the provided languages."`
	Shell        string         `short:"s" env:"SHELL" help:"Shell to run command with."`
	Number       int            `short:"n" default:"30" help:"Max number of repositories to operate on."`
	Interactive  bool           `short:"i" help:"Run commands sequentially and interactively."`
	Regex        *regexp.Regexp `short:"r" default:".*" help:"Filter via regex match on repo name."`
	NoConfirm    bool           `short:"N" help:"Don't ask for confirmation."`
	Cleanup      bool           `short:"c" help:"Remove temporary directory after running."`
	Command      []string       `arg:"" help:"The command to run."`
}

func (c *Run) Run(ctx *kong.Context) error {
	repos := api.GetRepos(c.Visibility, c.Affiliations, c.Languages, c.Number, *c.Regex)

	if !c.NoConfirm {
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

	if c.Interactive {
		copyWgs := make([]sync.WaitGroup, len(repos))
		var cleanupWgs []sync.WaitGroup
		if c.Cleanup {
			cleanupWgs = make([]sync.WaitGroup, len(repos))
		}
		for i, r := range repos {
			copyWgs[i].Add(1)

			go func(i int, r api.Repo) {
				defer copyWgs[i].Done()

				repo.CreateCopy(r, tmpDir)
			}(i, r)
		}
		for i, r := range repos {
			copyWgs[i].Wait()

			cmd := exec.Cmd{Path: c.Shell, Args: []string{c.Shell, "-c", strings.Join(c.Command, " ")},
				Dir: r.TmpDir(tmpDir), Stdout: os.Stdout, Stderr: os.Stderr, Stdin: os.Stdin}
			cmd.Run()

			if c.Cleanup {
				cleanupWgs[i].Add(1)

				go func(i int, r api.Repo) {
					defer cleanupWgs[i].Done()

					err := os.RemoveAll(r.TmpDir(tmpDir))
					if err != nil {
						log.Fatalln(err)
					}
				}(i, r)
			}
		}
		if c.Cleanup {
			for i := range cleanupWgs {
				cleanupWgs[i].Wait()
			}
			err := os.RemoveAll(tmpDir)
			if err != nil {
				log.Fatalln(err)
			}
		}
	} else {
		var wg sync.WaitGroup
		wg.Add(len(repos))

		for _, r := range repos {
			go func(r api.Repo) {
				defer wg.Done()

				repo.CreateCopy(r, tmpDir)

				cmd := exec.Cmd{Path: c.Shell, Args: []string{c.Shell, "-c", strings.Join(c.Command, " ")},
					Dir: r.TmpDir(tmpDir), Stdout: os.Stdout, Stderr: os.Stderr}
				cmd.Run()

				if c.Cleanup {
					os.RemoveAll(r.TmpDir(tmpDir))
				}
			}(r)
		}

		wg.Wait()

		if c.Cleanup {
			os.RemoveAll(tmpDir)
		}
	}

	return nil
}

// func NewRunCmd() *cobra.Command {
// 	opts := runOpts{}
//
// 	cmd := &cobra.Command{
// 		Use:   "run",
// 		Short: "execute a command",
// 		Long:  `Execute a command`,
// 		Run: func(cmd *cobra.Command, args []string) {
// 			validateOpts(&opts)
//
// 		},
// 	}
//
// 	// TODO: add user/organization options
// 	cmd.Flags().StringVarP(&opts.Visibility, "visibility", "v", "all", "filter by repo visibility, one of all, public, or private")
// 	cmd.Flags().StringArrayVarP(&opts.Affiliations, "affiliation", "a", []string{"owner"}, "filter by affiliation to repo, comma separated list of owner, collaborator, organization_member")
// 	cmd.Flags().StringArrayVarP(&opts.Languages, "languages", "l", nil, "filter by repos containing one or more of the comma separated list of languages")
// 	cmd.Flags().StringVarP(&opts.Shell, "shell", "s", os.Getenv("SHELL"), "shell to run command with")
// 	cmd.Flags().IntVarP(&opts.Number, "number", "n", 30, "max number of repositories operate on")
// 	cmd.Flags().BoolVarP(&opts.Interactive, "interactive", "i", false, "run commands sequentially and interactively")
// 	cmd.Flags().StringVarP(&opts.RegexString, "regex", "r", ".*", "filter via regex match on repo name")
// 	cmd.Flags().BoolVarP(&opts.NoConfirm, "no-confirm", "y", false, "don't ask for confirmation")
// 	cmd.Flags().BoolVarP(&opts.Cleanup, "cleanup", "c", false, "remove temporary ydirectory after running")
//
// 	return cmd
// }

// func validateOpts(opts *runOpts) {
// 	if !helper.ContainsString([]string{"all", "public", "private"}, opts.Visibility) {
// 		log.Fatalln("Invalid visibility: ", opts.Visibility)
// 	}
// 	validAffiliations := []string{"owner", "collaborator", "organization_member"}
// 	for _, v := range opts.Affiliations {
// 		if !helper.ContainsString(validAffiliations, v) {
// 			log.Fatalln("Invalid affiliation: ", v)
// 		}
// 	}
// 	path, err := which.Which(opts.Shell)
// 	if err != nil {
// 		log.Fatalln("Shell not found in path: ", opts.Shell)
// 	}
// 	if path != opts.Shell {
// 		opts.Shell = path
// 	}
// 	if opts.Number == 0 {
// 		log.Fatalln("Number must be non-zero.")
// 	} else if opts.Number > 100 {
// 		// TODO: support numbers greater than 100 by making multiple API calls
// 		log.Fatalln("Number must be at most 100.")
// 	}
// 	opts.Regex, err = regexp.Compile(opts.RegexString)
// 	if err != nil {
// 		log.Fatalln(err)
// 	}
// }
