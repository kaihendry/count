package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestMyHandler(t *testing.T) {
	server := httptest.NewServer(routes())
	defer server.Close()

	for _, i := range []int{1, 2} {
		resp, err := http.Get(server.URL)
		if err != nil {
			t.Fatal(err)
		}
		if resp.StatusCode != 200 {
			t.Fatalf("Received non-200 response: %d\n", resp.StatusCode)
		}
		expected := fmt.Sprintf("<title>Count: %d</title>", i)
		actual, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			t.Fatal(err)
		}
		// https://github.com/golang/go/wiki/CodeReviewComments#useful-test-failures
		if !strings.Contains(string(actual), expected) {
			t.Errorf("got %s, want %s", actual, expected)
		}
	}
}
