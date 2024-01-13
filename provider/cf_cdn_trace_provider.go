package provider

import (
	"fmt"
	"net"
	"strings"

	"github.com/xxxsen/common/utils"
)

type cfCDNTraceProvider struct {
	url string
}

func (t *cfCDNTraceProvider) Name() string {
	return ProviderCFCDNTrace
}

func (t *cfCDNTraceProvider) findIPFromText(raw []byte) (string, error) {
	lines := strings.Split(string(raw), "\n")
	for _, line := range lines {
		if !strings.HasPrefix(line, "ip=") {
			continue
		}
		ct := strings.Split(line, "=")
		if len(ct) != 2 {
			return "", fmt.Errorf("invalid ip line, ip line data:%s", line)
		}
		if ip := net.ParseIP(ct[1]); ip == nil {
			return "", fmt.Errorf("invalid ip data:%s", ct[1])
		}
		return ct[1], nil
	}
	return "", fmt.Errorf("no ip line found")
}

func (t *cfCDNTraceProvider) Get() (string, error) {
	return findIPByURL(t.url, t.findIPFromText)
}

func createCFCDNTraceProvider(args interface{}) (IProvider, error) {
	c := &cfCDNTraceProviderConfig{}
	if err := utils.ConvStructJson(args, c); err != nil {
		return nil, err
	}
	if len(c.URL) == 0 {
		return nil, fmt.Errorf("invalid url")
	}
	return &cfCDNTraceProvider{url: c.URL}, nil
}

func init() {
	Register(ProviderCFCDNTrace, createCFCDNTraceProvider)
}
