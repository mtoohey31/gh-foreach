// cspell:ignore fatih cleanup yoinked

package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"sync"

	"github.com/alecthomas/kong"
	"github.com/fatih/color"
)

type run struct {
	Visibility   string         `short:"v" enum:"all,public,private" default:"all" help:"Filter by repo visibility."`
	Affiliations []string       `short:"a" enum:"owner,collaborator,organization_member" default:"owner" help:"Filter by affiliation to repo."`
	Languages    []string       `short:"l" help:"Filter by repos containing one or more of the provided languages."`
	Shell        string         `short:"s" env:"SHELL" help:"Shell to run command with."`
	Interactive  bool           `short:"i" help:"Run commands sequentially and interactively."`
	Regex        *regexp.Regexp `short:"r" default:".*" help:"Filter via regex match on repo name."`
	NoConfirm    bool           `short:"N" help:"Don't ask for confirmation."`
	Cleanup      bool           `short:"c" help:"Remove temporary directory after running."`
	Command      []string       `arg:"" help:"The command to run."`
}

// shamelessly yoinked from: https://github.com/docker/compose/blob/7c47673d4af41d79900e6c70bc1a3f9f17bdd387/pkg/utils/writer.go#L1-L62
type repoOutput struct {
	name       string
	maxNameLen int
	color      color.Attribute
	out        io.Writer
	buffer     bytes.Buffer
}

func (o repoOutput) Write(b []byte) (n int, err error) {
	n, err = o.buffer.Write(b)
	if err != nil {
		return n, err
	}
	for {
		b = o.buffer.Bytes()
		index := bytes.Index(b, []byte{'\n'})
		if index < 0 {
			break
		}
		_, err := color.New(o.color).Fprintf(o.out, "%s %s|", o.name,
			strings.Repeat(" ", o.maxNameLen-len(o.name)))
		if err != nil {
			return n, err
		}
		_, err = o.out.Write([]byte{' '})
		if err != nil {
			return n, err
		}
		_, err = o.out.Write(o.buffer.Next(index + 1))
		if err != nil {
			return n, err
		}
	}
	return n, nil
}

func NewOutputPair(name string, maxNameLen int, index int) (repoOutput, repoOutput) {
	var c int
	wrappedIndex := index % 14
	if wrappedIndex < 7 {
		c = 31 + wrappedIndex
	} else {
		c = 84 + wrappedIndex
	}
	return repoOutput{name, maxNameLen, color.Attribute(c), os.Stdout, bytes.Buffer{}},
		repoOutput{name, maxNameLen, color.Attribute(c), os.Stderr, bytes.Buffer{}}
}

func (c *run) Run(ctx *kong.Context) error {
	repos := getRepos(c.Visibility, c.Affiliations, c.Languages, *c.Regex)

	if !c.NoConfirm {
		names := make([]string, len(repos))
		for i, r := range repos {
			names[i] = fmt.Sprintf("%s/%s", r.Owner.Login, r.Name)
		}
		fmt.Printf("%s\nContinue? [Y/n] ", strings.Join(names, "\n"))
		var userResponse string
		fmt.Scanln(&userResponse)
		if !containsString([]string{"", "y", "yes"}, strings.ToLower(userResponse)) {
			log.Fatalf("User cancelled with input %s\n", userResponse)
		}
	}

	tmpDir := createTmpDir()

	if c.Interactive {
		copyWgs := make([]sync.WaitGroup, len(repos))
		var cleanupWgs []sync.WaitGroup
		if c.Cleanup {
			cleanupWgs = make([]sync.WaitGroup, len(repos))
		}
		for i, r := range repos {
			copyWgs[i].Add(1)

			go func(i int, r repo) {
				defer copyWgs[i].Done()

				createCopy(r, tmpDir)
			}(i, r)
		}
		for i, r := range repos {
			copyWgs[i].Wait()

			cmd := exec.Cmd{Path: c.Shell, Args: []string{c.Shell, "-c", strings.Join(c.Command, " ")},
				Dir: r.tmpDir(tmpDir), Stdout: os.Stdout, Stderr: os.Stderr, Stdin: os.Stdin}
			cmd.Run()

			if c.Cleanup {
				cleanupWgs[i].Add(1)

				go func(i int, r repo) {
					defer cleanupWgs[i].Done()

					err := os.RemoveAll(r.tmpDir(tmpDir))
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
		maxNameLen := 0

		for _, r := range repos {
			if len(r.Name) > maxNameLen {
				maxNameLen = len(r.Name)
			}
		}

		var wg sync.WaitGroup
		wg.Add(len(repos))

		for i, r := range repos {
			go func(r repo, i int) {
				defer wg.Done()

				createCopy(r, tmpDir)

				stdout, stderr := NewOutputPair(r.Name, maxNameLen, i)
				cmd := exec.Cmd{Path: c.Shell, Args: []string{c.Shell, "-c", strings.Join(c.Command, " ")},
					Dir: r.tmpDir(tmpDir), Stdout: stdout, Stderr: stderr}
				cmd.Run()

				if c.Cleanup {
					os.RemoveAll(r.tmpDir(tmpDir))
				}
			}(r, i)
		}

		wg.Wait()

		if c.Cleanup {
			os.RemoveAll(tmpDir)
		}
	}

	return nil
}
