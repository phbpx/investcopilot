package market

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type HistoricalPrice struct {
	Date  time.Time
	Close float64
}

type historicalResponse struct {
	Results []struct {
		Symbol               string `json:"symbol"`
		HistoricalDataPrice  []struct {
			Date  int64   `json:"date"` // unix timestamp
			Close float64 `json:"close"`
		} `json:"historicalDataPrice"`
	} `json:"results"`
}

// GetHistoricalPrices fetches 1-year daily historical prices for a ticker.
// Returns a slice sorted by date ascending.
func (c *Client) GetHistoricalPrices(ticker string) ([]HistoricalPrice, error) {
	url := fmt.Sprintf("%s%s?interval=1d&range=1y&fundamental=false", brapiBaseURL, ticker)
	if c.token != "" {
		url += "&token=" + c.token
	}

	resp, err := c.http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("fetching historical prices for %s: %w", ticker, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized {
		return nil, fmt.Errorf("brapi 401: configure market.brapi_token no config.yaml")
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("brapi returned status %d for %s", resp.StatusCode, ticker)
	}

	var result historicalResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("decoding response for %s: %w", ticker, err)
	}

	if len(result.Results) == 0 || len(result.Results[0].HistoricalDataPrice) == 0 {
		return nil, fmt.Errorf("no historical data for %s", ticker)
	}

	raw := result.Results[0].HistoricalDataPrice
	prices := make([]HistoricalPrice, 0, len(raw))
	for _, p := range raw {
		prices = append(prices, HistoricalPrice{
			Date:  time.Unix(p.Date, 0).UTC(),
			Close: p.Close,
		})
	}

	return prices, nil
}

// PriceAt returns the closest available price on or before the given date.
func PriceAt(history []HistoricalPrice, date time.Time) (float64, bool) {
	var best *HistoricalPrice
	for i := range history {
		h := &history[i]
		if !h.Date.After(date) {
			if best == nil || h.Date.After(best.Date) {
				best = h
			}
		}
	}
	if best == nil {
		return 0, false
	}
	return best.Close, true
}
