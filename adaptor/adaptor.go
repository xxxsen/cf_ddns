package adaptor

import (
	"cf_ddns/client"
	"cf_ddns/notifier"
	"cf_ddns/refresher"
	"context"
	"fmt"
	"net"
	"strings"
)

func fetchZoneRecId(ctx context.Context, cli *client.Client, zone string, record string) (string, string, error) {
	zoneRsp, err := cli.GetZoneIdentifier(ctx, &client.GetZoneIdentifierRequest{
		ZoneName: zone,
	})
	if err != nil {
		return "", "", err
	}
	//
	recRsp, err := cli.GetRecordIdentifier(ctx, &client.GetRecordIdentifierRequest{
		ZoneIdentify: zoneRsp.Identifier,
		RecordName:   record,
	})
	if err != nil && err != client.ErrIdentifierNotFound {
		return "", "", err
	}
	//如果记录不存在, 那么尝试创建一条新的记录并重新获取
	if err == client.ErrIdentifierNotFound {
		_, _ = cli.CreateRecord(ctx, &client.CreateRecordRequest{
			ZoneIdentify: zoneRsp.Identifier,
			RecordType:   "A",
			RecordName:   record,
			IP:           "127.0.0.1",
			TTL:          60,
			Proxied:      false,
		})
		recRsp, err = cli.GetRecordIdentifier(ctx, &client.GetRecordIdentifierRequest{
			ZoneIdentify: zoneRsp.Identifier,
			RecordName:   record,
		})
	}
	if err != nil {
		return "", "", err
	}
	return zoneRsp.Identifier, recRsp.Identifier, nil
}

func precheckIP(typ string, rawip string) error {
	ip := net.ParseIP(rawip)
	if ip == nil {
		return fmt.Errorf("invalid ip:%s", rawip)
	}
	if (strings.EqualFold(typ, "aaaa") && !strings.Contains(rawip, ":")) ||
		(strings.EqualFold(typ, "a") && !strings.Contains(rawip, ".")) {
		return fmt.Errorf("ip:%s not match typ:%s", rawip, typ)
	}
	return nil
}

func CFClientToRefresherFunc(cli *client.Client, zone string, recordType string, record string, ttl int, proxied bool) refresher.RefresherFunc {
	var zoneid string
	var recid string

	return func(ctx context.Context, ip string) error {
		var err error
		if len(zoneid) == 0 || len(recid) == 0 {
			zoneid, recid, err = fetchZoneRecId(ctx, cli, zone, record)
			if err != nil {
				return err
			}
		}
		if err := precheckIP(recordType, ip); err != nil {
			return err
		}
		if _, err = cli.SetRecordIP(ctx, &client.SetRecordIPRequest{
			ZoneIdentify:   zoneid,
			RecordIdentify: recid,
			RecordType:     recordType,
			RecordName:     record,
			IP:             ip,
			TTL:            ttl,
			Proxied:        proxied,
		}); err != nil {
			return err
		}
		return nil
	}
}

func NotifierClientToRefreshCallback(cli notifier.INotifier) refresher.CallbackFunc {
	return func(ctx context.Context, name, domain, oldip, newip string) error {
		msg := fmt.Sprintf("[CF_DDNS: %s]: refresher: %s, refresh ip to: %s", domain, name, newip)
		if len(oldip) > 0 {
			msg += ", old ip: " + oldip
		}
		return cli.Notify(ctx, msg)
	}
}
