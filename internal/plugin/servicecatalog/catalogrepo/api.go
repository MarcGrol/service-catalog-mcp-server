package catalogrepo

import (
	"context"
	"encoding/json"
	"fmt"
)

// Cataloger defines the interface for interacting with the service catalog repository.
//
//go:generate go tool mockgen -source=api.go -destination=mock_cataloger.go -package=catalogrepo Cataloger
type Cataloger interface {
	Open(ctx context.Context) error
	Close(ctx context.Context) error
	ListDatabases(ctx context.Context) ([]string, error)
	ListTeams(ctx context.Context) ([]string, error)
	ListModules(ctx context.Context, keyword string) ([]Module, error)
	ListModulesByCompexity(ctx context.Context, limit int) ([]Module, error)
	ListModulesOfTeam(ctx context.Context, id string) ([]string, bool, error)
	GetModuleOnID(ctx context.Context, id string) (Module, bool, error)
	ListInterfaces(ctx context.Context, keyword string) ([]Interface, error)
	ListInterfacesByComplexity(ctx context.Context, limit int) ([]Interface, error)
	GetInterfaceOnID(ctx context.Context, id string) (Interface, bool, error)
	ListInterfaceConsumers(ctx context.Context, id string) ([]string, bool, error)
	ListDatabaseConsumers(ctx context.Context, id string) ([]string, bool, error)
	ListFlows(ctx context.Context) ([]string, error)
	ListParticpantsOfFlow(ctx context.Context, id string) ([]string, bool, error)
	ListKinds(ctx context.Context) ([]string, error)
	ListModulesWithKind(ctx context.Context, id string) ([]string, bool, error)
}

// Module represents a software module in the catalog.
type Module struct {
	Version            string   `db:"version"  json:",omitempty" yaml:",omitempty"`
	ModuleID           string   `db:"module_id"`
	Name               string   `db:"name"`
	Description        string   `db:"description"`
	Spec               string   `db:"specification"`
	FileCount          int      `db:"file_count"`
	LineCount          int      `db:"line_count"`
	ComplexityScore    float32  `json:",omitempty" yaml:",omitempty"`
	KindCount          *int     `db:"kind_count" json:",omitempty" yaml:",omitempty"`
	TeamCount          *int     `db:"team_count" json:",omitempty" yaml:",omitempty"`
	ExposedAPICount    *int     `db:"exposed_api_count" json:",omitempty" yaml:",omitempty"`
	ConsumedAPICount   *int     `db:"consumed_api_count" json:",omitempty" yaml:",omitempty"`
	DatabaseCount      *int     `db:"database_count" json:",omitempty" yaml:",omitempty"`
	JobCount           *int     `db:"job_count" json:",omitempty" yaml:",omitempty"`
	FlowCount          *int     `db:"flow_count" json:",omitempty" yaml:",omitempty"`
	ApplicationKinds   []string `db:"-" json:",omitempty" yaml:",omitempty"`
	Teams              []string `db:"-" json:",omitempty" yaml:",omitempty"`
	Flows              []string `db:"-" json:",omitempty" yaml:",omitempty"`
	ExposedInterfaces  []string `db:"-" json:",omitempty" yaml:",omitempty"`
	ConsumedInterfaces []string `db:"-" json:",omitempty" yaml:",omitempty"`
	Jobs               []string `db:"-" json:",omitempty" yaml:",omitempty"`
	Databases          []string `db:"-" json:",omitempty" yaml:",omitempty"`
}

const (
	lineCountWeight           float32 = 0.25
	databaseCountWeight       float32 = 0.20
	teamCountWeight           float32 = 0.15
	exposedAPICountWeight     float32 = 0.15
	consumedAPICountWeight    float32 = 0.15
	jobCountWeight            float32 = 0.10
	fileCountWeight           float32 = 0.10
	flowCountWeight           float32 = 0.05
	kindCountWeight           float32 = 0.05
	complexityScoreMultiplier float32 = 100
)

// CalculateComplexityScore calculates the complexity score for a module.
func (m Module) CalculateComplexityScore() float32 {
	complexityScore := ((float32(m.LineCount/1000) * lineCountWeight) +
		(valueOrZero(m.DatabaseCount) * databaseCountWeight) +
		(valueOrZero(m.TeamCount) * teamCountWeight) +
		(valueOrZero(m.ExposedAPICount) * exposedAPICountWeight) +
		(valueOrZero(m.ConsumedAPICount) * consumedAPICountWeight) +
		(valueOrZero(m.JobCount) * jobCountWeight) +
		(valueOrZero(m.FlowCount) * flowCountWeight) +
		(valueOrZero(m.KindCount) * kindCountWeight)) * complexityScoreMultiplier

	return complexityScore
}

func valueOrZero(value *int) float32 {
	if value == nil {
		return 0
	}
	return float32(*value)
}

func (m Module) String() string {
	asJSON, _ := json.Marshal(m)
	return fmt.Sprintf("%s\n", asJSON)
}

// Interface represents a web API in the catalog.
type Interface struct {
	ModuleID      string   `db:"module_id" yaml:",omitempty"`
	InterfaceID   string   `db:"interface_id" yaml:",omitempty"`
	Description   string   `db:"description" yaml:",omitempty"`
	Kind          string   `db:"kind" yaml:",omitempty"`
	OpenAPISpecs  *string  `db:"openapi_specification" yaml:",omitempty"`
	RPLSpecs      *string  `db:"rpl_specification" yaml:",omitempty"`
	MethodCount   int      `db:"method_count" yaml:",omitempty"`
	Methods       []string `db:"-" yaml:",omitempty"`
	MethodBasedID string   `db:"method_based_interface_id" yaml:",omitempty"`
}
