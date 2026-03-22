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

// MacroCache is the subset of cache.DB used by the macro engine.
type MacroCache interface {
	GetMacroFloat(key string, ttl time.Duration) (float64, bool)
	SetMacroFloat(key string, value float64) error
	GetMacroText(key string, ttl time.Duration) (string, bool)
	SetMacroText(key string, value string) error
}

const macroCacheTTL = 6 * time.Hour

// macro cache keys
const (
	keySelicMeta   = "selic_meta"
	keySelicTrend  = "selic_trend"
	keyIPCA12m     = "ipca_12m"
	keyUSDRBL      = "usd_brl"
	keyFocusIPCA   = "focus_ipca"
	keyFocusSelic  = "focus_selic"
)

// Context holds current macro indicators and forward-looking expectations.
type Context struct {
	SelicMeta  float64
	SelicTrend Trend
	IPCA12m    float64
	USDBRL     float64

	// Focus market expectations (end of current year)
	FocusIPCA  *float64
	FocusSelic *float64

	FetchedAt time.Time
	Errors    []string
}

// Fetch retrieves all macro indicators from BCB (free, no auth required).
// cache may be nil to skip caching.
func Fetch(cache MacroCache) *Context {
	ctx := &Context{FetchedAt: time.Now()}
	year := ctx.FetchedAt.Year()

	// Selic + trend
	if cache != nil {
		if v, ok := cache.GetMacroFloat(keySelicMeta, macroCacheTTL); ok {
			ctx.SelicMeta = v
		}
		if t, ok := cache.GetMacroText(keySelicTrend, macroCacheTTL); ok {
			ctx.SelicTrend = Trend(t)
		}
	}
	if ctx.SelicMeta == 0 {
		selic, trend, err := selicTrend()
		if err != nil {
			ctx.Errors = append(ctx.Errors, fmt.Sprintf("Selic: %v", err))
		} else {
			ctx.SelicMeta = selic
			ctx.SelicTrend = trend
			if cache != nil {
				_ = cache.SetMacroFloat(keySelicMeta, selic)
				_ = cache.SetMacroText(keySelicTrend, string(trend))
			}
		}
	}

	// IPCA 12m
	if cache != nil {
		if v, ok := cache.GetMacroFloat(keyIPCA12m, macroCacheTTL); ok {
			ctx.IPCA12m = v
		}
	}
	if ctx.IPCA12m == 0 {
		ipca, err := ipcaAccumulated12m()
		if err != nil {
			ctx.Errors = append(ctx.Errors, fmt.Sprintf("IPCA: %v", err))
		} else {
			ctx.IPCA12m = ipca
			if cache != nil {
				_ = cache.SetMacroFloat(keyIPCA12m, ipca)
			}
		}
	}

	// USD/BRL
	if cache != nil {
		if v, ok := cache.GetMacroFloat(keyUSDRBL, macroCacheTTL); ok {
			ctx.USDBRL = v
		}
	}
	if ctx.USDBRL == 0 {
		usd, err := latestFloat(seriesUSD)
		if err != nil {
			ctx.Errors = append(ctx.Errors, fmt.Sprintf("USD/BRL: %v", err))
		} else {
			ctx.USDBRL = usd
			if cache != nil {
				_ = cache.SetMacroFloat(keyUSDRBL, usd)
			}
		}
	}

	// Focus IPCA
	if cache != nil {
		if v, ok := cache.GetMacroFloat(keyFocusIPCA, macroCacheTTL); ok {
			ctx.FocusIPCA = &v
		}
	}
	if ctx.FocusIPCA == nil {
		if v, err := fetchFocusExpectation("IPCA", year); err == nil {
			ctx.FocusIPCA = &v
			if cache != nil {
				_ = cache.SetMacroFloat(keyFocusIPCA, v)
			}
		} else {
			ctx.Errors = append(ctx.Errors, fmt.Sprintf("Focus IPCA: %v", err))
		}
	}

	// Focus Selic
	if cache != nil {
		if v, ok := cache.GetMacroFloat(keyFocusSelic, macroCacheTTL); ok {
			ctx.FocusSelic = &v
		}
	}
	if ctx.FocusSelic == nil {
		if v, err := fetchFocusExpectation("Selic", year); err == nil {
			ctx.FocusSelic = &v
			if cache != nil {
				_ = cache.SetMacroFloat(keyFocusSelic, v)
			}
		} else {
			ctx.Errors = append(ctx.Errors, fmt.Sprintf("Focus Selic: %v", err))
		}
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
