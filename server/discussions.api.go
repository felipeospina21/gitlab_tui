package server

import (
	"encoding/json"
	"gitlab_tui/config"
	"gitlab_tui/internal/icon"
	"gitlab_tui/tui/components/table"
	"slices"
	"strconv"
	"strings"
)

type DiscussionNotes struct {
	GetMergeRequestCommentsResponse
	Resolvable bool   `json:"resolvable"`
	ResolvedBy Author `json:"resolved_by"`
}

type GetMergeRequestDiscussionsResponse struct {
	ID            string            `json:"id"`
	HasSingleNote bool              `json:"individual_notes"`
	Notes         []DiscussionNotes `json:"notes"`
}

func GetMergeRequestDiscussions(url string) ([]table.Row, error) {
	token := config.GlobalConfig.APIToken

	responseData, _, err := fetchData(url, fetchConfig{method: "GET", token: token})
	if err != nil {
		return nil, err
	}

	var r []GetMergeRequestDiscussionsResponse
	if err = json.Unmarshal(responseData, &r); err != nil {
		return nil, err
	}

	// transforms response interface to match tbl Row
	var rows []table.Row
	var discussion []table.Row
	for _, item := range r {
		createdAt := table.FormatTime(item.Notes[0].CreatedAt)
		t := item.Notes[0].Type
		author := item.Notes[0].Author.Name
		comments := len(item.Notes)

		n := table.Row{
			createdAt,
			author,
			t,
			strconv.Itoa(comments),
			item.Notes[0].Body,
			item.ID,
		}

		if t != "" {
			discussion = append(discussion, n)
		} else {
			rows = append(rows, n)
		}

	}

	return slices.Concat(discussion, rows), nil
}

func GetMergeRequestSingleDiscussion(url string) (string, error) {
	token := config.GlobalConfig.APIToken

	responseData, _, err := fetchData(url, fetchConfig{method: "GET", token: token})
	if err != nil {
		return "", err
	}

	var r GetMergeRequestDiscussionsResponse
	if err = json.Unmarshal(responseData, &r); err != nil {
		return "", err
	}

	var content strings.Builder
	var resolved strings.Builder
	for _, item := range r.Notes {
		createdAt := table.FormatTime(item.CreatedAt)
		author := item.Author.Name
		body := item.Body
		separator := strings.Repeat("-", 5)

		if item.Resolved {
			writeString(writeStringArgs{builder: &resolved, s: icon.Check + " ", icon: true})
			// writeString(writeStringArgs{builder: &resolved, s: icon.Discussion + " ", icon: true})
			writeString(writeStringArgs{builder: &resolved, s: body, italize: true})
			writeString(writeStringArgs{builder: &resolved, s: separator})
		} else {
			writeString(writeStringArgs{builder: &content, s: icon.Clock + " ", icon: true})
			writeString(writeStringArgs{builder: &content, s: createdAt})
			writeString(writeStringArgs{builder: &content, s: icon.Person + " ", icon: true})
			writeString(writeStringArgs{builder: &content, s: author})
			writeString(writeStringArgs{builder: &content, s: icon.Discussion + " ", icon: true})
			writeString(writeStringArgs{builder: &content, s: body})
			writeString(writeStringArgs{builder: &content, s: separator})

		}

	}

	content.WriteString(resolved.String())
	return content.String(), nil
}

type writeStringArgs struct {
	builder *strings.Builder
	s       string
	italize bool
	icon    bool
}

func writeString(args writeStringArgs) *strings.Builder {
	if args.italize {
		args.builder.WriteString("*")
		args.builder.WriteString(args.s)
		args.builder.WriteString("*")

	} else {
		args.builder.WriteString(args.s)
	}

	if !args.icon {
		args.builder.WriteString("\n\n")
	}
	return args.builder
}
