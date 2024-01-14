package refresher

import (
	"cf_ddns/model"
	"context"
	"fmt"
	"time"

	"github.com/xxxsen/common/logutil"
	"go.uber.org/zap"
)

type Refresher struct {
	c      *config
	lastip string
	logger *zap.Logger
}

func New(opts ...Option) (*Refresher, error) {
	c := &config{
		cb: noCB,
	}
	for _, opt := range opts {
		opt(c)
	}
	if c.duration == 0 {
		c.duration = 60 * time.Second
	}
	if c.fn == nil || c.pv == nil {
		return nil, fmt.Errorf("invalid refresher fn or ip provider")
	}
	return &Refresher{c: c, logger: logutil.GetLogger(context.Background()).With(zap.String("name", c.name))}, nil
}

func (r *Refresher) Start() {
	r.logger.Info("start refrestar")
	ticker := time.NewTicker(r.c.duration)
	for range ticker.C {
		r.refresh()
	}
}

func (r *Refresher) refresh() {
	ctx := context.Background()
	logger := r.logger.With(zap.String("old_ip", r.lastip))
	newip, err := r.c.pv.Get()
	if err != nil {
		logger.Error("get new ip failed", zap.Error(err))
		return
	}
	if len(newip) == 0 {
		logger.Error("invalid new ip, empty")
		return
	}
	logger = logger.With(zap.String("new_ip", newip))
	if newip == r.lastip {
		logger.Debug("ip not change, skip next")
		return
	}
	if err := r.c.fn(ctx, newip); err != nil {
		logger.Error("update ip failed", zap.Error(err))
		return
	}
	if err := r.c.cb(ctx, &model.Notification{
		Title:     "[CF_DDNS] IP Change Notification",
		Refresher: r.c.name,
		Domain:    r.c.record,
		Time:      time.Now(),
		NewIP:     newip,
		OldIP:     r.lastip,
	}); err != nil {
		logger.Error("notify ip changed failed", zap.Error(err))
	}
	r.lastip = newip
	logger.Info("update ip succ")
}
