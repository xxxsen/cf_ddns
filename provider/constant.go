package provider

type IPType string

const (
	IPV4 IPType = "ipv4"
	IPV6 IPType = "ipv6"
)

const (
	ProviderCFCDNTrace = "cf_cdn"
	ProviderSendev     = "sendev"
	ProviderRawText    = "raw_text"
)
