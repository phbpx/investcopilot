package playbook

import (
	"fmt"

	"github.com/phbpx/investcopilot/internal/macro"
	"github.com/phbpx/investcopilot/internal/portfolio"
	"github.com/phbpx/investcopilot/internal/risk"
)

// SignalLevel indicates the nature of a playbook signal.
type SignalLevel string

const (
	LevelOpportunity SignalLevel = "OPPORTUNITY"
	LevelCaution     SignalLevel = "CAUTION"
	LevelInfo        SignalLevel = "INFO"
)

// Signal is a single actionable insight from the playbook.
type Signal struct {
	Level   SignalLevel
	Class   portfolio.AssetClass // affected class, empty if general
	Message string
	Reason  string
}

// Report is the output of the Playbook Engine.
type Report struct {
	Signals []Signal
}

// Analyze combines macro context + portfolio state to generate signals.
func Analyze(macroCtx *macro.Context, p *portfolio.Portfolio, riskReport *risk.Report) *Report {
	r := &Report{}

	r.applyMacroRules(macroCtx)
	r.applyPortfolioRules(p, riskReport, macroCtx)

	return r
}

func (r *Report) add(level SignalLevel, class portfolio.AssetClass, message, reason string) {
	r.Signals = append(r.Signals, Signal{
		Level:   level,
		Class:   class,
		Message: message,
		Reason:  reason,
	})
}

// --- Macro rules ---

func (r *Report) applyMacroRules(m *macro.Context) {
	if m == nil || m.SelicMeta == 0 {
		return
	}

	// Selic high and stable/rising → renda fixa atrativa
	if m.SelicMeta >= 12 && m.SelicTrend != macro.TrendFalling {
		r.add(LevelOpportunity, portfolio.ClassRendaFixa,
			fmt.Sprintf("Renda fixa atrativa (Selic %.2f%% a.a., %s)", m.SelicMeta, m.SelicTrendLabel()),
			"Taxa real elevada favorece alocação em renda fixa",
		)
	}

	// Selic falling → equities tend to benefit
	if m.SelicTrend == macro.TrendFalling {
		r.add(LevelOpportunity, portfolio.ClassEquitiesBR,
			fmt.Sprintf("Ciclo de queda da Selic (%.2f%% a.a.)", m.SelicMeta),
			"Queda de juros historicamente beneficia bolsa e FIIs",
		)
		r.add(LevelOpportunity, portfolio.ClassRealEstate,
			"FIIs tendem a se valorizar em ciclo de queda de juros",
			"Custo de oportunidade reduz, múltiplos expandem",
		)
	}

	// IPCA high → inflation protection
	if m.IPCA12m >= 6 {
		r.add(LevelCaution, "",
			fmt.Sprintf("Inflação elevada: IPCA 12m em %.2f%%", m.IPCA12m),
			"Checar exposição a ativos indexados ao IPCA (Tesouro IPCA+, FIIs de tijolo)",
		)
	}

	// Focus IPCA above target (4.5% upper band)
	if m.FocusIPCA != nil && *m.FocusIPCA > 4.5 {
		r.add(LevelCaution, "",
			fmt.Sprintf("Expectativa de IPCA acima da meta: %.2f%% (Focus)", *m.FocusIPCA),
			"Mercado projeta inflação fora do centro da meta — revisar proteção inflacionária",
		)
	}

	// Focus Selic above current → market expects more hikes
	if m.FocusSelic != nil && *m.FocusSelic > m.SelicMeta {
		r.add(LevelCaution, portfolio.ClassEquitiesBR,
			fmt.Sprintf("Mercado espera Selic mais alta: %.2f%% ao final do ano (Focus)", *m.FocusSelic),
			"Expectativa de aperto monetário adicional pesa sobre bolsa",
		)
	}

	// Focus Selic below current → market expects cuts
	if m.FocusSelic != nil && *m.FocusSelic < m.SelicMeta-0.5 {
		r.add(LevelOpportunity, portfolio.ClassEquitiesBR,
			fmt.Sprintf("Mercado espera cortes de juros: Selic projetada em %.2f%% (Focus)", *m.FocusSelic),
			"Expectativa de afrouxamento monetário favorece risco",
		)
	}

	// USD/BRL high
	if m.USDBRL >= 5.5 {
		r.add(LevelInfo, portfolio.ClassEquitiesGlobal,
			fmt.Sprintf("Câmbio elevado: USD/BRL %.2f", m.USDBRL),
			"Exposição internacional em BRL beneficiada — avaliar se já está bem posicionado",
		)
	}
}

// --- Portfolio rules ---

func (r *Report) applyPortfolioRules(p *portfolio.Portfolio, riskReport *risk.Report, m *macro.Context) {
	for _, dev := range riskReport.ClassDeviations {
		// class significantly above target AND macro reinforces caution
		if dev.Deviation > 10 {
			r.add(LevelCaution, dev.Class,
				fmt.Sprintf("%s: %.1fpp acima do target — evitar novos aportes nessa classe", dev.Class, dev.Deviation),
				"Rebalancear progressivamente direcionando aportes para classes deficitárias",
			)
		}

		// class significantly below target AND macro creates opportunity
		if dev.Deviation < -10 {
			msg := fmt.Sprintf("%s: %.1fpp abaixo do target — priorizar nos próximos aportes", dev.Class, -dev.Deviation)
			reason := "Desvio expressivo do target de alocação"

			// reinforce if macro aligns
			if m != nil && dev.Class == portfolio.ClassRendaFixa && m.SelicMeta >= 12 {
				reason = "Desvio expressivo + Selic elevada: janela de entrada atrativa"
			}

			r.add(LevelOpportunity, dev.Class, msg, reason)
		}
	}
}
