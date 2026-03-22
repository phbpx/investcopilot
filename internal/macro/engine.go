package macro

import (
	"fmt"
	"time"
)

// Trend represents the direction of a macro indicator.
type Trend string

const (
	TrendRising  Trend = "rising"
	TrendFalling Trend = "falling"
	TrendStable  Trend = "stable"
)

// Context holds current macro indicators and forward-looking expectations.
type Context struct {
	// Current indicators
	SelicMeta   float64
	SelicTrend  Trend
	IPCA12m     float64
	USDBRL      float64

	// Focus market expectations (end of current year)
	FocusIPCA  *float64
	FocusSelic *float64

	FetchedAt time.Time
	Errors    []string // non-fatal fetch errors
}

// Fetch retrieves all macro indicators from BCB (free, no auth required).
func Fetch() *Context {
	ctx := &Context{FetchedAt: time.Now()}
	year := ctx.FetchedAt.Year()

	selic, trend, err := selicTrend()
	if err != nil {
		ctx.Errors = append(ctx.Errors, fmt.Sprintf("Selic: %v", err))
	} else {
		ctx.SelicMeta = selic
		ctx.SelicTrend = trend
	}

	ipca, err := ipcaAccumulated12m()
	if err != nil {
		ctx.Errors = append(ctx.Errors, fmt.Sprintf("IPCA: %v", err))
	} else {
		ctx.IPCA12m = ipca
	}

	usd, err := latestFloat(seriesUSD)
	if err != nil {
		ctx.Errors = append(ctx.Errors, fmt.Sprintf("USD/BRL: %v", err))
	} else {
		ctx.USDBRL = usd
	}

	if v, err := fetchFocusExpectation("IPCA", year); err == nil {
		ctx.FocusIPCA = &v
	}
	if v, err := fetchFocusExpectation("Selic", year); err == nil {
		ctx.FocusSelic = &v
	}

	return ctx
}

// SelicTrendLabel returns a human-readable trend label.
func (c *Context) SelicTrendLabel() string {
	switch c.SelicTrend {
	case TrendRising:
		return "em alta"
	case TrendFalling:
		return "em queda"
	default:
		return "estável"
	}
}
