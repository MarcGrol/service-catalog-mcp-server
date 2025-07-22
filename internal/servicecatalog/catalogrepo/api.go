package catalogrepo

import (
	"context"
	"encoding/json"
	"fmt"
)

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
	GroupInterfaces(ctx context.Context) (map[string][]Interface, error)
	GetInterfaceOnID(ctx context.Context, id string) (Interface, bool, error)
	ListInterfaceConsumers(ctx context.Context, id string) ([]string, bool, error)
	ListDatabaseConsumers(ctx context.Context, id string) ([]string, bool, error)
}

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
	ExposedApiCount    *int     `db:"exposed_api_count" json:",omitempty" yaml:",omitempty"`
	ConsumedApiCount   *int     `db:"consumed_api_count" json:",omitempty" yaml:",omitempty"`
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

func (m Module) CalculateComplexityScore() float32 {
	complexityScore := ((float32(m.LineCount) * 0.25) +
		(valueOrZero(m.DatabaseCount) * 0.20) +
		(valueOrZero(m.TeamCount) * 0.15) +
		(valueOrZero(m.ExposedApiCount) * 0.15) +
		(valueOrZero(m.ConsumedApiCount) * 0.15) +
		(valueOrZero(m.JobCount) * 0.10) +
		(float32(m.FileCount) * 0.10) +
		(valueOrZero(m.FlowCount) * 0.05) +
		(valueOrZero(m.KindCount) * 0.05)) * 100

	return complexityScore
}

func valueOrZero(value *int) float32 {
	if value == nil {
		return 0
	}
	return float32(*value)
}

func (m Module) String() string {
	asJson, _ := json.Marshal(m)
	return fmt.Sprintf("%s\n", asJson)
}

type Interface struct {
	InterfaceID   string   `db:"interface_id" yaml:",omitempty"`
	Description   string   `db:"description" yaml:",omitempty"`
	Kind          string   `db:"kind" yaml:",omitempty"`
	Spec          string   `db:"specification" yaml:",omitempty"`
	MethodCount   int      `db:"method_count" yaml:",omitempty"`
	Methods       []string `db:"-" yaml:",omitempty"`
	MethodBasedID string   `db:"method_based_interface_id" yaml:",omitempty"`
}
