package server

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"slices"
	"testing"
)

func Test_fetchData(t *testing.T) {
	cases := []struct {
		name     string
		config   fetchConfig
		status   int
		response string
	}{
		{"200", fetchConfig{method: http.MethodGet}, http.StatusOK, "{'data': 'dummy'}"},
		{"401", fetchConfig{method: http.MethodGet}, http.StatusUnauthorized, "{'error': 'not authorized'}"},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			handler := func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(c.status)
				w.Write([]byte(c.response))
			}
			server := httptest.NewServer(http.HandlerFunc(handler))
			defer server.Close()

			resData, res, err := fetchData(server.URL, fetchConfig{method: c.config.method})

			if res.StatusCode != c.status {
				t.Errorf("%s: Expected %v, got %v", c.name, c.status, res.StatusCode)
			}

			if c.status != http.StatusOK {
				resError := fmt.Errorf("%s", c.response)
				if errors.Is(err, resError) {
					t.Errorf("%s: Expected %s, got %s", c.name, err, resError)
				}

				if resData != nil {
					t.Errorf("%s: Expected nil, got %s", c.name, resData)
				}

			} else {
				if err != nil {
					t.Errorf("%s: Expected nil, got %v", c.name, err)
				}

				if string(resData) != c.response {
					t.Errorf("%s: Expected %s, got %s", c.name, c.response, resData)
				}
			}
		})
	}
}

func Test_renderIcon(t *testing.T) {
	cases := []struct {
		name string
		b    bool
		icon string
	}{
		{"true", true, "yes"},
		{"false", false, "no"},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			res := renderIcon(c.b, c.icon)

			if c.b && res != "yes" {
				t.Errorf("%s: expected 'yes', got %s", c.name, c.icon)
			}

			if !c.b && res != "" {
				t.Errorf("%s: expected '', got %s", c.name, c.icon)
			}
		})
	}
}

func Test_getPagesLinks(t *testing.T) {
	links := []string{
		"<https://gitlab.example.com/api/v4/projects/8/issues/8/notes?page=1&per_page=3>; rel=\"prev\"",
		"<https://gitlab.example.com/api/v4/projects/8/issues/8/notes?page=3&per_page=3>; rel=\"next\"",
		"<https://gitlab.example.com/api/v4/projects/8/issues/8/notes?page=1&per_page=3>; rel=\"first\"",
		"<https://gitlab.example.com/api/v4/projects/8/issues/8/notes?page=3&per_page=3>; rel=\"last\"",
	}

	linksWithoutPrev := links[1:]
	linksWithoutNext := slices.Delete(links, 1, 2)

	cases := []struct {
		name  string
		links []string
		prev  string
		next  string
	}{
		{"all links", links, links[0], links[1]},
		{"without prev", linksWithoutPrev, "", linksWithoutPrev[0]},
		{"without next", linksWithoutNext, linksWithoutNext[0], ""},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			prev, next, err := getPagesLinks(c.links)

			expectedPrev, _ := parseLink(c.prev)
			expectedNext, _ := parseLink(c.next)

			if prev != expectedPrev {
				t.Errorf("%s: Expected %s, got %s", c.name, expectedPrev, prev)
			}

			if next != expectedNext {
				t.Errorf("%s: Expected %s, got %s", c.name, expectedNext, next)
			}

			if err != nil {
				if c.prev != "" {
					t.Errorf("%s: Expected %s, got %s", c.name, c.prev, prev)
				}

				if c.next != "" {
					t.Errorf("%s: Expected %s, got %s", c.name, c.next, next)
				}
			}
		})
	}
}
