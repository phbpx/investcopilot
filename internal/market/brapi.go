package market

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

const brapiBaseURL = "https://brapi.dev/api/quote/"

// Client fetches current market prices from brapi.dev.
type Client struct {
	http  *http.Client
	token string
}

// NewClient creates a new market data client.
// token is optional — get a free token at brapi.dev.
func NewClient(token string) *Client {
	return &Client{
		http:  &http.Client{Timeout: 10 * time.Second},
		token: token,
	}
}

type brapiResponse struct {
	Results []struct {
		Symbol             string  `json:"symbol"`
		RegularMarketPrice float64 `json:"regularMarketPrice"`
	} `json:"results"`
}

// GetPrices fetches current prices for the given tickers.
// Returns a map of ticker → price.
// If the batch request fails, falls back to fetching one ticker at a time.
func (c *Client) GetPrices(tickers []string) (map[string]float64, error) {
	if len(tickers) == 0 {
		return map[string]float64{}, nil
	}

	prices, err := c.fetchBatch(tickers)
	if err == nil {
		return prices, nil
	}

	// batch failed — fall back to one-by-one to identify the problematic ticker
	prices = make(map[string]float64, len(tickers))
	for _, ticker := range tickers {
		p, err := c.fetchBatch([]string{ticker})
		if err != nil {
			return nil, fmt.Errorf("ticker %s: %w", ticker, err)
		}
		for k, v := range p {
			prices[k] = v
		}
	}
	return prices, nil
}

func (c *Client) fetchBatch(tickers []string) (map[string]float64, error) {
	url := brapiBaseURL + strings.Join(tickers, ",")
	if c.token != "" {
		url += "?token=" + c.token
	}

	resp, err := c.http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("fetching prices: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized {
		return nil, fmt.Errorf("brapi returned 401: configure um token em market.brapi_token no config.yaml (obtenha gratuitamente em brapi.dev)")
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("brapi returned status %d", resp.StatusCode)
	}

	var result brapiResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("decoding response: %w", err)
	}

	prices := make(map[string]float64, len(result.Results))
	for _, r := range result.Results {
		prices[r.Symbol] = r.RegularMarketPrice
	}
	return prices, nil
}
