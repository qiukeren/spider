package utils

import (
	"errors"
	"net/url"
	"strings"
)

func FormatUrl(url1, site string) (string, error) {

	if strings.HasPrefix(url1, "javascript:") {
		return "", errors.New("starts with 'javascript:'")
	}

	u, err := url.Parse(url1)
	if err != nil {
		return "", err
	}
	if u.Scheme == "" {
		u.Scheme = "http"

	}
	if u.Host == "" {
		u.Host = site
	}
	return u.String(), nil
}

func ParseUrl(url1 string) (*url.URL, error) {
	return url.Parse(url1)
}

func IsCurrentSite(url1, site string) bool {
	u, err := url.Parse(url1)
	if err != nil {
		return false
	}
	if u.Host == site || u.Host == "" {
		return true
	}
	return false
}
