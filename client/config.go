package client

type config struct {
	key   string
	email string
}

type Option func(c *config)

func WithAuth(key string, mail string) Option {
	return func(c *config) {
		c.key = key
		c.email = mail
	}
}
