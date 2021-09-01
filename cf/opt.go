package cf

type Config struct {
	authKey    string
	authMail   string
	zoneName   string
	recName    string
	ipProvider string
	callback   RefreshCallback
}

type RefreshCallback func(newip string, oldip string, err error)

type Option func(c *Config)

func WithRefreshCallback(cb RefreshCallback) Option {
	return func(c *Config) {
		c.callback = cb
	}
}

func WithAuthKey(key string) Option {
	return func(c *Config) {
		c.authKey = key
	}
}

func WithAuthMail(mail string) Option {
	return func(c *Config) {
		c.authMail = mail
	}
}

func WithZoneName(zone string) Option {
	return func(c *Config) {
		c.zoneName = zone
	}
}

func WithRecName(rec string) Option {
	return func(c *Config) {
		c.recName = rec
	}
}

func WithIPProvider(ips string) Option {
	return func(c *Config) {
		c.ipProvider = ips
	}
}
