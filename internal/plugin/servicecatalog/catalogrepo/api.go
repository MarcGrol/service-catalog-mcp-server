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
	ListMethods(ctx context.Context) ([]string, error)
	ListParticpantsOfFlow(ctx context.Context, id string) ([]string, bool, error)
	ListKinds(ctx context.Context) ([]string, error)
	ListModulesWithKind(ctx context.Context, id string) ([]string, bool, error)
}

// Module represents a software module in the catalog.
type Module struct {
	Version            string   `db:"version" json:"version"`
	ModuleID           string   `db:"module_id" json:"moduleID"`
	Name               string   `db:"name" json:"name"`
	Description        string   `db:"description" json:"description"`
	Spec               string   `db:"specification" json:"specification"`
	FileCount          int      `db:"file_count" json:"fileCount"`
	LineCount          int      `db:"line_count"  json:"lineCount"`
	ComplexityScore    float32  `json:",omitempty" json:"complexityScore,omitempty"`
	KindCount          *int     `db:"kind_count" json:"kindCount,omitempty"`
	TeamCount          *int     `db:"team_count" json:"teamCount,omitempty"`
	ExposedAPICount    *int     `db:"exposed_api_count" json:",omitempty"`
	ConsumedAPICount   *int     `db:"consumed_api_count" json:"exposedAPICount,omitempty"`
	DatabaseCount      *int     `db:"database_count" json:"databaseCount,omitempty"`
	JobCount           *int     `db:"job_count" json:"jobCount,omitempty"`
	FlowCount          *int     `db:"flow_count" json:"flowCount,omitempty"`
	ApplicationKinds   []string `db:"-" json:"applicationKinds,omitempty"`
	Teams              []string `db:"-" json:"teams,omitempty"`
	Flows              []string `db:"-" json:"flows,omitempty"`
	ExposedInterfaces  []string `db:"-" json:"exposedInterfaces,omitempty"`
	ConsumedInterfaces []string `db:"-" json:"consumedInterfaces,omitempty"`
	Jobs               []string `db:"-" json:"jobs,omitempty"`
	Databases          []string `db:"-" json:"databases,omitempty"`
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
	ModuleID      string   `db:"module_id" json:"moduleID,omitempty"`
	InterfaceID   string   `db:"interface_id" json:"interfaceID,omitempty"`
	Description   string   `db:"description" json:"description,omitempty"`
	Kind          string   `db:"kind" json:"kind,omitempty"`
	OpenAPISpecs  *string  `db:"openapi_specification" json:"-"` // API can not deal with null returned for string
	RPLSpecs      *string  `db:"rpl_specification" json:"-"`     // API can not deal with null returned for string
	MethodCount   int      `db:"method_count" json:"methodCount,omitempty"`
	Methods       []string `db:"-" json:"methods,omitempty"`
	MethodBasedID string   `db:"method_based_interface_id" json:"methodBasedID,omitempty"`
}
