package notifier

import (
	"context"
	"fmt"
)

type INotifier interface {
	Name() string
	Notify(ctx context.Context, msg string) error
}

type NotifierCreator func(args interface{}) (INotifier, error)

var mp = make(map[string]NotifierCreator)

func Register(name string, c NotifierCreator) {
	mp[name] = c
}

func MakeNotifier(name string, data interface{}) (INotifier, error) {
	c, ok := mp[name]
	if !ok {
		return nil, fmt.Errorf("notifier:%s not found", name)
	}
	return c(data)
}

func List() []string {
	rs := make([]string, 0, len(mp))
	for k := range mp {
		rs = append(rs, k)
	}
	return rs
}
