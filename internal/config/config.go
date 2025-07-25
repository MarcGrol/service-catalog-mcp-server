package config

// Config holds the application's configuration.
type Config struct {
	UseSSE        bool
	UseStreamable bool
	Port          string
	BaseURL       string
	APIKey        string
	PluginConfigs map[string]string
}
