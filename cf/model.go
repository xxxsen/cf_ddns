package cf

//{\"id\":\"$zone_identifier\",\"type\":\"A\",\"name\":\"$record_name\",\"content\":\"$ip\"}
type DDNSUpdateReq struct {
	ID      string `json:"id"`
	Type    string `json:"type"`
	Name    string `json:"record_name"`
	Content string `json:"content"`
}

type DDNSUpdateRsp struct {
	Success  bool          `json:"success"`
	Errors   []interface{} `json:"errors"`
	Messages []interface{} `json:"messages"`
}

type GetZoneIDTReq struct {
}

/*
{
	"success": true,
	"errors": [],
	"messages": [],
	"result": [
	  {
		"id": "023e105f4ecef8ad9ca31a8372d0c353",
		"name": "example.com",
	  }
	  ...
	]
}
*/
type GetZoneIDTRsp struct {
	Success  bool          `json:"success"`
	Errors   []interface{} `json:"errors"`
	Messages []interface{} `json:"messages"`
	Result   []struct {
		Id   string `json:"id"`
		Name string `json:"name"`
	} `json:"result"`
}

type GetRecordIDTReq struct {
}

/*
{
  "success": true,
  "errors": [],
  "messages": [],
  "result": {
    "id": "372e67954025e0ba6aaa6d586b9e0b59",
    "type": "A",
    "name": "example.com",
    "content": "198.51.100.4",
  }
}
*/
type GetRecordIDTRsp struct {
	Success  bool          `json:"success"`
	Errors   []interface{} `json:"errors"`
	Messages []interface{} `json:"messages"`
	Result   []struct {
		Id      string `json:"id"`
		Type    string `json:"type"`
		Name    string `json:"name"`
		Content string `json:"content"`
	} `json:"result"`
}
