package provider

import (
	"fmt"

	"github.com/xxxsen/common/utils"
)

type rawTextProvider struct {
	name string
	url  string
}

func (p *rawTextProvider) Name() string {
	return ProviderRawText
}

func (p *rawTextProvider) Get() (string, error) {
	return findIPByURL(p.url, func(s []byte) (string, error) {
		return string(s), nil
	})
}

func createRawTextProvider(args interface{}) (IProvider, error) {
	c := &rawTextProviderConfig{}
	if err := utils.ConvStructJson(args, c); err != nil {
		return nil, err
	}
	if len(c.URL) == 0 {
		return nil, fmt.Errorf("invalid url")
	}
	return &rawTextProvider{url: c.URL}, nil
}

func init() {
	Register(ProviderRawText, createRawTextProvider)
}
