package report

import (
	"fmt"
	"html/template"
	"io"
	"math"
	"strings"
	"time"

	"github.com/phbpx/investcopilot/internal/analytics"
	"github.com/phbpx/investcopilot/internal/benchmark"
	"github.com/phbpx/investcopilot/internal/macro"
	"github.com/phbpx/investcopilot/internal/playbook"
	"github.com/phbpx/investcopilot/internal/portfolio"
	"github.com/phbpx/investcopilot/internal/recommendation"
	"github.com/phbpx/investcopilot/internal/risk"
)

// Data holds all the information needed to render the report.
type Data struct {
	GeneratedAt    time.Time
	Portfolio      *portfolio.Portfolio
	Config         *portfolio.Config
	Risk           *risk.Report
	Performance    *analytics.Performance
	Benchmarks     []*benchmark.Result
	Macro          *macro.Context
	Playbook       *playbook.Report
	Recommendation *recommendation.Report
	Contribution   float64
}

// Render writes a self-contained HTML report to w.
func Render(w io.Writer, d *Data) error {
	tmpl, err := template.New("report").Funcs(funcMap(d)).Parse(htmlTemplate)
	if err != nil {
		return fmt.Errorf("parsing template: %w", err)
	}
	if err := tmpl.Execute(w, d); err != nil {
		return fmt.Errorf("executing template: %w", err)
	}
	return nil
}

func funcMap(d *Data) template.FuncMap {
	return template.FuncMap{
		"brl": func(v float64) string {
			return formatBRL(v)
		},
		"pct": func(v float64) string {
			return fmt.Sprintf("%.2f%%", v)
		},
		"pctSigned": func(v float64) string {
			if v >= 0 {
				return fmt.Sprintf("+%.2f%%", v)
			}
			return fmt.Sprintf("%.2f%%", v)
		},
		"ppSigned": func(v float64) string {
			if v >= 0 {
				return fmt.Sprintf("+%.1fpp", v)
			}
			return fmt.Sprintf("%.1fpp", v)
		},
		"returnStr": func(v *float64) string {
			if v == nil {
				return "—"
			}
			if *v >= 0 {
				return fmt.Sprintf("+%.2f%%", *v)
			}
			return fmt.Sprintf("%.2f%%", *v)
		},
		"returnClass": func(v *float64) string {
			if v == nil {
				return "neutral"
			}
			if *v >= 0 {
				return "positive"
			}
			return "negative"
		},
		"deviationClass": func(dev float64) string {
			switch {
			case dev > 5:
				return "above"
			case dev < -5:
				return "below"
			default:
				return "ok"
			}
		},
		"signalClass": func(level playbook.SignalLevel) string {
			switch level {
			case playbook.LevelOpportunity:
				return "opportunity"
			case playbook.LevelCaution:
				return "caution"
			default:
				return "info"
			}
		},
		"alertClass": func(level risk.AlertLevel) string {
			if level == risk.LevelWarn {
				return "warn"
			}
			return "info"
		},
		"allocationBar": func(current, target float64) template.HTML {
			return template.HTML(renderAllocationBar(current, target))
		},
		"date": func(t time.Time) string {
			return t.Format("02/01/2006")
		},
		"selicTrendIcon": func(t macro.Trend) string {
			switch t {
			case macro.TrendRising:
				return "↑"
			case macro.TrendFalling:
				return "↓"
			default:
				return "→"
			}
		},
		"upper": strings.ToUpper,
		"string": func(v interface{}) string {
			return fmt.Sprintf("%s", v)
		},
		"deref": func(v *float64) float64 {
			if v == nil {
				return 0
			}
			return *v
		},
	}
}

func renderAllocationBar(current, target float64) string {
	const width = 300
	const height = 20

	currentW := math.Min(current/100*width, width)
	targetX := math.Min(target/100*width, width)

	barColor := "#22c55e" // green
	if current-target > 10 {
		barColor = "#f59e0b" // amber
	} else if current-target < -10 {
		barColor = "#3b82f6" // blue
	}

	return fmt.Sprintf(
		`<svg width="%d" height="%d" xmlns="http://www.w3.org/2000/svg">
			<rect x="0" y="4" width="%d" height="%d" fill="#e5e7eb" rx="3"/>
			<rect x="0" y="4" width="%.1f" height="%d" fill="%s" rx="3"/>
			<line x1="%.1f" y1="0" x2="%.1f" y2="%d" stroke="#374151" stroke-width="2" stroke-dasharray="3,2"/>
		</svg>`,
		width, height,
		width, height-8,
		currentW, height-8, barColor,
		targetX, targetX, height,
	)
}

func formatBRL(v float64) string {
	negative := v < 0
	if negative {
		v = -v
	}
	s := fmt.Sprintf("%.2f", v)
	parts := strings.Split(s, ".")
	intPart := parts[0]
	result := ""
	for i, c := range intPart {
		if i > 0 && (len(intPart)-i)%3 == 0 {
			result += "."
		}
		result += string(c)
	}
	formatted := result + "," + parts[1]
	if negative {
		return "-" + formatted
	}
	return formatted
}
