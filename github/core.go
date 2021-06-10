package github

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"
)
type Repository string

func (repository Repository) Author() string {
	return strings.Split(string(repository), "/")[0]
}

func (repository Repository) Name() string {
	return strings.Split(string(repository), "/")[1]
}

func NewGithubRequest(url string, data interface{}) error {
	response, err := http.Get("https://api.github.com/" + url)

	if err != nil {
		return err
	}

	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(response.Body)

	body, err := io.ReadAll(response.Body)

	if err != nil {
		return err
	}

	err = json.Unmarshal(body, &data)

	if err != nil {
		return err
	}

	return nil
}

