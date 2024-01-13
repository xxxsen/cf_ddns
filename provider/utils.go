package provider

import (
	"fmt"
	"io"
	"net/http"
)

func makeRequest(uri string) ([]byte, error) {
	req, err := http.NewRequest(http.MethodGet, uri, nil)
	if err != nil {
		return nil, err
	}
	rsp, err := defaultClient().Do(req)
	if err != nil {
		return nil, err
	}
	defer rsp.Body.Close()
	if rsp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("invalid status_code:%d", rsp.StatusCode)
	}
	raw, err := io.ReadAll(rsp.Body)
	if err != nil {
		return nil, err
	}
	return raw, nil
}

func findIPByURL(uri string, parser func([]byte) (string, error)) (string, error) {
	data, err := makeRequest(uri)
	if err != nil {
		return "", err
	}
	return parser(data)
}
