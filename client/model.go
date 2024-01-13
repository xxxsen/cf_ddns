package client

import "fmt"

type iErrorable interface {
	ConvertToError() error
}

type BaseResponse struct {
	Success  bool          `json:"success"`
	Errors   []interface{} `json:"errors"`
	Messages []interface{} `json:"messages"`
}

func (r *BaseResponse) ConvertToError() error {
	if !r.Success {
		return fmt.Errorf("response failed, err:%+v, msg:%+v", r.Errors, r.Messages)
	}
	return nil
}

//

type SetRecordIPRequest struct {
	ZoneIdentify   string
	RecordIdentify string
	RecordType     string
	RecordName     string
	IP             string
	TTL            int
	Proxied        bool
}

type SetRecordIPResponse struct {
}

type rpcSetRecordIPRequest struct {
	Type    string `json:"type"`
	Name    string `json:"name"`
	Content string `json:"content"`
	TTL     int    `json:"ttl"`
	Proxied bool   `json:"proxied"`
}

type rpcSetRecordIPResponse struct {
	BaseResponse
}

//

type GetZoneIdentifierRequest struct {
	ZoneName string
}

type rpcGetZoneIdentifierResponse struct {
	BaseResponse
	Result []struct {
		Id   string `json:"id"`
		Name string `json:"name"`
	} `json:"result"`
}

type GetZoneIdentifierResponse struct {
	Exist      bool
	Identifier string
}

type GetRecordIdentifierRequest struct {
	ZoneIdentify string
	RecordName   string
}

type rpcGetRecordIdentifierResponse struct {
	BaseResponse
	Result []struct {
		Id      string `json:"id"`
		Type    string `json:"type"`
		Name    string `json:"name"`
		Content string `json:"content"`
	} `json:"result"`
}

type GetRecordIdentifierResponse struct {
	Exist      bool
	Identifier string
}
