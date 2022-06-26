package net

import "net/http"

// CheckHTTPConnection checks if the given HTTP client can connect to the given
// URL.
func CheckHTTPConnection(client *http.Client, url string) bool {
	resp, err := http.Get(url)
	if err != nil {
		return false
	}
	switch resp.StatusCode {
	case http.StatusOK, http.StatusNoContent:
		return true
	}
	return false
}
