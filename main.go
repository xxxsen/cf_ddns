package main

import (
	"cf_ddns/cf"
	"cf_ddns/config"
	"cf_ddns/provider"
	"context"
	"flag"
	"fmt"

	"github.com/xxxsen/common/logutil"
	"github.com/xxxsen/runner"
	"go.uber.org/zap"
)

var cfg = flag.String("config", "./config", "config file")

func createProviderMap(items []config.ProviderConfig) (map[string]provider.IPProvider, error) {
	rs := make(map[string]provider.IPProvider)
	for _, item := range items {
		p, err := provider.MakeProvider(item.Name, item.Data)
		if err != nil {
			return nil, err
		}
		rs[item.Name] = p
	}
	return rs, nil
}

func findProviderFromMap(ps []string, pm map[string]provider.IPProvider) ([]provider.IPProvider, error) {
	rs := make([]provider.IPProvider, 0, len(ps))
	for _, p := range ps {
		v, ok := pm[p]
		if !ok {
			return nil, fmt.Errorf("provider:%s not found", p)
		}
		rs = append(rs, v)
	}
	return rs, nil
}

func buildRefreshClient(refresherConfigList []config.RefreshCongfig, pm map[string]provider.IPProvider) ([]*cf.Client, error) {
	rs := make([]*cf.Client, 0, len(refresherConfigList))
	for _, item := range refresherConfigList {
		pvList, err := findProviderFromMap(item.Providers, pm)
		if err != nil {
			return nil, fmt.Errorf("find providers failed, err:%v", err)
		}
		if len(pvList) == 0 {
			return nil, fmt.Errorf("no provider found, name:%s", item.Name)
		}
		client, err := cf.New(
			cf.WithAuthKey(item.CloudflareConfig.Key),
			cf.WithAuthMail(item.CloudflareConfig.EMail),
			cf.WithZoneName(item.CloudflareConfig.ZoneName),
			cf.WithRecName(item.CloudflareConfig.RecordName),
			cf.WithRefreshInterval(item.RefreshInterval),
			cf.WithProvider(provider.NewGroup(pvList...)),
			cf.WithClientName(item.Name),
			cf.WithRecordType(item.CloudflareConfig.RecordType),
		)
		if err != nil {
			return nil, fmt.Errorf("create provider failed, name:%s, err:%v", item.Name, err)
		}
		rs = append(rs, client)
	}
	return rs, nil
}

func main() {
	logger := logutil.GetLogger(context.Background())
	logger.Info("support providers", zap.Strings("providers", provider.List()))
	flag.Parse()
	c, err := config.Parse(*cfg)
	if err != nil {
		logger.Fatal("parse config failed", zap.Error(err))
	}
	ps, err := createProviderMap(c.ProviderList)
	if err != nil {
		logger.Fatal("create provider map failed", zap.Error(err))
	}
	clients, err := buildRefreshClient(c.RefreshCongfig, ps)
	if err != nil {
		logger.Fatal("create refresh client failed", zap.Error(err))
	}
	run := runner.New(len(clients))
	for idx, client := range clients {
		client := client
		run.Add(fmt.Sprintf("client_%d", idx), func(ctx context.Context) error {
			client.Start()
			return nil
		})
	}
	run.Run(context.Background())
}
