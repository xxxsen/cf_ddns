package provider

type CFCDNTraceProviderConfig struct {
	URL string `json:"url"`
}

type RawTextProviderConfig struct {
	URL string `json:"url"`
}

type JsonProviderConfig struct {
	URL  string `json:"url"`
	Path string `json:"path"`
}

type MachineConfig struct {
	Type         string `json:"type"`      //ipv4 or ipv6
	Interface    string `json:"interface"` //example: eth0
	AllowPrivate bool   `json:"allow_private"`
}
