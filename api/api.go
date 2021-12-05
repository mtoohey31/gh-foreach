package api

import (
	"log"
	"net/url"
	"path"
	"strings"

	"github.com/cli/go-gh"
	"github.com/cli/go-gh/pkg/api"
	"github.com/mtoohey31/gh-foreach/helper"
)

type Repo struct {
	Name          string
	Owner         struct{ Login string }
	URL           string
	Clone_URL     string
	Languages_URL string
}

type Languages map[string]int

func (repo Repo) CacheDir() string {
	return path.Join(helper.GetCacheDir(), repo.Owner.Login, repo.Name)
}

func (repo Repo) TmpDir(tmpRoot string) string {
	return path.Join(tmpRoot, repo.Owner.Login, repo.Name)
}

func GetRepos(visibility string, affiliations []string, languages []string) []Repo {
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

	if languages == nil {
		return response
	}

	filteredResponse := []Repo{}

	for _, repo := range response {
		if repo.containsSomeLanguage(client, languages) {
			filteredResponse = append(filteredResponse, repo)
		}
	}

	return filteredResponse
}

func (repo Repo) containsSomeLanguage(client api.RESTClient, languages []string) bool {
	response := Languages{}

	err := client.Get(repo.Languages_URL, &response)
	if err != nil {
		log.Fatalln(err)
	}

	lowerLanguages := []string{}

	for _, v := range languages {
		lowerLanguages = append(lowerLanguages, strings.ToLower(v))
	}

	for language := range response {
		if helper.ContainsString(languages, strings.ToLower(language)) {
			return true
		}
	}
	return false
}
