package notifier

type config struct {
	user string
	code string
	addr string
}

type Option func(c *config)

func WithAuth(user string, code string) Option {
	return func(c *config) {
		c.user = user
		c.code = code
	}
}

func WithAddr(addr string) Option {
	return func(c *config) {
		c.addr = addr
	}
}

type jsConfig struct {
	User string `json:"user"`
	Code string `json:"code"`
	Addr string `json:"addr"`
}
