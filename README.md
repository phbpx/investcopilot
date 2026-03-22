# Investment Copilot

A CLI tool for structured, data-driven investment decision-making.

Investment Copilot does not try to predict the market. It measures risk, identifies concentration, interprets macro context, applies explicit rules, and suggests actions.

## Overview

```
$ investcopilot analyze --transactions transactions.csv --config config.yaml

Portfolio Summary (2026-03-22)
──────────────────────────────────────────────────────
Total:  R$ 87.432,00

Allocation vs Target:
  Renda Fixa        28% / 30%   -2pp
  Equities BR       32% / 25%   +7pp  [WARN] above target
  Equities Global   19% / 20%   -1pp
  Real Estate       15% / 15%    ok
  Commodities        6% / 10%   -4pp

Risk Alerts:
  [WARN] PETR4: 22% of portfolio (limit: 20%)

Next Contribution: R$ 2.000
  → R$ 1.200 to Renda Fixa
  → R$   800 to Commodities
  → Avoid Equities BR until rebalanced
```

## Features

- **Portfolio Engine** — consolidates transactions, calculates current positions and allocations
- **Risk Engine** — measures concentration, identifies hidden risk
- **Analytics Engine** — returns (1m/6m/12m), volatility, max drawdown
- **Benchmark Engine** — compares performance against CDI, IBOV, IPCA+5%, S&P 500
- **Playbook Engine** — applies explicit decision rules based on portfolio state and macro context
- **Recommendation Engine** — generates concrete next-action suggestions
- **Report Engine** — exports a static HTML report for browser viewing

## Installation

```bash
go install github.com/phbpx/investcopilot/cmd/investcopilot@latest
```

Or build from source:

```bash
git clone https://github.com/phbpx/investcopilot
cd investcopilot
go build -o investcopilot ./cmd/investcopilot
```

## Usage

### Analyze portfolio

```bash
investcopilot analyze --transactions transactions.csv --config config.yaml
```

### Generate HTML report

```bash
investcopilot report --transactions transactions.csv --config config.yaml --output report.html
```

## Configuration

Create a `config.yaml` file with your target allocation and risk rules:

```yaml
target_allocation:
  renda_fixa: 30
  equities_br: 25
  equities_global: 20
  real_estate: 15
  commodities: 10

risk_rules:
  max_single_asset: 20
  max_top3: 50
  max_sector: 35
```

## Transactions CSV format

```csv
date,ticker,type,quantity,price,fees
2024-01-15,PETR4,BUY,100,38.50,5.00
2024-03-01,IVVB11,BUY,50,285.00,5.00
2024-06-10,PETR4,SELL,30,42.00,5.00
```

## Roadmap

- [x] Project definition (PRD)
- [ ] Portfolio Engine + CLI base
- [ ] Risk Engine
- [ ] Recommendation Engine
- [ ] Analytics Engine + Benchmark Engine
- [ ] Playbook Engine (macro rules)
- [ ] Report Engine (HTML export)

## License

MIT
