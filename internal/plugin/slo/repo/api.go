package repo

import context "context"

// SLORepo defines the interface for SLO repository operations.
//
//go:generate mockgen -source=api.go -destination=mock_repo.go -package=repo SLORepo
type SLORepo interface {
	Open(ctx context.Context) error
	Close(ctx context.Context) error
	ListSLOs(ctx context.Context) ([]SLO, error)
	GetSLOByID(ctx context.Context, id string) (SLO, bool, error)
	ListSLOsByTeam(ctx context.Context, id string) ([]SLO, bool, error)
	ListSLOsByApplication(ctx context.Context, id string) ([]SLO, bool, error)
}

// SLO represents the structure of the SLO data.
type SLO struct {
	UID                  string  `json:"uid" db:"UID"`
	CreatedAt            string  `json:"created_at" db:"CreatedAt"`
	LastModified         string  `json:"last_modified" db:"LastModified"`
	ModificationCount    int     `json:"modification_count" db:"ModificationCount"`
	Filename             string  `json:"filename" db:"Filename"`
	DisplayName          string  `json:"display_name" db:"DisplayName"`
	Team                 string  `json:"team" db:"Team"`
	Application          string  `json:"application" db:"Application"`
	Service              string  `json:"service" db:"Service"`
	Component            string  `json:"component" db:"Component"`
	Category             string  `json:"category" db:"Category"`
	RelativeThroughput   float64 `json:"relative_throughput" db:"RelativeThroughput"`
	PromQLQuery          string  `json:"promql_query" db:"PromQLQuery"`
	TargetSLO            float64 `json:"target_slo" db:"TargetSLO"`
	Duration             string  `json:"duration" db:"Duration"`
	SLI                  float64 `json:"sli" db:"SLI"`
	DashboardLinkCount   int     `json:"dashboard_link_count" db:"DashboardLinkCount"`
	AlertLinkCount       int     `json:"alert_link_count" db:"AlertLinkCount"`
	EmailChannelCount    int     `json:"email_channel_count" db:"EmailChannelCount"`
	ChatChannelCount     int     `json:"chat_channel_count" db:"ChatChannelCount"`
	IsEnriched           bool    `json:"is_enriched" db:"IsEnriched"`
	IsCritical           bool    `json:"is_critical" db:"IsCritical"`
	IsFrontdoor          bool    `json:"is_frontdoor" db:"IsFrontdoor"`
	IsOnlinePaymentsFlow bool    `json:"is_online_payments_flow" db:"IsOnlinePaymentsFlow"`
	IsIPPPaymentsFlow    bool    `json:"is_ipp_payments_flow" db:"IsIPPPaymentsFlow"`
	IsPayoutFlow         bool    `json:"is_payout_flow" db:"IsPayoutFlow"`
	IsReportingFlow      bool    `json:"is_reporting_flow" db:"IsReportingFlow"`
	IsOnboardingFlow     bool    `json:"is_onboarding_flow" db:"IsOnboardingFlow"`
	IsCustomerPortalFlow bool    `json:"is_customer_portal_flow" db:"IsCustomerPortalFlow"`
	CriticalFlows        string  `json:"critical_flows" db:"CriticalFlows"`
	OperationalReadiness float64 `json:"operational_readiness" db:"-"`
	BusinessCriticality  float64 `json:"business_criticality" db:"-"`
}

// CalculateOperationalReadinessMultiplier calculates a score based on operational support.
func (s SLO) calculateOperationalReadinessMultiplier() float64 {
	multiplier := 1.0

	if s.DashboardLinkCount > 0 {
		multiplier += 0.1
	}
	if s.AlertLinkCount > 0 {
		multiplier += 0.1
	}
	if s.EmailChannelCount > 0 || s.ChatChannelCount > 0 {
		multiplier += 0.05
	}
	if s.IsEnriched {
		multiplier += 0.1
	}

	return multiplier
}

// CalculateBusinessCriticalityMultiplier calculates a score based on business criticality.
func (s SLO) calculateBusinessCriticalityMultiplier() float64 {
	multiplier := 1.0

	if s.IsCritical {
		multiplier += 0.5
	}
	if s.IsFrontdoor {
		multiplier += 0.5
	}

	return multiplier * s.RelativeThroughput
}
