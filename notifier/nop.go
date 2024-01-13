package notifier

import "context"

const (
	defaultNopNotifierName = "nop"
)

var (
	NopNotifier = &nopNotifier{}
)

type nopNotifier struct {
}

func (n *nopNotifier) Name() string {
	return defaultNopNotifierName
}

func (n *nopNotifier) Notify(ctx context.Context, msg string) error {
	return nil
}

func init() {
	creator := func(interface{}) (INotifier, error) {
		return NopNotifier, nil
	}
	Register(defaultNopNotifierName, creator)
}
