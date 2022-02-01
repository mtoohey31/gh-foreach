package main

import (
	"errors"
	"log"
	"os"

	"github.com/go-git/go-git/v5"
	"github.com/otiai10/copy"
)

func createCopy(r repo, tmpDir string) {
	ensureUpToDate(r)
	err := copy.Copy(r.cacheDir(), r.tmpDir(tmpDir))
	if err != nil {
		log.Fatalln(err)
	}
}

func ensureUpToDate(r repo) {
	if exists(r) {
		pull(r)
	} else {
		clone(r)
	}
}

func exists(r repo) bool {
	_, err := os.Stat(r.cacheDir())
	return err == nil
}

func clone(r repo) {
	_, err := git.PlainClone(r.cacheDir(), false, &git.CloneOptions{
		URL:      r.Clone_URL,
		Progress: os.Stdout,
	})
	if err != nil {
		log.Fatalln(err)
	}
}

func pull(r repo) {
	gr, err := git.PlainOpen(r.cacheDir())
	if err != nil {
		log.Fatalln(err)
	}
	w, err := gr.Worktree()
	if err != nil {
		log.Fatalln(err)
	}
	err = w.Pull(&git.PullOptions{})
	if err != nil && !errors.Is(err, git.NoErrAlreadyUpToDate) {
		log.Fatalln(err)
	}
}
