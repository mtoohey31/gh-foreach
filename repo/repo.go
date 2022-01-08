package repo

import (
	"errors"
	"log"
	"os"

	"github.com/go-git/go-git/v5"
	"github.com/otiai10/copy"
	"mtoohey.com/gh-foreach/api"
)

func CreateCopy(repo api.Repo, tmpDir string) {
	ensureUpToDate(repo)
	err := copy.Copy(repo.CacheDir(), repo.TmpDir(tmpDir))
	if err != nil {
		log.Fatalln(err)
	}
}

func ensureUpToDate(repo api.Repo) {
	if exists(repo) {
		pull(repo)
	} else {
		clone(repo)
	}
}

func exists(repo api.Repo) bool {
	_, err := os.Stat(repo.CacheDir())
	return err == nil
}

func clone(repo api.Repo) {
	_, err := git.PlainClone(repo.CacheDir(), false, &git.CloneOptions{
		URL:      repo.Clone_URL,
		Progress: os.Stdout,
	})
	if err != nil {
		log.Fatalln(err)
	}
}

func pull(repo api.Repo) {
	r, err := git.PlainOpen(repo.CacheDir())
	if err != nil {
		log.Fatalln(err)
	}
	w, err := r.Worktree()
	if err != nil {
		log.Fatalln(err)
	}
	err = w.Pull(&git.PullOptions{})
	if err != nil && !errors.Is(err, git.NoErrAlreadyUpToDate) {
		log.Fatalln(err)
	}
}
