package cache

import (
	"fmt"
	"time"

	"github.com/phbpx/investcopilot/internal/market"
)

// GetCurrentPrice returns the cached price for ticker on today's date.
// Returns (price, true) if found and within TTL, (0, false) otherwise.
func (c *DB) GetCurrentPrice(ticker string, ttl time.Duration) (float64, bool) {
	today := time.Now().Format("2006-01-02")
	cutoff := time.Now().Add(-ttl).Format(time.RFC3339)

	var price float64
	err := c.db.QueryRow(`
		SELECT price FROM price_cache
		WHERE ticker = ? AND date = ? AND fetched_at > ?`,
		ticker, today, cutoff,
	).Scan(&price)

	if err != nil {
		return 0, false
	}
	return price, true
}

// SetCurrentPrice stores the current price for ticker.
func (c *DB) SetCurrentPrice(ticker string, price float64) error {
	today := time.Now().Format("2006-01-02")
	now := time.Now().Format(time.RFC3339)

	_, err := c.db.Exec(`
		INSERT INTO price_cache (ticker, date, price, fetched_at)
		VALUES (?, ?, ?, ?)
		ON CONFLICT(ticker, date) DO UPDATE SET price=excluded.price, fetched_at=excluded.fetched_at`,
		ticker, today, price, now,
	)
	return err
}

// GetHistoricalPrices returns all cached historical prices for ticker if fresh.
// "Fresh" means at least one record was fetched within TTL (implies a full sync happened).
func (c *DB) GetHistoricalPrices(ticker string, ttl time.Duration) ([]market.HistoricalPrice, bool) {
	cutoff := time.Now().Add(-ttl).Format(time.RFC3339)

	// check if the cache has a recent fetch for this ticker
	var count int
	err := c.db.QueryRow(`
		SELECT COUNT(*) FROM price_cache
		WHERE ticker = ? AND fetched_at > ?`, ticker, cutoff,
	).Scan(&count)
	if err != nil || count == 0 {
		return nil, false
	}

	rows, err := c.db.Query(`
		SELECT date, price FROM price_cache
		WHERE ticker = ?
		ORDER BY date ASC`, ticker,
	)
	if err != nil {
		return nil, false
	}
	defer rows.Close()

	var prices []market.HistoricalPrice
	for rows.Next() {
		var dateStr string
		var price float64
		if err := rows.Scan(&dateStr, &price); err != nil {
			continue
		}
		date, err := time.Parse("2006-01-02", dateStr)
		if err != nil {
			continue
		}
		prices = append(prices, market.HistoricalPrice{Date: date, Close: price})
	}

	if len(prices) == 0 {
		return nil, false
	}
	return prices, true
}

// SetHistoricalPrices stores a batch of historical prices for ticker.
func (c *DB) SetHistoricalPrices(ticker string, prices []market.HistoricalPrice) error {
	now := time.Now().Format(time.RFC3339)

	tx, err := c.db.Begin()
	if err != nil {
		return fmt.Errorf("beginning transaction: %w", err)
	}
	defer tx.Rollback() //nolint:errcheck

	stmt, err := tx.Prepare(`
		INSERT INTO price_cache (ticker, date, price, fetched_at)
		VALUES (?, ?, ?, ?)
		ON CONFLICT(ticker, date) DO UPDATE SET price=excluded.price, fetched_at=excluded.fetched_at`)
	if err != nil {
		return fmt.Errorf("preparing statement: %w", err)
	}
	defer stmt.Close()

	for _, p := range prices {
		if _, err := stmt.Exec(ticker, p.Date.Format("2006-01-02"), p.Close, now); err != nil {
			return fmt.Errorf("inserting price for %s on %s: %w", ticker, p.Date.Format("2006-01-02"), err)
		}
	}

	return tx.Commit()
}
