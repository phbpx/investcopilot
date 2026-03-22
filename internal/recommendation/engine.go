package recommendation

import (
	"sort"

	"github.com/phbpx/investcopilot/internal/portfolio"
	"github.com/phbpx/investcopilot/internal/risk"
)

// Action represents a single recommended action for a contribution.
type Action struct {
	Class  portfolio.AssetClass
	Amount float64 // suggested BRL amount for this class
	Reason string
}

// Report is the output of the Recommendation Engine.
type Report struct {
	Actions []Action
	Avoid   []portfolio.AssetClass
	Monitor []string // tickers to watch
}

// Generate produces contribution recommendations based on portfolio state and risk analysis.
func Generate(p *portfolio.Portfolio, riskReport *risk.Report, contribution float64) *Report {
	report := &Report{}

	type candidate struct {
		class     portfolio.AssetClass
		deviation float64 // negative = below target
	}

	var below []candidate
	for _, dev := range riskReport.ClassDeviations {
		if dev.Deviation < 0 {
			below = append(below, candidate{class: dev.Class, deviation: dev.Deviation})
		}
		if dev.Deviation > 5 {
			report.Avoid = append(report.Avoid, dev.Class)
		}
	}

	// sort most below target first
	sort.Slice(below, func(i, j int) bool {
		return below[i].deviation < below[j].deviation
	})

	// distribute contribution proportionally to the deficit
	totalDeficit := 0.0
	for _, b := range below {
		totalDeficit += -b.deviation
	}

	if totalDeficit > 0 && contribution > 0 {
		for _, b := range below {
			weight := (-b.deviation) / totalDeficit
			report.Actions = append(report.Actions, Action{
				Class:  b.class,
				Amount: contribution * weight,
				Reason: "abaixo do target",
			})
		}
	}

	// flag tickers with notable concentration
	for _, pos := range p.Positions {
		if pos.Allocation > 15 {
			report.Monitor = append(report.Monitor, pos.Ticker)
		}
	}

	return report
}
