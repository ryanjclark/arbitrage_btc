package config

// ForexConfig holds configuration details for a forex exchange connection.
type ForexConfig struct {
	HostURL *string
	Scheme  *string
	Path    *string
	Symbols []string
	APIKey  *string
}

// NewForexConfig hydrates a new config and returns a pointer to it.
func NewForexConfig(hostURL *string, scheme *string, path *string, symbols []string, apiKey *string, needTicker bool) *ForexConfig {
	return &ForexConfig{
		HostURL: hostURL,
		Scheme:  scheme,
		Path:    path,
		Symbols: symbols,
		APIKey:  apiKey,
	}
}
