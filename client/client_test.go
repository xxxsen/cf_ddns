package client

import (
	"context"
	"encoding/json"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

type testData struct {
	Key        string `json:"key"`
	Mail       string `json:"mail"`
	ZoneName   string `json:"zone_name"`
	RecordName string `json:"record_name"`
}

func mustGetConfig() *testData {
	raw, err := os.ReadFile("/tmp/test_data.json")
	if err != nil {
		panic(err)
	}
	data := &testData{}
	if err := json.Unmarshal(raw, data); err != nil {
		panic(err)
	}
	return data
}

func mustNewClient(key, mail string) *Client {
	client, err := New(WithAuth(key, mail))
	if err != nil {
		panic(err)
	}
	return client
}

func TestUpdateIP(t *testing.T) {
	cfg := mustGetConfig()
	client := mustNewClient(cfg.Key, cfg.Mail)
	ctx := context.Background()
	zoneRsp, err := client.GetZoneIdentifier(ctx, &GetZoneIdentifierRequest{
		ZoneName: cfg.ZoneName,
	})
	assert.NoError(t, err)
	t.Logf("%s", zoneRsp.Identifier)
	assert.True(t, zoneRsp.Exist)
	recRsp, err := client.GetRecordIdentifier(ctx, &GetRecordIdentifierRequest{
		ZoneIdentify: zoneRsp.Identifier,
		RecordName:   cfg.RecordName,
	})
	assert.NoError(t, err)
	assert.True(t, recRsp.Exist)
	t.Logf("%s", recRsp.Identifier)
	_, err = client.SetRecordIP(ctx, &SetRecordIPRequest{
		ZoneIdentify:   zoneRsp.Identifier,
		RecordIdentify: recRsp.Identifier,
		RecordType:     "A",
		RecordName:     cfg.RecordName,
		IP:             "5.6.7.8",
	})
	assert.NoError(t, err)
}
