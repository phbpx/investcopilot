# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Commands

```bash
make build        # compile binary to ./bin/investcopilot
make lint         # go vet + staticcheck
make test         # go test -race -count=1 ./...
make test/cover   # run tests with coverage HTML report
make tidy         # go mod tidy && go mod verify
make all          # lint + test + build (default)
```

Run a single test:
```bash
go test -run TestFunctionName ./internal/package/
```

`staticcheck` must be installed separately: `go install honnef.co/go/tools/cmd/staticcheck@latest`

Run the tool:
```bash
./bin/investcopilot analyze -t transactions.csv -c config.yaml [-a 2000] [-o report.html]
```

## Architecture

The app is a CLI portfolio analyzer for Brazilian investments. It ingests a transaction CSV and config YAML, fetches market/macro data, and outputs a terminal report or HTML file.

### Package layout

```
cmd/investcopilot/
  main.go              # Cobra root command
  analyze/
    command.go         # Orchestrates the full analysis pipeline
    printer.go         # Colored terminal output (tablewriter)

internal/
  portfolio/           # Transaction CSV parsing → position consolidation
  market/              # brapi.dev price fetching + fallback chain
  cache/               # SQLite persistence (~/.investcopilot/cache.db)
  macro/               # BCB API (Selic, IPCA, USD/BRL, Focus expectations)
  risk/                # Concentration & class-deviation alerts
  analytics/           # Period returns (modified Dietz), total return
  benchmark/           # CDI (BCB) and IBOV comparisons
  recommendation/      # Allocate contributions to under-weight classes
  playbook/            # Decision signals from macro + portfolio state
  report/              # Embedded HTML template rendering

examples/
  config.yaml          # Reference configuration
```

### Data flow

`analyze/command.go` orchestrates a sequential pipeline:

1. **Load config + transactions** — YAML config defines `target_allocation`, `risk_rules`, and per-ticker `assets` mapping (class + sector); CSV rows are BUY / SELL / INCOME operations
2. **Build portfolio** (`portfolio.Build`) — processes transactions into holdings, then fetches current prices via the price source
3. **Price source fallback chain**: manual prices (config) → SQLite cache (1h TTL) → brapi.dev API
4. **Analysis engines** run in sequence: risk → recommendations → analytics → benchmarks → macro → playbook
5. **Output**: terminal printer or HTML report

### Key design points

- **Cache-aside pattern**: all external data (prices, macro) is stored in SQLite; TTLs are 1h (current prices), 24h (historical), 6h (macro)
- **Interfaces for the cache**: `PriceCache` and `MacroCache` interfaces in `market/` and `macro/` let the SQLite DB satisfy both without coupling
- **Transaction state machine** (`portfolio/engine.go`): BUY accumulates cost basis, SELL reduces it proportionally, INCOME tracks dividends/FII rendimentos/JCP separately
- **Analytics require a brapi token**: historical price fetching (needed for 1m/6m/12m returns and IBOV benchmark) requires a free token from brapi.dev; without it, performance sections are skipped
- **Pure Go SQLite**: `modernc.org/sqlite` — no CGO, cross-platform

### External APIs

| Source | Data | Auth |
|--------|------|------|
| brapi.dev | Current + historical BR equity prices | Optional free token |
| BCB SGS | Selic (series 432), IPCA (433), USD/BRL (1) | None |
| BCB Olinda (Focus) | Year-end Selic/IPCA expectations | None |
