package utils

import (
	"os"
	"strings"
)

func IsDifferentDomain(url string) bool {
	domain := os.Getenv("Domain")

	if url == domain {
		return false
	}

	cleanURl := strings.TrimPrefix(url, "http://")
	cleanURl = strings.TrimPrefix(cleanURl, "https://")
	cleanURl = strings.TrimPrefix(cleanURl, "www.")
	cleanURl = strings.Split(cleanURl, "/")[0]

	return cleanURl != domain

}

func EnsureHTTPPrefix(url string) string {
	if !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") {
		return "http://" + url
	}

	return url
}
