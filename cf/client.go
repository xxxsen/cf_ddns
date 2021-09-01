package cf

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"time"
)

const (
	ddnsUpdateURI    = "https://api.cloudflare.com/client/v4/zones/%s/dns_records/%s"
	ddnsZoneIDTURI   = "https://api.cloudflare.com/client/v4/zones?name=%s"
	ddnsRecordIDTURI = "https://api.cloudflare.com/client/v4/zones/%s/dns_records?name=%s"
)

type identify struct {
	Zone   string
	Record string
}

type Client struct {
	c      *Config
	netcli *http.Client
	lastip string
	idt    *identify
}

func New(opts ...Option) *Client {
	c := &Config{}
	for _, opt := range opts {
		opt(c)
	}
	return &Client{c: c, netcli: &http.Client{
		Timeout: 10 * time.Second,
	}}
}

func (c *Client) getZoneIDTURI(zonename string) string {
	return fmt.Sprintf(ddnsZoneIDTURI, zonename)
}

func (c *Client) buildZone(zonename string) (string, error) {
	httpReq, err := http.NewRequest(http.MethodGet, c.getZoneIDTURI(zonename), nil)
	if err != nil {
		return "", err
	}
	c.buildHeader(httpReq)
	req := &GetZoneIDTReq{}
	rsp := &GetZoneIDTRsp{}
	if err := c.jsonCaller(httpReq, req, rsp); err != nil {
		return "", fmt.Errorf("build zone req fail, err:%v", err)
	}
	if !rsp.Success {
		return "", fmt.Errorf("logic fail, code:%v, msg:%v", rsp.Errors, rsp.Messages)
	}
	if len(rsp.Result) != 1 {
		return "", fmt.Errorf("invalid zone list")
	}
	return rsp.Result[0].Id, nil
}

func (c *Client) getRecordIDTURI(zoneidt string, recordname string) string {
	return fmt.Sprintf(ddnsRecordIDTURI, zoneidt, recordname)
}

func (c *Client) buildRecord(zoneidt string, recordname string) (string, error) {
	httpReq, err := http.NewRequest(http.MethodGet, c.getRecordIDTURI(zoneidt, recordname), nil)
	if err != nil {
		return "", err
	}
	c.buildHeader(httpReq)
	req := &GetRecordIDTReq{}
	rsp := &GetRecordIDTRsp{}
	if err := c.jsonCaller(httpReq, req, rsp); err != nil {
		return "", fmt.Errorf("build record req fail, err:%v", err)
	}
	if !rsp.Success {
		return "", fmt.Errorf("logic fail, code:%v, msg:%v", rsp.Errors, rsp.Messages)
	}
	if len(rsp.Result) != 1 {
		return "", fmt.Errorf("invalid record list")
	}
	return rsp.Result[0].Id, nil
}

func (c *Client) buildIdentify() error {
	zoneidt, err := c.buildZone(c.c.zoneName)
	if err != nil {
		return err
	}
	recidt, err := c.buildRecord(zoneidt, c.c.recName)
	if err != nil {
		return err
	}
	c.idt = &identify{
		Zone:   zoneidt,
		Record: recidt,
	}
	return nil
}

func (c *Client) getDDNSUpdateURI(idt *identify) string {
	return fmt.Sprintf(ddnsUpdateURI, idt.Zone, idt.Record)
}

func (c *Client) jsonCaller(httpReq *http.Request, req interface{}, rsp interface{}) error {
	raw, err := json.Marshal(req)
	if err != nil {
		return fmt.Errorf("encode fail, err:%v", err)
	}
	httpReq.Body = io.NopCloser(bytes.NewReader(raw))
	httpRsp, err := c.netcli.Do(httpReq)
	if err != nil {
		return err
	}
	defer httpRsp.Body.Close()
	if httpRsp.StatusCode != http.StatusOK {
		return fmt.Errorf("read invalid status code, code:%d", httpRsp.StatusCode)
	}
	data, err := ioutil.ReadAll(httpRsp.Body)
	if err != nil {
		return fmt.Errorf("read data fail, err:%v", err)
	}
	if err := json.Unmarshal(data, rsp); err != nil {
		return fmt.Errorf("decode fail, err:%v", err)
	}
	return nil
}

func (c *Client) buildHeader(req *http.Request) {
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Auth-Key", c.c.authKey)
	req.Header.Set("X-Auth-Email", c.c.authMail)
}

func (c *Client) RefreshDNS(ipstr string) error {
	if net.ParseIP(ipstr) == nil {
		return fmt.Errorf("invalid ip:%s", ipstr)
	}
	if c.idt == nil {
		if err := c.buildIdentify(); err != nil {
			return err
		}
	}

	httpReq, err := http.NewRequest(http.MethodPut, c.getDDNSUpdateURI(c.idt), nil)
	if err != nil {
		return err
	}
	c.buildHeader(httpReq)
	req := &DDNSUpdateReq{
		ID:      c.idt.Zone,
		Type:    "A",
		Name:    c.c.recName,
		Content: ipstr,
	}
	rsp := &DDNSUpdateRsp{}
	if err := c.jsonCaller(httpReq, req, rsp); err != nil {
		return fmt.Errorf("do update fail, err:%v", err)
	}
	if !rsp.Success {
		return fmt.Errorf("logic fail, code:%v, msg:%v", rsp.Errors, rsp.Messages)
	}
	return nil
}

func (c *Client) GetMyPublicIP() (string, error) {
	request, err := http.NewRequest(http.MethodGet, c.c.ipProvider, nil)
	if err != nil {
		return "", err
	}
	rsp, err := c.netcli.Do(request)
	if err != nil {
		return "", err
	}
	defer rsp.Body.Close()
	data, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func (c *Client) onCallback(newip string, oldip string, err error) {
	if newip == oldip && err == nil {
		return
	}
	if c.c.callback != nil {
		c.c.callback(newip, oldip, err)
	}
}

func (c *Client) refresh() {
	var newip string
	var err error
	var oldip = c.lastip
	defer func() {
		c.onCallback(newip, oldip, err)
	}()
	newip, err = c.GetMyPublicIP()
	if err != nil {
		log.Printf("get ip fail, provider:%s, err:%v", c.c.ipProvider, err)
		return
	}
	if newip == oldip && len(oldip) != 0 {
		return
	}
	err = c.RefreshDNS(newip)
	if err != nil {
		log.Printf("refresh ip failed, old ip:%s, new ip:%s, err:%v", oldip, newip, err)
		return
	}
	log.Printf("refresh ip succ, old ip:%s, new ip:%s", oldip, newip)
	c.lastip = newip

}

func (c *Client) StartRefresh(duration time.Duration) {
	if duration == 0 {
		duration = 10 * time.Second
	}
	for {
		c.refresh()
		time.Sleep(duration)
	}
}
