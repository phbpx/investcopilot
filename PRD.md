# PRD — Investment Copilot

---

# 1. Visão do Produto

**Investment Copilot** é um sistema de apoio à decisão que transforma uma carteira de investimentos em um processo estruturado, disciplinado e orientado por dados.

O sistema não tenta prever o mercado. Ele:

- mede risco
- identifica concentração
- interpreta contexto macroeconômico
- aplica regras explícitas e configuráveis
- sugere ações concretas

---

# 2. Objetivo

Melhorar:

- qualidade das decisões de aporte e rebalanceamento
- consistência de alocação em relação a uma estratégia definida
- controle de risco e concentração
- disciplina de investimento ao longo do tempo

---

# 3. Usuário

## Perfil atual (fase de aprendizado)

- investidor individual
- perfil agressivo
- técnico (engenheiro)
- toma decisões próprias
- usa o sistema como ferramenta pessoal de análise

## Perfil futuro (após validação do produto)

- investidores não-técnicos que queiram um copiloto de carteira
- pessoas que precisam de um processo mas não têm tempo para montar um

> Decisão: começar como CLI pessoal, aprender sobre o produto, depois avaliar expansão.

---

# 4. Problemas

| Problema | Descrição |
|----------|-----------|
| Decisão inconsistente | aportes sem critério claro, baseados em feeling |
| Falta de visão consolidada | não sabe exposição real por ativo e classe |
| Concentração invisível | risco oculto que só aparece quando consolida tudo |
| Falta de disciplina | decisões emocionais em momentos de volatilidade |
| Ignorar contexto macro | timing ruim por desconhecer o ambiente |
| Sem benchmark | não sabe se está performando bem ou mal |

---

# 5. Proposta de Valor

> Transformar investimento em um sistema de decisão repetível, orientado por dados e por uma estratégia explícita.

---

# 6. Produto

## Interface

**CLI-first.** O sistema é uma ferramenta de linha de comando.

- sem servidor, sem infraestrutura, sem dependências de runtime
- output rico no terminal (tabelas, alertas, cores)
- exportação de relatório HTML estático para visualização no browser

Não há interface gráfica no roadmap atual. Se o produto evoluir para outros usuários, avaliar na época.

## Stack

- **Linguagem:** Go
- **CLI:** cobra
- **Dados de mercado (BR):** brapi.dev (gratuito, sem autenticação)
- **Dados de mercado (global):** Yahoo Finance
- **Dados macro:** API Focus (Banco Central do Brasil)
- **Configuração:** YAML

---

# 7. Escopo

## MVP (Fase 1) — CLI funcional

- importação de transações via CSV
- configuração de alocação-alvo via YAML
- posição consolidada com preço atual (via brapi.dev)
- análise de risco e concentração
- desvio vs alocação-alvo
- recomendação de próximo aporte
- output no terminal

## Fase 2 — Analytics e Benchmark

- retorno histórico (1m / 6m / 12m)
- volatilidade e drawdown máximo
- comparação com benchmarks (CDI, IBOV, IPCA+5%, S&P 500 em BRL)
- exportação de relatório HTML

## Fase 3 — Contexto Macro e Playbook

- integração com API Focus (BCB)
- regime de mercado (Selic, IPCA, câmbio)
- regras do playbook baseadas em contexto macro
- scoring de oportunidade por classe de ativo

## Fora de escopo (por enquanto)

- execução automática de trades
- machine learning / modelos preditivos
- integração direta com corretoras
- interface gráfica / aplicação web

---

# 8. Funcionalidades

---

## 8.0 Configuração de Alocação-Alvo

O usuário define sua estratégia em `config.yaml`. Essa é a base de todas as decisões do sistema.

```yaml
target_allocation:
  renda_fixa: 30
  equities_br: 25
  equities_global: 20
  real_estate: 15
  commodities: 10

risk_rules:
  max_single_asset: 20    # % máximo por ativo individual
  max_top3: 50            # % máximo dos 3 maiores ativos
  max_sector: 35          # % máximo por setor
```

---

## 8.1 Portfolio Engine

Consolida as transações e calcula a posição atual.

### Input

Arquivo CSV de transações:

```csv
date,ticker,type,quantity,price,fees
2024-01-15,PETR4,BUY,100,38.50,5.00
2024-03-01,IVVB11,BUY,50,285.00,5.00
2024-06-10,PETR4,SELL,30,42.00,5.00
```

### Output

- posição consolidada por ativo (quantidade, preço médio, valor atual)
- alocação percentual por ativo
- alocação percentual por classe
- desvio vs alocação-alvo

---

## 8.2 Risk Engine

Avalia risco e concentração da carteira.

### Métricas

- concentração por ativo (%)
- concentração por classe (%)
- peso dos top 3 ativos
- exposição por setor

### Regras (configuráveis)

| Condição | Alerta |
|----------|--------|
| Ativo > 20% da carteira | WARN: concentração excessiva |
| Top 3 ativos > 50% | WARN: carteira concentrada |
| Setor > 35% | WARN: exposição setorial alta |
| Desvio de classe > 10pp do target | INFO: rebalanceamento sugerido |

---

## 8.3 Analytics Engine

Mede performance da carteira ao longo do tempo.

### Métricas

- retorno em 1m / 6m / 12m / desde o início
- volatilidade anualizada
- drawdown máximo
- evolução do patrimônio

> Requer histórico de preços. Disponível a partir da Fase 2.

---

## 8.4 Benchmark Engine

Compara a performance da carteira com referências relevantes.

### Benchmarks

| Benchmark | Referência |
|-----------|------------|
| CDI | Renda fixa / custo de oportunidade |
| IBOVESPA | Equities brasileiros |
| IPCA+5% | Retorno real mínimo aceitável |
| S&P 500 (BRL) | Equities globais |

### Output

- retorno da carteira vs cada benchmark no período
- alpha gerado (ou destruído)

---

## 8.5 Playbook Engine

Aplica regras explícitas para gerar sinais de decisão.

### Regras de alocação

| Condição | Sinal |
|----------|-------|
| Classe X% acima do target | REDUZIR: direcionar próximo aporte para outro ativo |
| Classe X% abaixo do target | AUMENTAR: priorizar essa classe no próximo aporte |
| Drawdown > 15% em ativo | REVISAR: checar tese de investimento |
| Volatilidade da carteira > limite | ALERTA: revisar exposição a risco |

### Regras macro (Fase 3)

| Condição | Sinal |
|----------|-------|
| Selic > 13% e subindo | AUMENTAR renda fixa |
| IPCA > 6% | CHECAR proteção à inflação (IPCA+, FIIs) |
| Dólar > limiar configurado | REVISAR exposição internacional |
| Expectativa Focus deteriorando | CAUTELA: reduzir risco |

---

## 8.6 Recommendation Engine

Consolida todos os sinais e gera a recomendação de ação.

### Output esperado (terminal)

```
Próximo aporte: R$ 2.000

  1. Aportar R$ 1.200 em Renda Fixa          (+6pp abaixo do target)
  2. Aportar R$   800 em Equities Global     (+1pp abaixo do target)

  Evitar: Equities BR — já 7pp acima do target
  Monitorar: PETR4 — 22% da carteira (limite: 20%)
```

---

## 8.7 Report Engine

Gera um relatório HTML estático para visualização no browser.

```
$ investcopilot report --output report.html
Relatório gerado: report.html
```

### Conteúdo do relatório

- posição atual vs target (gráfico de barras)
- performance vs benchmark (linha do tempo)
- alertas ativos
- recomendação de aporte
- evolução histórica do patrimônio

---

# 9. Fluxo de uso

```
$ investcopilot analyze --transactions transactions.csv --config config.yaml

Portfolio Summary (2026-03-22)
──────────────────────────────────────────────────────
Patrimônio total:  R$ 87.432,00

Alocação atual vs target:
  Renda Fixa        28% / 30%   -2pp
  Equities BR       32% / 25%   +7pp  [WARN] acima do target
  Equities Global   19% / 20%   -1pp
  Real Estate       15% / 15%    ok
  Commodities        6% / 10%   -4pp

Alertas de risco:
  [WARN] PETR4: 22% da carteira (limite: 20%)

Próximo aporte: R$ 2.000
  → R$ 1.200 em Renda Fixa
  → R$   800 em Commodities
  → Evitar Equities BR até rebalancear
```

---

# 10. Decisões de produto registradas

| Decisão | Justificativa |
|---------|---------------|
| CLI em vez de web app | Uso pessoal, perfil técnico, zero infraestrutura |
| Sem interface gráfica no MVP | Complexidade não justificada para uso individual |
| Relatório HTML estático | Visualização sem servidor ou deploy |
| Go + cobra | Performance, binário único, sem dependências de runtime |
| brapi.dev para preços BR | Gratuito, sem necessidade de autenticação |
| Regras explícitas no playbook | Transparência e controle total sobre as decisões |
| Começar pessoal, escalar depois | Aprender sobre o produto antes de construir para outros |
