package provider

import "strings"

type groupProvider struct {
	lst []IProvider
}

func NewGroup(ps ...IProvider) IProvider {
	return &groupProvider{lst: ps}
}

func (p *groupProvider) Name() string {
	names := make([]string, 0, len(p.lst))
	for _, item := range p.lst {
		names = append(names, item.Name())
	}
	return "Group:[" + strings.Join(names, ",") + "]"
}

func (p *groupProvider) Get() (string, error) {
	var rerr error
	for _, item := range p.lst {
		v, err := item.Get()
		if err == nil {
			return v, nil
		}
		rerr = err
	}
	return "", rerr
}
