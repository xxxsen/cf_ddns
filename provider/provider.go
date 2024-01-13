package provider

import "fmt"

type IProvider interface {
	Name() string
	Get() (string, error)
}

type ProviderCreater func(args interface{}) (IProvider, error)

var mp = make(map[string]ProviderCreater)

func Register(name string, p ProviderCreater) {
	mp[name] = p
}

func keys(m map[string]ProviderCreater) []string {
	rs := make([]string, 0, len(m))
	for k := range m {
		rs = append(rs, k)
	}
	return rs
}

func List() []string {
	return keys(mp)
}

type ProviderMaker func(name string, data interface{}) (IProvider, error)

func MakeProvider(name string, data interface{}) (IProvider, error) {
	c, ok := mp[name]
	if !ok {
		return nil, fmt.Errorf("provider:%s not found", name)
	}
	return c(data)
}
