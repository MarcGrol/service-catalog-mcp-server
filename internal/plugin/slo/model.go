package slo

import (
	"github.com/MarcGrol/service-catalog-mcp-server/internal/plugin/slo/repo"
)

// List wraps a list into a single object (because the API does not allow lists)
type List struct {
	SLOs []repo.SLO `json:"slos"`
}
