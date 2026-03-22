package portfolio

import (
	"fmt"
	"sort"
)

// Position represents a consolidated holding for a single ticker.
type Position struct {
	Ticker       string
	Class        AssetClass
	Sector       string
	Quantity     float64
	AvgPrice     float64 // average cost basis
	CurrentPrice float64
	CurrentValue float64
	TotalIncome  float64 // accumulated dividends / FII rendimentos / JCP (net of taxes)
	Allocation   float64 // % of total portfolio
}

// Portfolio is the consolidated view of all holdings.
type Portfolio struct {
	Positions       []*Position
	TotalValue      float64
	TotalIncome     float64
	ClassAllocation map[AssetClass]float64 // % per class
}

// PriceSource is a function that returns current prices for a set of tickers.
type PriceSource func(tickers []string) (map[string]float64, error)

// Build consolidates transactions into a Portfolio using the given price source.
func Build(txs []Transaction, cfg *Config, prices PriceSource) (*Portfolio, error) {
	type state struct {
		quantity    float64
		costBasis   float64 // total cost (quantity * avg price + fees)
		totalIncome float64 // accumulated proventos net of taxes
	}

	holdings := make(map[string]*state)

	for _, tx := range txs {
		s, ok := holdings[tx.Ticker]
		if !ok {
			s = &state{}
			holdings[tx.Ticker] = s
		}

		switch tx.Type {
		case Buy:
			s.costBasis += tx.Quantity*tx.Price + tx.Fees
			s.quantity += tx.Quantity
		case Sell:
			if tx.Quantity > s.quantity {
				return nil, fmt.Errorf("ticker %s: selling %.2f but only have %.2f", tx.Ticker, tx.Quantity, s.quantity)
			}
			ratio := tx.Quantity / s.quantity
			s.costBasis -= s.costBasis * ratio
			s.quantity -= tx.Quantity
		case Income:
			// price = total received, fees = IR withheld
			s.totalIncome += tx.Price - tx.Fees
		}
	}

	// collect tickers with non-zero positions
	tickers := make([]string, 0, len(holdings))
	for ticker, s := range holdings {
		if s.quantity > 0 {
			tickers = append(tickers, ticker)
		}
	}

	// fetch current prices
	currentPrices, err := prices(tickers)
	if err != nil {
		return nil, fmt.Errorf("fetching prices: %w", err)
	}

	// build positions
	positions := make([]*Position, 0, len(tickers))
	totalValue := 0.0

	for _, ticker := range tickers {
		s := holdings[ticker]
		price, ok := currentPrices[ticker]
		if !ok {
			return nil, fmt.Errorf("no price available for %s", ticker)
		}

		avgPrice := 0.0
		if s.quantity > 0 {
			avgPrice = s.costBasis / s.quantity
		}

		value := s.quantity * price
		totalValue += value

		assetCfg := cfg.Assets[ticker]

		positions = append(positions, &Position{
			Ticker:       ticker,
			Class:        assetCfg.Class,
			Sector:       assetCfg.Sector,
			Quantity:     s.quantity,
			AvgPrice:     avgPrice,
			CurrentPrice: price,
			CurrentValue: value,
			TotalIncome:  s.totalIncome,
		})
	}

	// calculate allocations and total income
	classAllocation := make(map[AssetClass]float64)
	totalIncome := 0.0
	for _, p := range positions {
		if totalValue > 0 {
			p.Allocation = (p.CurrentValue / totalValue) * 100
		}
		classAllocation[p.Class] += p.Allocation
		totalIncome += p.TotalIncome
	}

	// also count income from fully sold positions
	for ticker, s := range holdings {
		if s.quantity == 0 && s.totalIncome > 0 {
			// already sold — income still counts but position not shown
			_ = ticker
			totalIncome += s.totalIncome
		}
	}

	// sort positions by value descending
	sort.Slice(positions, func(i, j int) bool {
		return positions[i].CurrentValue > positions[j].CurrentValue
	})

	return &Portfolio{
		Positions:       positions,
		TotalValue:      totalValue,
		TotalIncome:     totalIncome,
		ClassAllocation: classAllocation,
	}, nil
}
