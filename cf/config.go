package cf

import (
	"cf_ddns/provider"
	"time"
)

type cfConfig struct {
	authKey    string
	authMail   string
	zoneName   string
	recName    string
	recType    string
	provider   provider.IPProvider
	callback   RefreshCallback
	interval   time.Duration
	clientName string
}

type RefreshCallback func(newip string, oldip string, err error)

var noCallback = func(string, string, error) {}

type Option func(c *cfConfig)

func WithRefreshCallback(cb RefreshCallback) Option {
	return func(c *cfConfig) {
		c.callback = cb
	}
}

func WithAuthKey(key string) Option {
	return func(c *cfConfig) {
		c.authKey = key
	}
}

func WithAuthMail(mail string) Option {
	return func(c *cfConfig) {
		c.authMail = mail
	}
}

func WithZoneName(zone string) Option {
	return func(c *cfConfig) {
		c.zoneName = zone
	}
}

func WithRecName(rec string) Option {
	return func(c *cfConfig) {
		c.recName = rec
	}
}

func WithRecordType(typ string) Option {
	return func(c *cfConfig) {
		c.recType = typ
	}
}

func WithProvider(pv provider.IPProvider) Option {
	return func(c *cfConfig) {
		c.provider = pv
	}
}

func WithRefreshInterval(ts time.Duration) Option {
	return func(c *cfConfig) {
		c.interval = ts
	}
}

func WithClientName(name string) Option {
	return func(c *cfConfig) {
		c.clientName = name
	}
}
