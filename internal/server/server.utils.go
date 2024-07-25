package server

import (
	"fmt"
	"gitlab_tui/internal/icon"
	"io"
	"net/http"
	"strings"
)

type fetchConfig struct {
	method string
	params string
	token  string
}

func fetchData(url string, config fetchConfig) ([]byte, *http.Response, error) {
	req, err := http.NewRequest(config.method, url+config.params, nil)
	if err != nil {
		return []byte{}, nil, err
	}

	req.Header.Add("PRIVATE-TOKEN", config.token)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return []byte{}, res, err
	}

	defer res.Body.Close()

	responseData, err := io.ReadAll(res.Body)

	if res.StatusCode != 200 {
		err = fmt.Errorf("%s", responseData)
		responseData = nil
	}

	return responseData, res, err
}

func renderIcon(b bool, i string) string {
	if b {
		return i
	}

	return icon.Empty
}

func getPagesLinks(links []string) (string, string) {
	// First page, only has Next page link
	if len(links) == 3 {
		return "", parseLink(links[0])
	}

	return parseLink(links[0]), parseLink(links[1])
}

func parseLink(l string) string {
	b, _, ok := strings.Cut(l, ";")
	if ok {
		trimed := strings.TrimSpace(b)
		parsed := strings.TrimSuffix(strings.TrimPrefix(trimed, "<"), ">")
		return parsed
	}

	return l
}
