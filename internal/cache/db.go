package cache

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	_ "modernc.org/sqlite"
)

const schema = `
CREATE TABLE IF NOT EXISTS price_cache (
    ticker     TEXT NOT NULL,
    date       TEXT NOT NULL,  -- YYYY-MM-DD
    price      REAL NOT NULL,
    fetched_at TEXT NOT NULL,  -- RFC3339
    PRIMARY KEY (ticker, date)
);

CREATE TABLE IF NOT EXISTS macro_cache (
    key        TEXT PRIMARY KEY,
    value_num  REAL,
    value_txt  TEXT,
    fetched_at TEXT NOT NULL   -- RFC3339
);

CREATE INDEX IF NOT EXISTS idx_price_ticker ON price_cache(ticker);
`

// DB is a persistent cache backed by SQLite.
type DB struct {
	db *sql.DB
}

// DefaultPath returns ~/.investcopilot/cache.db
func DefaultPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("resolving home dir: %w", err)
	}
	dir := filepath.Join(home, ".investcopilot")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return "", fmt.Errorf("creating cache dir: %w", err)
	}
	return filepath.Join(dir, "cache.db"), nil
}

// Open opens (or creates) the cache database at path.
func Open(path string) (*DB, error) {
	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, fmt.Errorf("opening cache db: %w", err)
	}

	db.SetMaxOpenConns(1) // SQLite is single-writer

	if _, err := db.Exec(schema); err != nil {
		db.Close()
		return nil, fmt.Errorf("initializing schema: %w", err)
	}

	return &DB{db: db}, nil
}

// Close closes the database.
func (c *DB) Close() error {
	return c.db.Close()
}
