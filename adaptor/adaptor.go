package adaptor

import (
	"cf_ddns/client"
	"cf_ddns/refresher"
	"context"
	"fmt"
	"net"
	"strings"
)

func buildCreateTmpRecordRequest(zident string, typ string, name string) *client.CreateRecordRequest {
	ip := "0.0.0.0"
	if strings.EqualFold(typ, "aaaa") {
		ip = "::"
	}
	return &client.CreateRecordRequest{
		ZoneIdentify: zident,
		RecordType:   typ,
		RecordName:   name,
		IP:           ip,
		TTL:          60,
		Proxied:      false,
	}
}

func fetchZoneRecId(ctx context.Context, cli *client.Client, zone string, typ string, record string) (string, string, error) {
	zoneRsp, err := cli.GetZoneIdentifier(ctx, &client.GetZoneIdentifierRequest{
		ZoneName: zone,
	})
	if err != nil {
		return "", "", err
	}
	//
	recReq := &client.GetRecordIdentifierRequest{
		ZoneIdentify: zoneRsp.Identifier,
		RecordName:   record,
		RecordType:   typ,
	}
	recRsp, err := cli.GetRecordIdentifier(ctx, recReq)
	if err != nil && err != client.ErrIdentifierNotFound {
		return "", "", err
	}
	//如果记录不存在, 那么尝试创建一条新的记录并重新获取
	if err == client.ErrIdentifierNotFound {
		_, _ = cli.CreateRecord(ctx, buildCreateTmpRecordRequest(zoneRsp.Identifier, typ, record))
		recRsp, err = cli.GetRecordIdentifier(ctx, recReq)
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
			zoneid, recid, err = fetchZoneRecId(ctx, cli, zone, recordType, record)
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
