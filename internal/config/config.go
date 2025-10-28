package config

// Mode represents what mcp-services to run
type Mode string

const (
	// Both runs everything
	Both Mode = "both"
	// SLO runs the SLO  mcp-service
	SLO Mode = "slo"
	// ServiceCatalog runs the servoce-catalog related mcp-service
	ServiceCatalog Mode = "service-catalog"
)

// Config holds the application's configuration.
type Config struct {
	UseSSE        bool
	UseStreamable bool
	Port          string
	BaseURL       string
	APIKey        string
	Mode          Mode
	PluginConfigs map[string]string
}
