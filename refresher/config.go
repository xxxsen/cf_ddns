package refresher

import (
	"cf_ddns/provider"
	"context"
	"time"
)

type RefresherFunc func(ctx context.Context, ip string) error
type CallbackFunc func(ctx context.Context, name string, oldip string, newip string) error

var noCB CallbackFunc = func(ctx context.Context, name, oldip, newip string) error { return nil }

type config struct {
	name     string
	fn       RefresherFunc
	pv       provider.IProvider
	duration time.Duration
	cb       CallbackFunc
}

type Option func(c *config)

func WithRefresherFunc(fn RefresherFunc) Option {
	return func(c *config) {
		c.fn = fn
	}
}

func WithIPProvider(pv provider.IProvider) Option {
	return func(c *config) {
		c.pv = pv
	}
}

func WithInterval(t time.Duration) Option {
	return func(c *config) {
		c.duration = t
	}
}

func WithName(name string) Option {
	return func(c *config) {
		c.name = name
	}
}

func WithCallback(cb CallbackFunc) Option {
	return func(c *config) {
		c.cb = cb
	}
}
