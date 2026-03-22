# Investment Copilot

Uma ferramenta de linha de comando para apoio à decisão de investimentos, estruturada e orientada por dados.

O Investment Copilot não tenta prever o mercado. Ele mede risco, identifica concentração, interpreta o contexto macro, aplica regras explícitas e sugere ações.

## Visão Geral

```
$ investcopilot analyze -t transactions.csv -c config.yaml -a 2000

  Investment Copilot
  2026-03-22

────────────────────────────────────────────────────────────
  Patrimônio total:  R$ 59.596,00
────────────────────────────────────────────────────────────

  Posições

  TICKER          CLASSE             QTD     PREÇO MÉDIO   PREÇO ATUAL   VALOR          ALOCAÇÃO
  IVVB11          equities_global    50.00   R$ 301.20     R$ 325.00     R$ 16.250,00   27.3%
  PETR4           equities_br       300.00   R$ 37.30      R$ 41.20      R$ 12.360,00   20.7%
  TESOURO_SELIC   renda_fixa          2.00   R$ 5100.00    R$ 5350.00    R$ 10.700,00   18.0%
  ...

  Alocação vs Target

  equities_br           38.4% /  25.0%  ▲  +13.4pp [WARN]
  equities_global       27.3% /  20.0%  ▲  +7.3pp  [WARN]
  renda_fixa            18.0% /  30.0%  ▼  -12.0pp
  real_estate            8.3% /  15.0%  ▼  -6.7pp
  commodities            8.1% /  10.0%  ▼  ok

  Alertas de Risco

  [WARN] PETR4 representa 20.7% da carteira (limite: 20%)
  [WARN] Top 3 ativos representam 66.0% da carteira (limite: 50%)

  Performance

  Retorno total:   +R$ 3.331,00  (+5.92%)
  Total investido: R$ 56.265,00

  BENCHMARK   1 MÊS    6 MESES   12 MESES
  Carteira    —        —         —
  CDI         +1.05%   +7.07%    +14.65%
  IBOV        —        —         —

  Contexto Macro

  Selic meta:    14.75% a.a.  (em queda)
  IPCA 12m:       3.81%
  USD/BRL:        5.28

  Expectativas Focus (mediana, fim do ano):
  IPCA projetado:  4.10%
  Selic projetada: 12.25%

  Sinais do Playbook

  [OPPORTUNITY] Ciclo de queda da Selic (14.75% a.a.)
                Queda de juros historicamente beneficia bolsa e FIIs

  [CAUTION]     equities_br: 13.4pp acima do target
                Rebalancear progressivamente

────────────────────────────────────────────────────────────
  Próximo aporte: R$ 2.000,00
────────────────────────────────────────────────────────────

  1. Aportar R$ 1.165,12   em renda_fixa
  2. Aportar R$   651,54   em real_estate
  3. Aportar R$   183,33   em commodities

  Evitar: equities_br — acima do target
```

## Módulos

| Módulo | O que faz |
|--------|-----------|
| **Portfolio** | Consolida transações, calcula posição (preço médio, valor atual, alocação %) |
| **Risk** | Mede concentração por ativo e classe, detecta desvios do target |
| **Analytics** | Retorno total desde o início e por período (1m/6m/12m) |
| **Benchmark** | Compara performance com CDI (BCB, gratuito) e IBOV (brapi, com token) |
| **Macro** | Busca Selic, IPCA, USD/BRL e expectativas Focus diretamente do BCB |
| **Playbook** | Aplica regras explícitas cruzando macro + carteira, gera sinais OPPORTUNITY/CAUTION/INFO |
| **Recommendation** | Distribui o valor do aporte proporcionalmente ao déficit de cada classe |
| **Report** | Exporta relatório HTML auto-contido (sem dependências externas) |

## Instalação

```bash
go install github.com/phbpx/investcopilot/cmd/investcopilot@latest
```

Ou compilar do fonte:

```bash
git clone https://github.com/phbpx/investcopilot
cd investcopilot
go build -o investcopilot ./cmd/investcopilot
```

Requer Go 1.21+. Sem CGO — binário único, sem dependências de runtime.

## Uso

### Analisar carteira

```bash
investcopilot analyze -t transactions.csv -c config.yaml
```

### Com valor de aporte

```bash
investcopilot analyze -t transactions.csv -c config.yaml -a 2000
```

### Exportar relatório HTML

```bash
investcopilot analyze -t transactions.csv -c config.yaml -a 2000 -o report.html
```

### Flags

| Flag | Atalho | Padrão | Descrição |
|------|--------|--------|-----------|
| `--transactions` | `-t` | — | Caminho para o CSV de transações (obrigatório) |
| `--config` | `-c` | `config.yaml` | Caminho para o arquivo de configuração |
| `--contribution` | `-a` | `0` | Valor disponível para o próximo aporte (BRL) |
| `--output` | `-o` | — | Exportar relatório HTML para este arquivo |

## Configuração

Crie um `config.yaml` com sua estratégia de alocação:

```yaml
target_allocation:
  renda_fixa: 30
  equities_br: 25
  equities_global: 20
  real_estate: 15
  commodities: 10

risk_rules:
  max_single_asset: 20   # % máximo por ativo individual
  max_top3: 50           # % máximo dos 3 maiores ativos
  max_sector: 35         # % máximo por setor

assets:
  PETR4:
    class: equities_br
    sector: energia
  IVVB11:
    class: equities_global
    sector: diversificado
  MXRF11:
    class: real_estate
    sector: papeis
  TESOURO_SELIC:
    class: renda_fixa
    sector: soberano

market:
  # Token gratuito para preços ao vivo e histórico via brapi.dev
  # Sem token: use manual_prices como fallback
  # brapi_token: seu-token-aqui

  manual_prices:
    TESOURO_SELIC: 5350.00
```

### Classes de ativo suportadas

| Classe | Chave no config |
|--------|----------------|
| Renda Fixa | `renda_fixa` |
| Ações BR | `equities_br` |
| Ações Global | `equities_global` |
| Fundos Imobiliários | `real_estate` |
| Commodities | `commodities` |

## Formato do CSV de transações

```csv
date,ticker,type,quantity,price,fees
2024-01-15,PETR4,BUY,100,38.50,5.00
2024-03-01,IVVB11,BUY,50,285.00,5.00
2024-06-10,PETR4,SELL,30,42.00,5.00
```

| Campo | Formato | Descrição |
|-------|---------|-----------|
| `date` | `YYYY-MM-DD` | Data da transação |
| `ticker` | string | Código do ativo |
| `type` | `BUY`, `SELL` ou `INCOME` | Tipo da operação |
| `quantity` | decimal | Quantidade |
| `price` | decimal | Preço unitário (BRL) |
| `fees` | decimal | Corretagem e taxas (BRL) |

## Fontes de dados

| Dado | Fonte | Auth | TTL do cache |
|------|-------|------|-------------|
| Preços atuais (BR) | [brapi.dev](https://brapi.dev) | Token gratuito | 1h |
| Preços históricos | [brapi.dev](https://brapi.dev) | Token gratuito | 24h |
| Selic, IPCA, USD/BRL | [BCB SGS](https://www.bcb.gov.br/estatisticas/tabelaespecial) | Nenhuma | 6h |
| Expectativas Focus | [BCB Olinda](https://olinda.bcb.gov.br) | Nenhuma | 6h |
| CDI (benchmark) | [BCB SGS](https://www.bcb.gov.br/estatisticas/tabelaespecial) | Nenhuma | 24h |

Os dados do BCB (Selic, IPCA, Focus, CDI) são sempre buscados sem autenticação. O token brapi é necessário apenas para preços ao vivo e retornos históricos por período.

## Cache

Os dados de mercado e macro são persistidos em `~/.investcopilot/cache.db` (SQLite). Isso evita requisições repetidas nas execuções seguintes enquanto os dados estão dentro do TTL.

Para inspecionar o cache:

```bash
sqlite3 ~/.investcopilot/cache.db "SELECT key, ROUND(value_num,4), fetched_at FROM macro_cache;"
sqlite3 ~/.investcopilot/cache.db "SELECT ticker, date, price FROM price_cache LIMIT 20;"
```

Para limpar o cache:

```bash
rm ~/.investcopilot/cache.db
```

## Roadmap

- [x] Portfolio Engine + CLI base
- [x] Risk Engine
- [x] Recommendation Engine
- [x] Analytics Engine (retorno total e por período)
- [x] Benchmark Engine (CDI, IBOV)
- [x] Macro Engine (Selic, IPCA, USD/BRL, Focus)
- [x] Playbook Engine (regras macro + carteira)
- [x] Report Engine (HTML export)
- [x] Cache persistido (SQLite)

## Licença

MIT
