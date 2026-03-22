package macro

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"strconv"
	"time"
)

const (
	bcbSGSBase  = "https://api.bcb.gov.br/dados/serie/bcdata.sgs"
	focusBase   = "https://olinda.bcb.gov.br/olinda/servico/Expectativas/versao/v1/odata"
	seriesSelic = 432  // Selic meta diária
	seriesIPCA  = 433  // IPCA mensal
	seriesUSD   = 1    // USD/BRL
)

var httpClient = &http.Client{Timeout: 10 * time.Second}

type sgsRecord struct {
	Date  string `json:"data"`
	Value string `json:"valor"`
}

func fetchSGS(series int, from, to time.Time) ([]sgsRecord, error) {
	u := fmt.Sprintf("%s.%d/dados?dataInicial=%s&dataFinal=%s&formato=json",
		bcbSGSBase, series,
		from.Format("02/01/2006"),
		to.Format("02/01/2006"),
	)
	resp, err := httpClient.Get(u) //nolint:gosec
	if err != nil {
		return nil, fmt.Errorf("BCB SGS series %d: %w", series, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("BCB SGS series %d: status %d", series, resp.StatusCode)
	}

	var records []sgsRecord
	if err := json.NewDecoder(resp.Body).Decode(&records); err != nil {
		return nil, fmt.Errorf("BCB SGS series %d decode: %w", series, err)
	}
	return records, nil
}

func latestFloat(series int) (float64, error) {
	now := time.Now()
	records, err := fetchSGS(series, now.AddDate(0, -1, 0), now)
	if err != nil {
		return 0, err
	}
	if len(records) == 0 {
		return 0, fmt.Errorf("no data for series %d", series)
	}
	last := records[len(records)-1]
	v, err := strconv.ParseFloat(last.Value, 64)
	if err != nil {
		return 0, fmt.Errorf("parsing value %q: %w", last.Value, err)
	}
	return v, nil
}

// selicTrend fetches Selic over the last 6 months and detects direction.
func selicTrend() (current float64, trend Trend, err error) {
	now := time.Now()
	records, err := fetchSGS(seriesSelic, now.AddDate(0, -6, 0), now)
	if err != nil {
		return 0, TrendStable, err
	}
	if len(records) == 0 {
		return 0, TrendStable, fmt.Errorf("no Selic data")
	}

	first, _ := strconv.ParseFloat(records[0].Value, 64)
	last, _ := strconv.ParseFloat(records[len(records)-1].Value, 64)
	current = last

	switch {
	case last > first:
		trend = TrendRising
	case last < first:
		trend = TrendFalling
	default:
		trend = TrendStable
	}
	return current, trend, nil
}

// ipcaAccumulated12m returns IPCA accumulated over the last 12 months.
func ipcaAccumulated12m() (float64, error) {
	now := time.Now()
	records, err := fetchSGS(seriesIPCA, now.AddDate(-1, 0, 0), now)
	if err != nil {
		return 0, err
	}
	if len(records) == 0 {
		return 0, fmt.Errorf("no IPCA data")
	}

	acc := 1.0
	for _, r := range records {
		v, _ := strconv.ParseFloat(r.Value, 64)
		acc *= 1 + v/100
	}
	return (acc - 1) * 100, nil
}

// focusRecord is the OData response envelope.
type focusEnvelope struct {
	Value []struct {
		Indicador      string  `json:"Indicador"`
		Data           string  `json:"Data"`
		DataReferencia string  `json:"DataReferencia"`
		Mediana        float64 `json:"Mediana"`
	} `json:"value"`
}

// fetchFocusExpectation returns the latest median expectation for an indicator/year.
func fetchFocusExpectation(indicator string, year int) (float64, error) {
	filter := fmt.Sprintf("Indicador eq '%s' and DataReferencia eq '%d'", indicator, year)
	// OData requires %20 for spaces; url.QueryEscape uses + which the BCB API rejects.
	encoded := strings.NewReplacer("+", "%20").Replace(url.QueryEscape(filter))
	u := fmt.Sprintf(
		"%s/ExpectativasMercadoAnuais?$filter=%s&$orderby=Data%%20desc&$top=1&$format=json&$select=Indicador,Data,DataReferencia,Mediana",
		focusBase, encoded,
	)

	resp, err := httpClient.Get(u) //nolint:gosec
	if err != nil {
		return 0, fmt.Errorf("focus API %s: %w", indicator, err)
	}
	defer resp.Body.Close()

	var env focusEnvelope
	if err := json.NewDecoder(resp.Body).Decode(&env); err != nil {
		return 0, fmt.Errorf("focus API decode: %w", err)
	}
	if len(env.Value) == 0 {
		return 0, fmt.Errorf("no Focus data for %s %d", indicator, year)
	}
	return env.Value[0].Mediana, nil
}
