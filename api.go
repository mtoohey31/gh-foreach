package main

import (
	"log"
	"net/url"
	"path"
	"regexp"
	"strconv"
	"strings"

	"github.com/cli/go-gh"
	"github.com/cli/go-gh/pkg/api"
)

type repo struct {
	Name          string
	Owner         struct{ Login string }
	URL           string
	Clone_URL     string
	Languages_URL string
	Page          string
}

func (r repo) cacheDir() string {
	return path.Join(getCacheDir(), r.Owner.Login, r.Name)
}

func (r repo) tmpDir(tmpRoot string) string {
	return path.Join(tmpRoot, r.Owner.Login, r.Name)
}

func getRepos(visibility string, affiliations []string, languages []string, regex regexp.Regexp) []repo {
	client, err := gh.RESTClient(nil)
	if err != nil {
		log.Fatalln(err)
	}

	values := url.Values{}

	values.Set("visibility", visibility)
	values.Set("affiliation", strings.Join(affiliations, ","))
	values.Set("per_page", "100")

	filteredResponse := []repo{}

	for i := 1; true; i++ {
		response := []repo{}

		values.Set("page", strconv.Itoa(i))

		err = client.Get("user/repos?"+values.Encode(), &response)

		if err != nil {
			log.Fatalln(err)
		}

		if len(response) == 0 {
			break
		}

		for _, r := range response {
			if regex.MatchString(r.Name) && (languages == nil || r.containsSomeLanguage(client, languages)) {
				filteredResponse = append(filteredResponse, r)
			}
		}
	}

	// TODO: handle interaction between language filtering and numbers
	return filteredResponse
}

func (r repo) containsSomeLanguage(client api.RESTClient, languages []string) bool {
	response := map[string]int{}

	err := client.Get(r.Languages_URL, &response)
	if err != nil {
		log.Fatalln(err)
	}

	lowerLanguages := []string{}

	for _, v := range languages {
		lowerLanguages = append(lowerLanguages, strings.ToLower(v))
	}

	for language := range response {
		if containsString(languages, strings.ToLower(language)) {
			return true
		}
	}
	return false
}
