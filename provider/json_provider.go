package provider

import (
	"encoding/json"
	"fmt"

	"github.com/xxxsen/common/utils"
	"github.com/yalp/jsonpath"
)

type jsonProvider struct {
	url  string
	path string
}

func (p *jsonProvider) Name() string {
	return ProviderJson
}

func (p *jsonProvider) Get() (string, error) {
	return findIPByURL(p.url, func(b []byte) (string, error) {
		var m interface{}
		if err := json.Unmarshal(b, &m); err != nil {
			return "", err
		}
		v, err := jsonpath.Read(m, p.path)
		if err != nil {
			return "", err
		}
		ip, ok := v.(string)
		if !ok {
			return "", fmt.Errorf("not string key")
		}

		return ip, nil
	})
}

func createJsonProvider(args interface{}) (IProvider, error) {
	c := &JsonProviderConfig{}
	if err := utils.ConvStructJson(args, c); err != nil {
		return nil, err
	}
	if len(c.URL) == 0 || len(c.Path) == 0 {
		return nil, fmt.Errorf("invalid params")
	}
	return &jsonProvider{url: c.URL, path: c.Path}, nil
}

func init() {
	Register(ProviderJson, createJsonProvider)
}
