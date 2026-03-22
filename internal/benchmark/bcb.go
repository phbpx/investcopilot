package benchmark

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

const bcbBaseURL = "https://api.bcb.gov.br/dados/serie/bcdata.sgs"

// series codes from BCB SGS
const (
	seriesCDI  = 12  // CDI daily rate
	seriesIPCA = 433 // IPCA monthly rate
)

type bcbRecord struct {
	Date  string `json:"data"`  // "DD/MM/YYYY"
	Value string `json:"valor"` // decimal as string
}

// fetchBCB fetches records from BCB SGS for a given series and date range.
func fetchBCB(series int, from, to time.Time) ([]bcbRecord, error) {
	url := fmt.Sprintf(
		"%s.%d/dados?dataInicial=%s&dataFinal=%s&formato=json",
		bcbBaseURL, series,
		from.Format("02/01/2006"),
		to.Format("02/01/2006"),
	)

	resp, err := http.Get(url) //nolint:gosec
	if err != nil {
		return nil, fmt.Errorf("fetching BCB series %d: %w", series, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("BCB API returned status %d for series %d", resp.StatusCode, series)
	}

	var records []bcbRecord
	if err := json.NewDecoder(resp.Body).Decode(&records); err != nil {
		return nil, fmt.Errorf("decoding BCB response: %w", err)
	}

	return records, nil
}

// cdiAccumulated returns the accumulated CDI return (%) for the period [from, to].
func cdiAccumulated(from, to time.Time) (float64, error) {
	records, err := fetchBCB(seriesCDI, from, to)
	if err != nil {
		return 0, err
	}
	if len(records) == 0 {
		return 0, fmt.Errorf("no CDI data for period")
	}

	accumulated := 1.0
	for _, r := range records {
		v, err := strconv.ParseFloat(r.Value, 64)
		if err != nil {
			continue
		}
		accumulated *= 1 + v/100
	}

	return (accumulated - 1) * 100, nil
}
