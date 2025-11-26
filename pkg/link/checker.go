package link

import (
	"net/http"
	"net/url"
)

const (
	Available    = "available"
	NotAvailable = "not_available"
)

func CheckLink(link string) string {
	if len(link) < 7 || (link[:7] != "http://" && link[:8] != "https://") {
		link = "https://" + link
	}

	_, err := url.Parse(link)
	if err != nil {
		return NotAvailable
	}

	resp, err := http.Head(link)
	if err != nil {
		return NotAvailable
	}

	if resp.StatusCode >= http.StatusOK && resp.StatusCode < http.StatusBadRequest {
		return Available
	}
	return NotAvailable
}
