package main

import (
	"cf_ddns/adaptor"
	"cf_ddns/client"
	"cf_ddns/config"
	"cf_ddns/notifier"
	_ "cf_ddns/notifier/tgmsger"
	"cf_ddns/provider"
	"cf_ddns/refresher"
	"context"
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/xxxsen/common/logger"
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

func createNotifierMap(items []config.NotifierConfig) (map[string]notifier.INotifier, error) {
	rs := make(map[string]notifier.INotifier)
	for _, item := range items {
		nt, err := notifier.MakeNotifier(item.Type, item.Data)
		if err != nil {
			return nil, err
		}
		rs[item.Name] = nt
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

func buildRefresher(refresherConfigList []config.RefreshCongfig, pm map[string]provider.IProvider, ntsm map[string]notifier.INotifier) ([]*refresher.Refresher, error) {
	rs := make([]*refresher.Refresher, 0, len(refresherConfigList))
	for _, item := range refresherConfigList {
		if item.Disable {
			continue
		}
		pvList, err := findProviderFromMap(item.Providers, pm)
		if err != nil {
			return nil, fmt.Errorf("find providers failed, err:%v", err)
		}
		if len(pvList) == 0 {
			return nil, fmt.Errorf("no provider found, name:%s", item.Name)
		}
		nt, ok := ntsm[item.Notifier]
		if !ok && len(item.Notifier) == 0 {
			nt = notifier.NopNotifier
			ok = true
		}
		if !ok {
			return nil, fmt.Errorf("notifier:%s invalid", item.Notifier)
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
			refresher.WithCallback(nt.Notify),
			refresher.WithDomain(item.CloudflareConfig.RecordName),
		)
		if err != nil {
			return nil, fmt.Errorf("create refresher failed, name:%s, err:%v", item.Name, err)
		}
		rs = append(rs, rfr)
	}
	return rs, nil
}

func main() {
	flag.Parse()
	c, err := config.Parse(*cfg)
	if err != nil {
		log.Fatalf("parse config failed, err:%v", err)
	}
	logger := logger.Init(c.LogConfig.File, c.LogConfig.Level, int(c.LogConfig.FileCount), int(c.LogConfig.FileSize), int(c.LogConfig.KeepDays), c.LogConfig.Console)
	logger.Info("support providers", zap.Strings("providers", provider.List()))
	logger.Info("support notifiers", zap.Strings("notifiers", notifier.List()))
	ps, err := createProviderMap(c.ProviderList)
	if err != nil {
		logger.Fatal("create provider map failed", zap.Error(err))
	}
	nts, err := createNotifierMap(c.NotifierList)
	if err != nil {
		logger.Fatal("create notifier map failed", zap.Error(err))
	}
	clients, err := buildRefresher(c.RefresherList, ps, nts)
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
