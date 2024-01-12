package provider

import "fmt"

type IPProvider interface {
	Name() string
	Get() (string, error)
}

type IPGetter func() (string, error)

type IPProviderCreater func(args interface{}) (IPProvider, error)

var mp = make(map[string]IPProviderCreater)

func Register(name string, p IPProviderCreater) {
	mp[name] = p
}

func keys(m map[string]IPProviderCreater) []string {
	rs := make([]string, 0, len(m))
	for k := range m {
		rs = append(rs, k)
	}
	return rs
}

func List() []string {
	return keys(mp)
}

type IPProviderMaker func(name string, data interface{}) (IPProvider, error)

func makeByMap(name string, data interface{}, m map[string]IPProviderCreater) (IPProvider, error) {
	c, ok := m[name]
	if !ok {
		return nil, fmt.Errorf("provider:%s not found", name)
	}
	return c(data)
}

func MakeProvider(name string, data interface{}) (IPProvider, error) {
	c, ok := mp[name]
	if !ok {
		return nil, fmt.Errorf("provider:%s not found", name)
	}
	return c(data)
}
