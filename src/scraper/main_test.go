package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestDownloadFile(t *testing.T) {
	t.Run("Download from server", func(t *testing.T) {
		want := "randomString123"
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprint(w, want)
		}))
		defer server.Close()
		resp, err := downloadFile(server.URL)
		if err != nil {
			t.Errorf("Expected err to be nil got %v", err)
		}
		respStr := string(resp)
		if respStr != want {
			t.Errorf("Expected response to be %v, got %v", want, respStr)
		}
	})
	t.Run("Get failed", func(t *testing.T) {
		_, err := downloadFile("http://server:123")
		if err == nil {
			t.Errorf("Expected an error but got nil")
		}
	})
	t.Run("Not OK status Code", func(t *testing.T) {
		want := "403 Forbidden"
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusForbidden)
		}))
		defer server.Close()
		_, err := downloadFile(server.URL)
		if err.Error() != fmt.Sprint(want) {
			t.Errorf("Expected error '%v' but got error '%v'", want, err)
		}
	})
}
