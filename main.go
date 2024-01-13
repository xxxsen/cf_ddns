package main

import (
	"cf_ddns/adaptor"
	"cf_ddns/client"
	"cf_ddns/config"
	"cf_ddns/provider"
	"cf_ddns/refresher"
	"context"
	"flag"
	"fmt"
	"time"

	"github.com/xxxsen/common/logutil"
	"github.com/xxxsen/runner"
	"go.uber.org/zap"
)

var cfg = flag.String("config", "./config", "config file")

func createProviderMap(items []config.ProviderConfig) (map[string]provider.IProvider, error) {
	rs := make(map[string]provider.IProvider)
	for _, item := range items {
		p, err := provider.MakeProvider(item.Type, item.Data)
		if err != nil {
			return nil, err
		}
		rs[item.Name] = p
	}
	return rs, nil
}

func findProviderFromMap(ps []string, pm map[string]provider.IProvider) ([]provider.IProvider, error) {
	rs := make([]provider.IProvider, 0, len(ps))
	for _, p := range ps {
		v, ok := pm[p]
		if !ok {
			return nil, fmt.Errorf("provider:%s not found", p)
		}
		rs = append(rs, v)
	}
	return rs, nil
}

func buildRefresher(refresherConfigList []config.RefreshCongfig, pm map[string]provider.IProvider) ([]*refresher.Refresher, error) {
	rs := make([]*refresher.Refresher, 0, len(refresherConfigList))
	for _, item := range refresherConfigList {
		pvList, err := findProviderFromMap(item.Providers, pm)
		if err != nil {
			return nil, fmt.Errorf("find providers failed, err:%v", err)
		}
		if len(pvList) == 0 {
			return nil, fmt.Errorf("no provider found, name:%s", item.Name)
		}
		cli, err := client.New(
			client.WithAuth(item.CloudflareConfig.Key, item.CloudflareConfig.EMail),
		)
		if err != nil {
			return nil, fmt.Errorf("create cf client failed, err:%v", err)
		}
		refresherFn := adaptor.CFClientToRefresherFunc(
			cli,
			item.CloudflareConfig.ZoneName,
			item.CloudflareConfig.RecordType,
			item.CloudflareConfig.RecordName,
			item.CloudflareConfig.TTL,
			item.CloudflareConfig.Proxied,
		)
		rfr, err := refresher.New(
			refresher.WithIPProvider(provider.NewGroup(pvList...)),
			refresher.WithInterval(time.Duration(item.RefreshInterval)*time.Second),
			refresher.WithRefresherFunc(refresherFn),
			refresher.WithName(item.Name),
		)
		if err != nil {
			return nil, fmt.Errorf("create refresher failed, name:%s, err:%v", item.Name, err)
		}
		rs = append(rs, rfr)
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
	clients, err := buildRefresher(c.RefresherList, ps)
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
