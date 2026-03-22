package portfolio

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// AssetClass represents a broad asset class category.
type AssetClass string

const (
	ClassRendaFixa       AssetClass = "renda_fixa"
	ClassEquitiesBR      AssetClass = "equities_br"
	ClassEquitiesGlobal  AssetClass = "equities_global"
	ClassRealEstate      AssetClass = "real_estate"
	ClassCommodities     AssetClass = "commodities"
)

// Config holds the user's investment strategy configuration.
type Config struct {
	TargetAllocation map[AssetClass]float64 `yaml:"target_allocation"`
	RiskRules        RiskRules              `yaml:"risk_rules"`
	Assets           map[string]AssetConfig `yaml:"assets"`
	Market           MarketConfig           `yaml:"market"`
}

// MarketConfig holds market data provider settings.
type MarketConfig struct {
	BrapiToken    string             `yaml:"brapi_token"`
	ManualPrices  map[string]float64 `yaml:"manual_prices"` // fallback: ticker → price
}

// RiskRules defines the risk thresholds.
type RiskRules struct {
	MaxSingleAsset float64 `yaml:"max_single_asset"`
	MaxTop3        float64 `yaml:"max_top3"`
	MaxSector      float64 `yaml:"max_sector"`
}

// AssetConfig maps a ticker to its asset class and sector.
type AssetConfig struct {
	Class  AssetClass `yaml:"class"`
	Sector string     `yaml:"sector"`
}

// LoadConfig reads and parses the config YAML file.
func LoadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading config: %w", err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("parsing config: %w", err)
	}

	if err := cfg.validate(); err != nil {
		return nil, fmt.Errorf("invalid config: %w", err)
	}

	return &cfg, nil
}

func (c *Config) validate() error {
	total := 0.0
	for _, pct := range c.TargetAllocation {
		total += pct
	}
	if total < 99 || total > 101 {
		return fmt.Errorf("target_allocation must sum to 100 (got %.1f)", total)
	}
	return nil
}

// ClassOf returns the asset class for a given ticker.
func (c *Config) ClassOf(ticker string) AssetClass {
	if a, ok := c.Assets[ticker]; ok {
		return a.Class
	}
	return ""
}
