package risk

import (
	"fmt"
	"sort"

	"github.com/phbpx/investcopilot/internal/portfolio"
)

// AlertLevel indicates the severity of a risk alert.
type AlertLevel string

const (
	LevelWarn AlertLevel = "WARN"
	LevelInfo AlertLevel = "INFO"
)

// Alert represents a risk signal.
type Alert struct {
	Level   AlertLevel
	Message string
}

// ClassDeviation represents the deviation of a class from its target.
type ClassDeviation struct {
	Class     portfolio.AssetClass
	Current   float64 // %
	Target    float64 // %
	Deviation float64 // current - target (pp)
}

// Report is the output of the Risk Engine.
type Report struct {
	Alerts          []Alert
	ClassDeviations []ClassDeviation
	Top3Weight      float64
}

// Analyze evaluates the portfolio against config rules and target allocation.
func Analyze(p *portfolio.Portfolio, cfg *portfolio.Config) *Report {
	report := &Report{}

	// --- class deviations ---
	for class, target := range cfg.TargetAllocation {
		current := p.ClassAllocation[class]
		dev := current - target
		report.ClassDeviations = append(report.ClassDeviations, ClassDeviation{
			Class:     class,
			Current:   current,
			Target:    target,
			Deviation: dev,
		})

		if dev > 10 {
			report.Alerts = append(report.Alerts, Alert{
				Level:   LevelWarn,
				Message: fmt.Sprintf("Classe %s %.1f%% acima do target (atual: %.1f%%, target: %.1f%%)", class, dev, current, target),
			})
		} else if dev < -10 {
			report.Alerts = append(report.Alerts, Alert{
				Level:   LevelInfo,
				Message: fmt.Sprintf("Classe %s %.1f%% abaixo do target (atual: %.1f%%, target: %.1f%%)", class, -dev, current, target),
			})
		}
	}

	// sort deviations by deviation descending (most above target first)
	sort.Slice(report.ClassDeviations, func(i, j int) bool {
		return report.ClassDeviations[i].Deviation > report.ClassDeviations[j].Deviation
	})

	// --- single asset concentration ---
	for _, pos := range p.Positions {
		if pos.Allocation > cfg.RiskRules.MaxSingleAsset {
			report.Alerts = append(report.Alerts, Alert{
				Level:   LevelWarn,
				Message: fmt.Sprintf("%s representa %.1f%% da carteira (limite: %.0f%%)", pos.Ticker, pos.Allocation, cfg.RiskRules.MaxSingleAsset),
			})
		}
	}

	// --- top 3 concentration ---
	top3 := 0.0
	for i, pos := range p.Positions {
		if i >= 3 {
			break
		}
		top3 += pos.Allocation
	}
	report.Top3Weight = top3
	if top3 > cfg.RiskRules.MaxTop3 {
		report.Alerts = append(report.Alerts, Alert{
			Level:   LevelWarn,
			Message: fmt.Sprintf("Top 3 ativos representam %.1f%% da carteira (limite: %.0f%%)", top3, cfg.RiskRules.MaxTop3),
		})
	}

	return report
}
