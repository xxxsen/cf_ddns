package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const (
	ddnsUpdateURI    = "https://api.cloudflare.com/client/v4/zones/%s/dns_records/%s"
	ddnsZoneIDTURI   = "https://api.cloudflare.com/client/v4/zones?name=%s"
	ddnsRecordIDTURI = "https://api.cloudflare.com/client/v4/zones/%s/dns_records?name=%s"
)

type Client struct {
	c      *config
	client *http.Client
}

func New(opts ...Option) (*Client, error) {
	c := &config{}
	for _, opt := range opts {
		opt(c)
	}
	if len(c.email) == 0 || len(c.key) == 0 {
		return nil, fmt.Errorf("invalid auth data")
	}
	cli := &http.Client{
		Transport: &http.Transport{
			MaxIdleConns:        10,
			MaxIdleConnsPerHost: 2,
			IdleConnTimeout:     120 * time.Second,
		},
		Timeout: 5 * time.Second,
	}
	return &Client{c: c, client: cli}, nil
}

func (c *Client) makeRequest(ctx context.Context, method string, uri string, request interface{}, response interface{}) error {
	var body io.Reader
	if request != nil {
		data, err := json.Marshal(request)
		if err != nil {
			return err
		}
		body = bytes.NewReader(data)
	}
	req, err := http.NewRequestWithContext(ctx, method, uri, body)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Auth-Key", c.c.key)
	req.Header.Set("X-Auth-Email", c.c.email)
	rsp, err := c.client.Do(req)
	if err != nil {
		return err
	}
	defer rsp.Body.Close()
	if rsp.StatusCode != http.StatusOK {
		return fmt.Errorf("invalid status code:%d", rsp.StatusCode)
	}
	raw, err := io.ReadAll(rsp.Body)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(raw, response); err != nil {
		return err
	}
	return nil
}

func (c *Client) makeRequestWrap(ctx context.Context, method string, uri string, req interface{}, rsp iErrorable) error {
	err := c.makeRequest(ctx, method, uri, req, rsp)
	if err != nil {
		return err
	}
	if err := rsp.ConvertToError(); err != nil {
		return err
	}
	return nil
}

func (c *Client) GetZoneIdentifier(ctx context.Context, req *GetZoneIdentifierRequest) (*GetZoneIdentifierResponse, error) {
	uri := fmt.Sprintf(ddnsZoneIDTURI, req.ZoneName)
	rpcRsp := &rpcGetZoneIdentifierResponse{}
	err := c.makeRequestWrap(ctx, http.MethodGet, uri, nil, rpcRsp)
	if err != nil {
		return nil, err
	}
	rsp := &GetZoneIdentifierResponse{Exist: false}
	if len(rpcRsp.Result) > 0 {
		rsp.Exist = true
		rsp.Identifier = rpcRsp.Result[0].Id
	}
	return rsp, nil

}

func (c *Client) GetRecordIdentifier(ctx context.Context, req *GetRecordIdentifierRequest) (*GetRecordIdentifierResponse, error) {
	uri := fmt.Sprintf(ddnsRecordIDTURI, req.ZoneIdentify, req.RecordName)
	rpcRsp := &rpcGetRecordIdentifierResponse{}
	if err := c.makeRequestWrap(ctx, http.MethodGet, uri, nil, rpcRsp); err != nil {
		return nil, err
	}
	rsp := &GetRecordIdentifierResponse{Exist: false}
	if len(rpcRsp.Result) > 0 {
		rsp.Exist = true
		rsp.Identifier = rpcRsp.Result[0].Id
	}
	return rsp, nil
}

func (c *Client) SetRecordIP(ctx context.Context, req *SetRecordIPRequest) (*SetRecordIPResponse, error) {
	uri := fmt.Sprintf(ddnsUpdateURI, req.ZoneIdentify, req.RecordIdentify)
	if req.TTL == 0 {
		req.TTL = 120
	}
	rpcReq := &rpcSetRecordIPRequest{
		Type:    req.RecordType,
		Name:    req.RecordName,
		Content: req.IP,
		TTL:     req.TTL,
		Proxied: req.Proxied,
	}
	rpcRsp := &rpcSetRecordIPResponse{}

	if err := c.makeRequestWrap(ctx, http.MethodPut, uri, rpcReq, rpcRsp); err != nil {
		return nil, err
	}
	return &SetRecordIPResponse{}, nil
}
