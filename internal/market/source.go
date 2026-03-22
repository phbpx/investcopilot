package market

import (
	"fmt"
	"time"
)

// PriceCache is the subset of cache.DB used by the price source.
type PriceCache interface {
	GetCurrentPrice(ticker string, ttl time.Duration) (float64, bool)
	SetCurrentPrice(ticker string, price float64) error
	GetHistoricalPrices(ticker string, ttl time.Duration) ([]HistoricalPrice, bool)
	SetHistoricalPrices(ticker string, prices []HistoricalPrice) error
}

const (
	CurrentPriceTTL  = 1 * time.Hour
	HistoricalPriceTTL = 24 * time.Hour
)

// NewPriceSource returns a PriceSource that resolves prices in order:
//  1. manual prices (config)
//  2. cache (SQLite)
//  3. brapi.dev API
//
// cache may be nil to skip caching.
func NewPriceSource(manual map[string]float64, token string, cache PriceCache) func(tickers []string) (map[string]float64, error) {
	client := NewClient(token)

	return func(tickers []string) (map[string]float64, error) {
		result := make(map[string]float64, len(tickers))

		var needFetch []string
		for _, t := range tickers {
			// 1. manual prices always win
			if price, ok := manual[t]; ok {
				result[t] = price
				continue
			}
			// 2. cache hit
			if cache != nil {
				if price, ok := cache.GetCurrentPrice(t, CurrentPriceTTL); ok {
					result[t] = price
					continue
				}
			}
			needFetch = append(needFetch, t)
		}

		if len(needFetch) > 0 {
			fetched, err := client.GetPrices(needFetch)
			if err != nil {
				return nil, fmt.Errorf("market data: %w", err)
			}
			for k, v := range fetched {
				result[k] = v
				if cache != nil {
					_ = cache.SetCurrentPrice(k, v)
				}
			}
		}

		return result, nil
	}
}

// NewCachedHistoricalSource returns a function that fetches historical prices
// with cache-aside: hit cache first, fall back to brapi, then populate cache.
// cache may be nil.
func NewCachedHistoricalSource(client *Client, cache PriceCache) func(ticker string) ([]HistoricalPrice, error) {
	return func(ticker string) ([]HistoricalPrice, error) {
		if cache != nil {
			if prices, ok := cache.GetHistoricalPrices(ticker, HistoricalPriceTTL); ok {
				return prices, nil
			}
		}

		prices, err := client.GetHistoricalPrices(ticker)
		if err != nil {
			return nil, err
		}

		if cache != nil {
			_ = cache.SetHistoricalPrices(ticker, prices)
		}

		return prices, nil
	}
}
