package provider

import (
	"net/http"
	"time"
)

var gClient *http.Client

func init() {
	gClient = &http.Client{
		Transport: &http.Transport{
			MaxIdleConns:        10,
			MaxIdleConnsPerHost: 2,
			IdleConnTimeout:     120 * time.Second,
		},
		Timeout: 5 * time.Second,
	}
}

func defaultClient() *http.Client {
	return gClient
}
