package cache

import (
	"time"
)

// GetMacroFloat returns a cached numeric macro value.
// Returns (value, true) if found and within TTL.
func (c *DB) GetMacroFloat(key string, ttl time.Duration) (float64, bool) {
	cutoff := time.Now().Add(-ttl).Format(time.RFC3339)

	var value float64
	err := c.db.QueryRow(`
		SELECT value_num FROM macro_cache
		WHERE key = ? AND fetched_at > ? AND value_num IS NOT NULL`,
		key, cutoff,
	).Scan(&value)

	if err != nil {
		return 0, false
	}
	return value, true
}

// SetMacroFloat stores a numeric macro value.
func (c *DB) SetMacroFloat(key string, value float64) error {
	now := time.Now().Format(time.RFC3339)
	_, err := c.db.Exec(`
		INSERT INTO macro_cache (key, value_num, fetched_at)
		VALUES (?, ?, ?)
		ON CONFLICT(key) DO UPDATE SET value_num=excluded.value_num, fetched_at=excluded.fetched_at`,
		key, value, now,
	)
	return err
}

// GetMacroText returns a cached text macro value.
func (c *DB) GetMacroText(key string, ttl time.Duration) (string, bool) {
	cutoff := time.Now().Add(-ttl).Format(time.RFC3339)

	var value string
	err := c.db.QueryRow(`
		SELECT value_txt FROM macro_cache
		WHERE key = ? AND fetched_at > ? AND value_txt IS NOT NULL`,
		key, cutoff,
	).Scan(&value)

	if err != nil {
		return "", false
	}
	return value, true
}

// SetMacroText stores a text macro value.
func (c *DB) SetMacroText(key string, value string) error {
	now := time.Now().Format(time.RFC3339)
	_, err := c.db.Exec(`
		INSERT INTO macro_cache (key, value_txt, fetched_at)
		VALUES (?, ?, ?)
		ON CONFLICT(key) DO UPDATE SET value_txt=excluded.value_txt, fetched_at=excluded.fetched_at`,
		key, value, now,
	)
	return err
}
