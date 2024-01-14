package notifier

import (
	"cf_ddns/model"
	"context"
)

var (
	NopNotifier = &nopNotifier{}
)

type nopNotifier struct {
}

func (n *nopNotifier) Name() string {
	return NameNop
}

func (n *nopNotifier) Notify(ctx context.Context, msg *model.Notification) error {
	return nil
}

func init() {
	creator := func(interface{}) (INotifier, error) {
		return NopNotifier, nil
	}
	Register(NameNop, creator)
}
