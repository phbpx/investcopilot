package benchmark

import (
	"time"

	"github.com/phbpx/investcopilot/internal/market"
)

// Result holds benchmark returns for multiple periods.
type Result struct {
	Name      string
	Return1M  *float64
	Return6M  *float64
	Return12M *float64
}

// IBOVHistoryFn fetches historical prices for IBOV (^BVSP).
type IBOVHistoryFn func(ticker string) ([]market.HistoricalPrice, error)

// Fetch retrieves benchmark results.
// ibovHistory is optional — pass nil to skip IBOV benchmark.
func Fetch(ibovHistory IBOVHistoryFn) []*Result {
	now := time.Now()
	results := []*Result{}

	// --- CDI (always available, BCB is free) ---
	cdi := &Result{Name: "CDI"}
	if r, err := cdiAccumulated(now.AddDate(0, -1, 0), now); err == nil {
		cdi.Return1M = &r
	}
	if r, err := cdiAccumulated(now.AddDate(0, -6, 0), now); err == nil {
		cdi.Return6M = &r
	}
	if r, err := cdiAccumulated(now.AddDate(-1, 0, 0), now); err == nil {
		cdi.Return12M = &r
	}
	results = append(results, cdi)

	// --- IBOV (requires brapi token) ---
	if ibovHistory != nil {
		ibov := &Result{Name: "IBOV"}
		history, err := ibovHistory("^BVSP")
		if err == nil {
			ibov.Return1M = priceReturn(history, now.AddDate(0, -1, 0), now)
			ibov.Return6M = priceReturn(history, now.AddDate(0, -6, 0), now)
			ibov.Return12M = priceReturn(history, now.AddDate(-1, 0, 0), now)
		}
		results = append(results, ibov)
	}

	return results
}

// priceReturn calculates the price return between two dates from a history slice.
func priceReturn(history []market.HistoricalPrice, from, to time.Time) *float64 {
	startPrice, ok := market.PriceAt(history, from)
	if !ok || startPrice == 0 {
		return nil
	}
	endPrice, ok := market.PriceAt(history, to)
	if !ok || endPrice == 0 {
		return nil
	}
	r := ((endPrice - startPrice) / startPrice) * 100
	return &r
}
