package api

import (
	"log"
	"net/url"
	"path"
	"strings"

	"github.com/cli/go-gh"
	"github.com/mtoohey31/gh-foreach/helper"
)

type Repo struct {
	Name      string
	Owner     struct{ Login string }
	URL       string
	Clone_URL string
}

func (repo Repo) CacheDir() string {
	return path.Join(helper.GetCacheDir(), repo.Owner.Login, repo.Name)
}

func (repo Repo) TmpDir(tmpRoot string) string {
	return path.Join(tmpRoot, repo.Owner.Login, repo.Name)
}

func GetRepos(visibility string, affiliations []string) []Repo {
	client, err := gh.RESTClient(nil)
	if err != nil {
		log.Fatalln(err)
	}

	values := url.Values{}

	values.Set("visibility", visibility)
	values.Set("affiliation", strings.Join(affiliations, ","))

	response := []Repo{}

	err = client.Get("user/repos?"+values.Encode(), &response)
	if err != nil {
		log.Fatalln(err)
	}

	return response
}
