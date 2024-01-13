package provider

type cfCDNTraceProviderConfig struct {
	URL string `json:"url"`
}

type rawTextProviderConfig struct {
	URL string `json:"url"`
}

type jsonProviderConfig struct {
	URL  string `json:"url"`
	Path string `json:"path"`
}
