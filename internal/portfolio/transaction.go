package portfolio

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"time"
)

// TransactionType represents a buy or sell operation.
type TransactionType string

const (
	Buy  TransactionType = "BUY"
	Sell TransactionType = "SELL"
)

// Transaction represents a single buy or sell event.
type Transaction struct {
	Date     time.Time
	Ticker   string
	Type     TransactionType
	Quantity float64
	Price    float64
	Fees     float64
}

// LoadTransactions reads and parses a CSV file of transactions.
// Expected columns: date,ticker,type,quantity,price,fees
func LoadTransactions(path string) ([]Transaction, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("opening transactions file: %w", err)
	}
	defer f.Close()

	r := csv.NewReader(f)
	r.TrimLeadingSpace = true

	records, err := r.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("reading csv: %w", err)
	}

	if len(records) < 2 {
		return nil, fmt.Errorf("transactions file is empty")
	}

	txs := make([]Transaction, 0, len(records)-1)
	for i, row := range records[1:] {
		if len(row) < 6 {
			return nil, fmt.Errorf("row %d: expected 6 columns, got %d", i+2, len(row))
		}

		date, err := time.Parse("2006-01-02", row[0])
		if err != nil {
			return nil, fmt.Errorf("row %d: invalid date %q: %w", i+2, row[0], err)
		}

		qty, err := strconv.ParseFloat(row[3], 64)
		if err != nil {
			return nil, fmt.Errorf("row %d: invalid quantity %q: %w", i+2, row[3], err)
		}

		price, err := strconv.ParseFloat(row[4], 64)
		if err != nil {
			return nil, fmt.Errorf("row %d: invalid price %q: %w", i+2, row[4], err)
		}

		fees, err := strconv.ParseFloat(row[5], 64)
		if err != nil {
			return nil, fmt.Errorf("row %d: invalid fees %q: %w", i+2, row[5], err)
		}

		txType := TransactionType(row[2])
		if txType != Buy && txType != Sell {
			return nil, fmt.Errorf("row %d: invalid type %q (expected BUY or SELL)", i+2, row[2])
		}

		txs = append(txs, Transaction{
			Date:     date,
			Ticker:   row[1],
			Type:     txType,
			Quantity: qty,
			Price:    price,
			Fees:     fees,
		})
	}

	return txs, nil
}
