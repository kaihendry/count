package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestCount(t *testing.T) {

	tests := []struct {
		name string
		uri  string
		want string
	}{
		{"HTML 1", "/", fmt.Sprintf("<title>Count: %d</title>", 1)},
		{"HTML 2", "/", fmt.Sprintf("<title>Count: %d</title>", 2)},
		{"JSON 3", "/inc", "3"},
		{"JSON 4", "/inc", "4"},
	}

	ts := httptest.NewServer(routes())
	defer ts.Close()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			url := ts.URL + tt.uri
			resp, _ := http.Get(url)
			if resp.StatusCode != http.StatusOK {
				t.Fatalf("got %d, want %d", resp.StatusCode, http.StatusOK)
			}
			respBody, _ := ioutil.ReadAll(resp.Body)
			resp.Body.Close()
			got := string(respBody)
			if !strings.Contains(got, tt.want) {
				t.Errorf("got %s, Want %s", got, tt.want)
			}
		})
	}

}
