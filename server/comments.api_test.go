package server

import (
	"encoding/json"
	"gitlab_tui/tui/components/table"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
)

func Test_getMergeRequestComments(t *testing.T) {
	mockAPIResponse := []GetMergeRequestCommentsResponse{
		{ID: 1, Type: "type", Body: "", CreatedAt: "03/12/2024", Resolved: false, Author: Author{Name: "name"}},
	}

	handler := func(w http.ResponseWriter, r *http.Request) {
		b, err := json.Marshal(mockAPIResponse)
		if err != nil {
			t.Errorf("Marshall error: %v", err)
		}

		w.Write(b)
	}
	server := httptest.NewServer(http.HandlerFunc(handler))
	defer server.Close()

	res, _ := GetMergeRequestComments(server.URL)

	if len(res) < 1 {
		t.Errorf("Expected response to have %v, got %v", len(mockAPIResponse), len(res))
	}

	for i, item := range res {
		if item[table.CommentsCols.ID.Idx] != strconv.Itoa(mockAPIResponse[i].ID) {
			t.Errorf("Expected %v, got %s", mockAPIResponse[i].ID, item[table.CommentsCols.ID.Idx])
		}
	}
}
