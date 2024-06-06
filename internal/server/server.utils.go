package server

import (
	"fmt"
	"gitlab_tui/internal/icon"
	"io"
	"net/http"
	"os"
)

type fetchConfig struct {
	method string
	params string
	token  string
}

func fetchData(url string, config fetchConfig) ([]byte, error) {
	req, err := http.NewRequest(config.method, url+config.params, nil)
	if err != nil {
		fmt.Print(err.Error())
		os.Exit(1)
	}
	req.Header.Add("PRIVATE-TOKEN", config.token)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Print(err.Error())
		os.Exit(1)
	}
	defer res.Body.Close()

	responseData, err := io.ReadAll(res.Body)
	return responseData, err
}

func renderIcon(b bool) string {
	// i := icon.Dash
	i := icon.Empty
	if b {
		i = icon.Check
	}

	return i
}
