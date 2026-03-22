package market

import "fmt"

// NewPriceSource returns a PriceSource that uses manual prices when available,
// falling back to brapi.dev for any tickers not covered manually.
func NewPriceSource(manual map[string]float64, token string) func(tickers []string) (map[string]float64, error) {
	client := NewClient(token)

	return func(tickers []string) (map[string]float64, error) {
		result := make(map[string]float64, len(tickers))

		var needFetch []string
		for _, t := range tickers {
			if price, ok := manual[t]; ok {
				result[t] = price
			} else {
				needFetch = append(needFetch, t)
			}
		}

		if len(needFetch) > 0 {
			fetched, err := client.GetPrices(needFetch)
			if err != nil {
				return nil, fmt.Errorf("market data: %w", err)
			}
			for k, v := range fetched {
				result[k] = v
			}
		}

		return result, nil
	}
}
