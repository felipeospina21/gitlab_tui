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

type Author struct {
	Name string `json:"name"`
}

type DiscussionNote struct {
	ID         int    `json:"id"`
	Type       string `json:"type"`
	Body       string `json:"body"`
	Author     Author
	CreatedAt  string `json:"created_at"`
	Resolved   bool   `json:"resolved"`
	Resolvable bool   `json:"resolvable"`
	ResolvedBy Author `json:"resolved_by"`
	System     bool   `json:"system"`
}

type GetMergeRequestDiscussionsResponse struct {
	ID            string           `json:"id"`
	HasSingleNote bool             `json:"individual_notes"`
	Notes         []DiscussionNote `json:"notes"`
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
		isResolved := func() bool {
			for _, item := range item.Notes {
				if item.Resolvable && item.Resolved {
					return true
				}
			}
			return false
		}()

		n := table.Row{
			createdAt,
			author,
			strconv.Itoa(comments),
			renderIcon(isResolved, icon.Check),
			t,
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

		if item.System {
			writeString(writeStringArgs{builder: &content, withSpace: true, sameLine: true, s: icon.Person})
			writeString(writeStringArgs{builder: &content, withSpace: true, sameLine: true, s: author})
			writeString(writeStringArgs{builder: &content, s: body})
			writeString(writeStringArgs{builder: &content, s: separator})
		} else {
			writeString(writeStringArgs{builder: &content, withSpace: true, sameLine: true, s: icon.Clock})
			writeString(writeStringArgs{builder: &content, withSpace: true, s: createdAt})
			writeString(writeStringArgs{builder: &content, withSpace: true, sameLine: true, s: icon.Person})
			writeString(writeStringArgs{builder: &content, withSpace: true, s: author})
			writeString(writeStringArgs{builder: &content, withSpace: true, sameLine: true, s: icon.Discussion})
			writeString(writeStringArgs{builder: &content, s: body})
			writeString(writeStringArgs{builder: &content, s: separator})

		}

	}

	content.WriteString(resolved.String())
	return content.String(), nil
}

type writeStringArgs struct {
	builder   *strings.Builder
	s         string
	sameLine  bool
	withSpace bool
}

func writeString(args writeStringArgs) *strings.Builder {
	if args.withSpace {
		args.builder.WriteString(args.s + " ")
	} else {
		args.builder.WriteString(args.s)
	}

	if !args.sameLine {
		args.builder.WriteString("\n\n")
	}
	return args.builder
}
