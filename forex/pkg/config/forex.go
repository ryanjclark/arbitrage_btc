package config

import "net/url"

// ForexConfig ....
type ForexConfig struct {
	HostURL string
	Scheme  string
	Path    string
	Symbols []string
	APIKey  string
}

// InitURL ....
func (f *ForexConfig) InitURL() *url.URL {
	return &url.URL{
		Scheme: f.Scheme,
		Host:   f.HostURL,
		Path:   f.Path,
	}
}

// NewForexConfig ....
func NewForexConfig(hostURL string, scheme string, path string, symbols []string, apiKey string) *ForexConfig {
	return &ForexConfig{
		HostURL: hostURL,
		Scheme:  scheme,
		Path:    path,
		Symbols: symbols,
		APIKey:  apiKey,
	}
}
