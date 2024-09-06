package server

import (
	"fmt"
	"gitlab_tui/config"
	"gitlab_tui/internal/icon"
	"io"
	"net/http"
	"regexp"
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

func getPagesLinks(links []string) (string, string, error) {
	// First page, only has Next page link
	if len(links) == 3 {
		l, rel := parseLink(links[0])
		ok, err := isPrevLink(rel)
		if err != nil {
			return "", "", err
		}
		if ok {
			return l, "", nil
		} else {
			return "", l, nil
		}
	}

	prev, _ := parseLink(links[0])
	next, _ := parseLink(links[1])

	return prev, next, nil
}

func parseLink(l string) (string, string) {
	b, a, ok := strings.Cut(l, ";")
	if ok {
		trimed := strings.TrimSpace(b)
		parsed := strings.TrimSuffix(strings.TrimPrefix(trimed, "<"), ">")
		return parsed, a
	}

	return l, ""
}

func isPrevLink(s string) (bool, error) {
	pattern := "prev"
	match, err := regexp.MatchString(pattern, s)
	if err != nil {
		return false, err
	}

	return match, nil
}

type requestsType = uint

const (
	MRListReq requestsType = iota
	MRCommentsReq
	MRDiscussionsReq
	MRSingleDiscussionReq
	MRPipelinesReq
	JobsReq
	IssuesReq
	MergeMRReq
)

type ReqData struct {
	ProjectID    string
	MrID         string
	MergeReqIDD  string
	PipelineID   string
	DiscussionID string
}

func BuildURL(t requestsType, reqData ReqData, config config.Config) string {
	switch t {
	case MRListReq:
		return fmt.Sprintf("%s/%s/projects/%s/merge_requests", config.BaseURL, config.APIVersion, reqData.ProjectID)

	case MRCommentsReq:
		return fmt.Sprintf("%s/%s/projects/%s/merge_requests/%s/notes", config.BaseURL, config.APIVersion, reqData.ProjectID, reqData.MrID)

	case MRDiscussionsReq:
		return fmt.Sprintf("%s/%s/projects/%s/merge_requests/%s/discussions", config.BaseURL, config.APIVersion, reqData.ProjectID, reqData.MrID)

	case MRSingleDiscussionReq:
		return fmt.Sprintf("%s/%s/projects/%s/merge_requests/%s/discussions/%s", config.BaseURL, config.APIVersion, reqData.ProjectID, reqData.MrID, reqData.DiscussionID)

	case MRPipelinesReq:
		return fmt.Sprintf("%s/%s/projects/%s/merge_requests/%s/pipelines", config.BaseURL, config.APIVersion, reqData.ProjectID, reqData.MrID)

	case JobsReq:
		return fmt.Sprintf("%s/%s/projects/%s/pipelines/%s/jobs", config.BaseURL, config.APIVersion, reqData.ProjectID, reqData.PipelineID)

	case IssuesReq:
		return fmt.Sprintf("%s/%s/projects/%s/issues", config.BaseURL, config.APIVersion, reqData.ProjectID)

	case MergeMRReq:
		return fmt.Sprintf("%s/%s/projects/%s/merge_requests/%s/merge", config.BaseURL, config.APIVersion, reqData.ProjectID, reqData.MergeReqIDD)

	default:
		return ""
	}
}
