package server

import (
	"fmt"
	"gitlab_tui/internal/icon"
	"io"
	"net/http"
)

type fetchConfig struct {
	method string
	params string
	token  string
}

func fetchData(url string, config fetchConfig) ([]byte, int, error) {
	req, err := http.NewRequest(config.method, url+config.params, nil)
	if err != nil {
		return []byte{}, 0, err
	}

	req.Header.Add("PRIVATE-TOKEN", config.token)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return []byte{}, 0, err
	}

	defer res.Body.Close()

	responseData, err := io.ReadAll(res.Body)

	if res.StatusCode != 200 {
		err = fmt.Errorf("%s", responseData)
		responseData = nil
	}

	return responseData, res.StatusCode, err
}

func renderIcon(b bool) string {
	// i := icon.Dash
	i := icon.Empty
	if b {
		i = icon.Check
	}

	return i
}
