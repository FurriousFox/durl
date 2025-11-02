package util

import (
	"errors"
	"fmt"
	"net/url"
	"regexp"
)

func ParseURL(Url string) (*url.URL, error) {
	var URL *url.URL
	var err error
	var matched bool

	if matched, err = regexp.MatchString(`^[a-zA-Z][a-zA-Z0-9+\-.]*://`, Url); err != nil {
		return nil, errors.New(fmt.Sprint("internal exception:", err))
	}

	if matched {
		URL, err = url.Parse(Url)
	} else {
		URL, err = url.Parse("https://" + Url)
	}

	if err != nil {
		var url2, err2 = url.Parse(Url)
		if err2 == nil {
			URL = url2
		} else {
			return nil, errors.New(fmt.Sprint("invalid url:", err))
		}
	}

	if URL.Scheme != "https" {
		return nil, fmt.Errorf("unsupported protocol '%s'", URL.Scheme)
	}

	if URL.Hostname() == "" {
		return nil, errors.New("hostname required")
	}

	return URL, nil
}
