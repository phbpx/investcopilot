package analytics

import (
	"time"

	"github.com/phbpx/investcopilot/internal/portfolio"
)

// Performance holds portfolio performance metrics.
type Performance struct {
	TotalInvested  float64
	CurrentValue   float64
	AbsoluteReturn float64
	PercentReturn  float64

	// Period returns — nil when historical data is unavailable.
	Return1M  *float64
	Return6M  *float64
	Return12M *float64
}

// HistoricalPriceFn returns the price of a ticker at a given date.
// Returns (price, true) if available, (0, false) if not.
type HistoricalPriceFn func(ticker string, date time.Time) (float64, bool)

// Calculate computes performance metrics for the portfolio.
// histPrices is optional — pass nil to skip period return calculations.
func Calculate(txs []portfolio.Transaction, currentValue float64, histPrices HistoricalPriceFn) *Performance {
	p := &Performance{
		CurrentValue: currentValue,
	}

	// net invested = sum(BUY cost) - sum(SELL proceeds)
	for _, tx := range txs {
		switch tx.Type {
		case portfolio.Buy:
			p.TotalInvested += tx.Quantity*tx.Price + tx.Fees
		case portfolio.Sell:
			p.TotalInvested -= tx.Quantity*tx.Price - tx.Fees
		}
	}

	p.AbsoluteReturn = currentValue - p.TotalInvested
	if p.TotalInvested > 0 {
		p.PercentReturn = (p.AbsoluteReturn / p.TotalInvested) * 100
	}

	if histPrices == nil {
		return p
	}

	now := time.Now()
	p.Return1M = periodReturn(txs, currentValue, histPrices, now.AddDate(0, -1, 0), now)
	p.Return6M = periodReturn(txs, currentValue, histPrices, now.AddDate(0, -6, 0), now)
	p.Return12M = periodReturn(txs, currentValue, histPrices, now.AddDate(-1, 0, 0), now)

	return p
}

// periodReturn calculates the money-weighted return for a period [start, end].
// It reconstructs the portfolio at `start`, values it with historical prices,
// then computes return accounting for contributions made during the period.
func periodReturn(
	txs []portfolio.Transaction,
	currentValue float64,
	histPrices HistoricalPriceFn,
	start, end time.Time,
) *float64 {
	// reconstruct holdings at start of period
	type holding struct {
		qty       float64
		costBasis float64
	}
	holdings := make(map[string]*holding)

	for _, tx := range txs {
		if tx.Date.After(start) {
			continue
		}
		h, ok := holdings[tx.Ticker]
		if !ok {
			h = &holding{}
			holdings[tx.Ticker] = h
		}
		switch tx.Type {
		case portfolio.Buy:
			h.costBasis += tx.Quantity*tx.Price + tx.Fees
			h.qty += tx.Quantity
		case portfolio.Sell:
			if h.qty > 0 {
				ratio := tx.Quantity / h.qty
				h.costBasis -= h.costBasis * ratio
			}
			h.qty -= tx.Quantity
		}
	}

	// value the portfolio at start using historical prices
	startValue := 0.0
	for ticker, h := range holdings {
		if h.qty <= 0 {
			continue
		}
		price, ok := histPrices(ticker, start)
		if !ok {
			return nil // can't compute without all prices
		}
		startValue += h.qty * price
	}

	if startValue <= 0 {
		return nil
	}

	// net contributions made during the period
	netContributions := 0.0
	for _, tx := range txs {
		if tx.Date.Before(start) || tx.Date.After(end) {
			continue
		}
		switch tx.Type {
		case portfolio.Buy:
			netContributions += tx.Quantity*tx.Price + tx.Fees
		case portfolio.Sell:
			netContributions -= tx.Quantity*tx.Price - tx.Fees
		}
	}

	// simple return adjusted for contributions
	// r = (end_value - start_value - contributions) / (start_value + contributions/2)
	denominator := startValue + netContributions/2
	if denominator <= 0 {
		return nil
	}

	ret := ((currentValue - startValue - netContributions) / denominator) * 100
	return &ret
}
