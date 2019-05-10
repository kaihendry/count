package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestHTML(t *testing.T) {
	server := httptest.NewServer(routes())
	defer server.Close()
	for _, i := range []int{1, 2} {
		resp, err := http.Get(server.URL)
		if err != nil {
			t.Fatal(err)
		}
		if resp.StatusCode != http.StatusOK {
			t.Fatalf("got %d, want %d", resp.StatusCode, http.StatusOK)
		}
		expected := fmt.Sprintf("<title>Count: %d</title>", i)
		actual, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			t.Fatal(err)
		}
		if !strings.Contains(string(actual), expected) {
			t.Errorf("got %s, want %s", actual, expected)
		}
	}
}

func TestJSON(t *testing.T) {
	server := httptest.NewServer(routes())
	defer server.Close()
	for _, i := range []int{1, 2} {
		resp, err := http.Get(server.URL + "/inc/")
		if err != nil {
			t.Fatal(err)
		}
		if resp.StatusCode != http.StatusOK {
			t.Fatalf("got %d, want %d", resp.StatusCode, http.StatusOK)
		}
		expected := fmt.Sprintf("%d", i)
		actual, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			t.Fatal(err)
		}
		if !strings.Contains(string(actual), expected) {
			t.Errorf("got %s, want %s", actual, expected)
		}
	}
}
