package analytics

import (
	"time"

	"github.com/phbpx/investcopilot/internal/portfolio"
)

// Performance holds portfolio performance metrics.
type Performance struct {
	TotalInvested  float64
	CurrentValue   float64
	TotalIncome    float64 // accumulated proventos (dividends, FII rendimentos, JCP)
	AbsoluteReturn float64 // (current_value + income) - invested
	PercentReturn  float64 // absolute_return / invested

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
func Calculate(txs []portfolio.Transaction, currentValue float64, totalIncome float64, histPrices HistoricalPriceFn) *Performance {
	p := &Performance{
		CurrentValue: currentValue,
		TotalIncome:  totalIncome,
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

	// total return includes income received (dividends, FII rendimentos, JCP)
	p.AbsoluteReturn = (currentValue + totalIncome) - p.TotalInvested
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
// Includes income received during the period in the numerator.
func periodReturn(
	txs []portfolio.Transaction,
	currentValue float64,
	histPrices HistoricalPriceFn,
	start, end time.Time,
) *float64 {
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

	startValue := 0.0
	for ticker, h := range holdings {
		if h.qty <= 0 {
			continue
		}
		price, ok := histPrices(ticker, start)
		if !ok {
			return nil
		}
		startValue += h.qty * price
	}

	if startValue <= 0 {
		return nil
	}

	// net contributions and income received during the period
	netContributions := 0.0
	periodIncome := 0.0
	for _, tx := range txs {
		if tx.Date.Before(start) || tx.Date.After(end) {
			continue
		}
		switch tx.Type {
		case portfolio.Buy:
			netContributions += tx.Quantity*tx.Price + tx.Fees
		case portfolio.Sell:
			netContributions -= tx.Quantity*tx.Price - tx.Fees
		case portfolio.Income:
			periodIncome += tx.Price - tx.Fees
		}
	}

	// modified Dietz: income is treated as a cash inflow at the end
	// r = (end_value + income - start_value - contributions) / (start_value + contributions/2)
	denominator := startValue + netContributions/2
	if denominator <= 0 {
		return nil
	}

	ret := ((currentValue + periodIncome - startValue - netContributions) / denominator) * 100
	return &ret
}
