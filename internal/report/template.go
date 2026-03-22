package report

const htmlTemplate = `<!DOCTYPE html>
<html lang="pt-BR">
<head>
<meta charset="UTF-8">
<meta name="viewport" content="width=device-width, initial-scale=1.0">
<title>Investment Copilot — {{ date .GeneratedAt }}</title>
<style>
*, *::before, *::after { box-sizing: border-box; margin: 0; padding: 0; }

:root {
  --bg:       #0f172a;
  --surface:  #1e293b;
  --border:   #334155;
  --text:     #f1f5f9;
  --muted:    #94a3b8;
  --green:    #22c55e;
  --yellow:   #f59e0b;
  --blue:     #3b82f6;
  --red:      #ef4444;
  --cyan:     #06b6d4;
}

body {
  font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, monospace;
  background: var(--bg);
  color: var(--text);
  font-size: 14px;
  line-height: 1.6;
  padding: 2rem;
  max-width: 1100px;
  margin: 0 auto;
}

h1 { font-size: 1.5rem; font-weight: 700; }
h2 { font-size: 1rem; font-weight: 600; text-transform: uppercase;
     letter-spacing: .08em; color: var(--muted); margin-bottom: 1rem; }

section { margin-bottom: 2.5rem; }

/* header */
.header { display: flex; justify-content: space-between; align-items: flex-start;
          padding-bottom: 1.5rem; border-bottom: 1px solid var(--border); margin-bottom: 2.5rem; }
.header-left .subtitle { color: var(--muted); font-size: .85rem; margin-top: .25rem; }
.total-value { text-align: right; }
.total-value .label { color: var(--muted); font-size: .8rem; }
.total-value .amount { font-size: 1.75rem; font-weight: 700; color: var(--green); }

/* cards row */
.cards { display: grid; grid-template-columns: repeat(auto-fit, minmax(160px, 1fr)); gap: 1rem; margin-bottom: 2.5rem; }
.card { background: var(--surface); border: 1px solid var(--border); border-radius: .5rem; padding: 1rem; }
.card .card-label { color: var(--muted); font-size: .75rem; text-transform: uppercase; letter-spacing: .06em; }
.card .card-value { font-size: 1.15rem; font-weight: 600; margin-top: .2rem; }
.card.green .card-value { color: var(--green); }
.card.yellow .card-value { color: var(--yellow); }
.card.blue .card-value { color: var(--blue); }
.card.red .card-value { color: var(--red); }

/* table */
table { width: 100%; border-collapse: collapse; }
th { text-align: left; padding: .5rem .75rem; font-size: .75rem; text-transform: uppercase;
     letter-spacing: .06em; color: var(--muted); border-bottom: 1px solid var(--border); }
td { padding: .6rem .75rem; border-bottom: 1px solid var(--border); }
tr:last-child td { border-bottom: none; }
tr:hover td { background: var(--surface); }
.mono { font-family: monospace; }
.right { text-align: right; }
.muted { color: var(--muted); }

/* allocation bars */
.alloc-table td { vertical-align: middle; }
.alloc-class { font-weight: 500; min-width: 150px; }
.alloc-pcts { color: var(--muted); font-size: .85rem; white-space: nowrap; min-width: 120px; }
.alloc-dev { font-size: .85rem; font-weight: 600; white-space: nowrap; min-width: 80px; }
.alloc-dev.above { color: var(--yellow); }
.alloc-dev.below { color: var(--blue); }
.alloc-dev.ok    { color: var(--green); }

/* alerts */
.alerts { display: flex; flex-direction: column; gap: .5rem; }
.alert { display: flex; gap: .75rem; align-items: flex-start;
         padding: .6rem .85rem; border-radius: .375rem; }
.alert.warn { background: rgba(245,158,11,.1); border: 1px solid rgba(245,158,11,.3); }
.alert.info { background: rgba(6,182,212,.1);  border: 1px solid rgba(6,182,212,.3); }
.alert-badge { font-size: .7rem; font-weight: 700; text-transform: uppercase;
               letter-spacing: .06em; padding: .15rem .4rem; border-radius: .25rem;
               white-space: nowrap; margin-top: .15rem; }
.alert.warn .alert-badge { background: var(--yellow); color: #000; }
.alert.info .alert-badge { background: var(--cyan);   color: #000; }

/* performance */
.perf-summary { display: grid; grid-template-columns: repeat(auto-fit, minmax(180px, 1fr));
                gap: 1rem; margin-bottom: 1.5rem; }
.perf-card { background: var(--surface); border: 1px solid var(--border);
             border-radius: .5rem; padding: 1rem; }
.perf-card .label { color: var(--muted); font-size: .75rem; text-transform: uppercase; letter-spacing: .06em; }
.perf-card .value { font-size: 1.2rem; font-weight: 600; margin-top: .2rem; }
.perf-card .sub { color: var(--muted); font-size: .8rem; margin-top: .15rem; }
.positive { color: var(--green); }
.negative { color: var(--red); }
.neutral  { color: var(--muted); }

/* benchmark table */
.bench-table th, .bench-table td { text-align: center; }
.bench-table th:first-child, .bench-table td:first-child { text-align: left; }
.bench-table tr.portfolio-row td { font-weight: 600; background: rgba(34,197,94,.05); }

/* macro */
.macro-grid { display: grid; grid-template-columns: repeat(auto-fit, minmax(200px, 1fr)); gap: 1rem; }
.macro-item { background: var(--surface); border: 1px solid var(--border);
              border-radius: .5rem; padding: 1rem; }
.macro-item .m-label { color: var(--muted); font-size: .75rem; text-transform: uppercase; letter-spacing: .06em; }
.macro-item .m-value { font-size: 1.3rem; font-weight: 700; margin-top: .25rem; }
.macro-item .m-value.yellow { color: var(--yellow); }
.macro-item .m-sub   { color: var(--muted); font-size: .8rem; margin-top: .1rem; }

/* playbook signals */
.signals { display: flex; flex-direction: column; gap: .75rem; }
.signal { padding: .85rem 1rem; border-radius: .5rem; border-left: 3px solid; }
.signal.opportunity { background: rgba(34,197,94,.08);  border-color: var(--green); }
.signal.caution     { background: rgba(245,158,11,.08); border-color: var(--yellow); }
.signal.info        { background: rgba(6,182,212,.08);  border-color: var(--cyan); }
.signal-header { display: flex; gap: .6rem; align-items: center; margin-bottom: .25rem; }
.signal-badge { font-size: .65rem; font-weight: 700; text-transform: uppercase;
                letter-spacing: .08em; padding: .15rem .45rem; border-radius: .25rem; }
.signal.opportunity .signal-badge { background: var(--green);  color: #000; }
.signal.caution     .signal-badge { background: var(--yellow); color: #000; }
.signal.info        .signal-badge { background: var(--cyan);   color: #000; }
.signal-message { font-weight: 500; }
.signal-reason  { color: var(--muted); font-size: .85rem; margin-top: .2rem; }

/* recommendation */
.rec-header { display: flex; align-items: baseline; gap: .75rem; margin-bottom: 1rem; }
.rec-total { font-size: 1.4rem; font-weight: 700; color: var(--green); }
.rec-actions { display: flex; flex-direction: column; gap: .5rem; margin-bottom: 1.25rem; }
.rec-action { display: flex; justify-content: space-between; align-items: center;
              padding: .6rem 1rem; background: var(--surface);
              border: 1px solid var(--border); border-radius: .375rem; }
.rec-action .rec-class { font-weight: 500; }
.rec-action .rec-amount { font-weight: 700; color: var(--green); font-family: monospace; }
.rec-avoid, .rec-monitor { margin-top: .75rem; color: var(--muted); font-size: .85rem; }
.rec-avoid span   { color: var(--yellow); font-weight: 600; }
.rec-monitor span { color: var(--muted); }

footer { margin-top: 3rem; padding-top: 1rem; border-top: 1px solid var(--border);
         color: var(--muted); font-size: .75rem; text-align: center; }
</style>
</head>
<body>

<!-- HEADER -->
<div class="header">
  <div class="header-left">
    <h1>Investment Copilot</h1>
    <div class="subtitle">Relatório gerado em {{ date .GeneratedAt }}</div>
  </div>
  <div class="total-value">
    <div class="label">PATRIMÔNIO TOTAL</div>
    <div class="amount">R$ {{ brl .Portfolio.TotalValue }}</div>
  </div>
</div>

<!-- SUMMARY CARDS -->
<div class="cards">
  <div class="card{{ if ge .Performance.PercentReturn 0.0 }} green{{ else }} red{{ end }}">
    <div class="card-label">Retorno Total</div>
    <div class="card-value">{{ pctSigned .Performance.PercentReturn }}</div>
  </div>
  <div class="card blue">
    <div class="card-label">Total Investido</div>
    <div class="card-value">R$ {{ brl .Performance.TotalInvested }}</div>
  </div>
  {{ if .Macro }}
  <div class="card yellow">
    <div class="card-label">Selic Meta</div>
    <div class="card-value">{{ printf "%.2f%%" .Macro.SelicMeta }}</div>
  </div>
  <div class="card{{ if ge .Macro.IPCA12m 6.0 }} yellow{{ else }} green{{ end }}">
    <div class="card-label">IPCA 12m</div>
    <div class="card-value">{{ printf "%.2f%%" .Macro.IPCA12m }}</div>
  </div>
  {{ end }}
</div>

<!-- POSITIONS -->
<section>
  <h2>Posições</h2>
  <table>
    <thead>
      <tr>
        <th>Ticker</th>
        <th>Classe</th>
        <th class="right">Qtd</th>
        <th class="right">Preço Médio</th>
        <th class="right">Preço Atual</th>
        <th class="right">Valor</th>
        <th class="right">Alocação</th>
      </tr>
    </thead>
    <tbody>
      {{ range .Portfolio.Positions }}
      <tr>
        <td><strong>{{ .Ticker }}</strong></td>
        <td class="muted">{{ .Class }}</td>
        <td class="right mono">{{ printf "%.2f" .Quantity }}</td>
        <td class="right mono">R$ {{ printf "%.2f" .AvgPrice }}</td>
        <td class="right mono">R$ {{ printf "%.2f" .CurrentPrice }}</td>
        <td class="right mono">R$ {{ brl .CurrentValue }}</td>
        <td class="right mono">{{ printf "%.1f%%" .Allocation }}</td>
      </tr>
      {{ end }}
    </tbody>
  </table>
</section>

<!-- ALLOCATION VS TARGET -->
<section>
  <h2>Alocação vs Target</h2>
  <table class="alloc-table">
    <thead>
      <tr>
        <th>Classe</th>
        <th>Atual / Target</th>
        <th>Desvio</th>
        <th style="min-width:320px">Distribuição</th>
      </tr>
    </thead>
    <tbody>
      {{ range .Risk.ClassDeviations }}
      <tr>
        <td class="alloc-class">{{ .Class }}</td>
        <td class="alloc-pcts mono">{{ printf "%.1f%%" .Current }} / {{ printf "%.1f%%" .Target }}</td>
        <td class="alloc-dev {{ deviationClass .Deviation }}">{{ ppSigned .Deviation }}</td>
        <td>{{ allocationBar .Current .Target }}</td>
      </tr>
      {{ end }}
    </tbody>
  </table>
</section>

<!-- RISK ALERTS -->
{{ if .Risk.Alerts }}
<section>
  <h2>Alertas de Risco</h2>
  <div class="alerts">
    {{ range .Risk.Alerts }}
    <div class="alert {{ alertClass .Level }}">
      <span class="alert-badge">{{ upper (string .Level) }}</span>
      <span>{{ .Message }}</span>
    </div>
    {{ end }}
  </div>
</section>
{{ end }}

<!-- PERFORMANCE -->
<section>
  <h2>Performance</h2>
  <div class="perf-summary">
    <div class="perf-card">
      <div class="label">Retorno Absoluto</div>
      <div class="value {{ if ge .Performance.AbsoluteReturn 0.0 }}positive{{ else }}negative{{ end }}">
        R$ {{ brl .Performance.AbsoluteReturn }}
      </div>
      <div class="sub">desde o início</div>
    </div>
    <div class="perf-card">
      <div class="label">Retorno %</div>
      <div class="value {{ if ge .Performance.PercentReturn 0.0 }}positive{{ else }}negative{{ end }}">
        {{ pctSigned .Performance.PercentReturn }}
      </div>
      <div class="sub">sobre capital investido</div>
    </div>
    {{ if .Performance.Return12M }}
    <div class="perf-card">
      <div class="label">Retorno 12m</div>
      <div class="value {{ returnClass .Performance.Return12M }}">{{ returnStr .Performance.Return12M }}</div>
    </div>
    {{ end }}
    {{ if .Performance.Return1M }}
    <div class="perf-card">
      <div class="label">Retorno 1m</div>
      <div class="value {{ returnClass .Performance.Return1M }}">{{ returnStr .Performance.Return1M }}</div>
    </div>
    {{ end }}
  </div>

  {{ if .Benchmarks }}
  <table class="bench-table">
    <thead>
      <tr>
        <th></th>
        <th>1 mês</th>
        <th>6 meses</th>
        <th>12 meses</th>
      </tr>
    </thead>
    <tbody>
      <tr class="portfolio-row">
        <td>Carteira</td>
        <td class="{{ returnClass .Performance.Return1M }}">{{ returnStr .Performance.Return1M }}</td>
        <td class="{{ returnClass .Performance.Return6M }}">{{ returnStr .Performance.Return6M }}</td>
        <td class="{{ returnClass .Performance.Return12M }}">{{ returnStr .Performance.Return12M }}</td>
      </tr>
      {{ range .Benchmarks }}
      <tr>
        <td>{{ .Name }}</td>
        <td class="{{ returnClass .Return1M }}">{{ returnStr .Return1M }}</td>
        <td class="{{ returnClass .Return6M }}">{{ returnStr .Return6M }}</td>
        <td class="{{ returnClass .Return12M }}">{{ returnStr .Return12M }}</td>
      </tr>
      {{ end }}
    </tbody>
  </table>
  {{ end }}
</section>

<!-- MACRO CONTEXT -->
{{ if .Macro }}
<section>
  <h2>Contexto Macro</h2>
  <div class="macro-grid">
    {{ if .Macro.SelicMeta }}
    <div class="macro-item">
      <div class="m-label">Selic Meta</div>
      <div class="m-value">{{ printf "%.2f%%" .Macro.SelicMeta }}</div>
      <div class="m-sub">{{ selicTrendIcon .Macro.SelicTrend }} {{ .Macro.SelicTrendLabel }}</div>
    </div>
    {{ end }}
    {{ if .Macro.IPCA12m }}
    <div class="macro-item">
      <div class="m-label">IPCA 12m</div>
      <div class="m-value{{ if ge .Macro.IPCA12m 6.0 }} yellow{{ end }}">{{ printf "%.2f%%" .Macro.IPCA12m }}</div>
      <div class="m-sub">acumulado 12 meses</div>
    </div>
    {{ end }}
    {{ if .Macro.USDBRL }}
    <div class="macro-item">
      <div class="m-label">USD / BRL</div>
      <div class="m-value">{{ printf "R$ %.2f" .Macro.USDBRL }}</div>
    </div>
    {{ end }}
    {{ if .Macro.FocusIPCA }}
    <div class="macro-item">
      <div class="m-label">IPCA Projetado (Focus)</div>
      <div class="m-value">{{ printf "%.2f%%" (deref .Macro.FocusIPCA) }}</div>
      <div class="m-sub">mediana, fim do ano</div>
    </div>
    {{ end }}
    {{ if .Macro.FocusSelic }}
    <div class="macro-item">
      <div class="m-label">Selic Projetada (Focus)</div>
      <div class="m-value">{{ printf "%.2f%%" (deref .Macro.FocusSelic) }}</div>
      <div class="m-sub">mediana, fim do ano</div>
    </div>
    {{ end }}
  </div>
</section>
{{ end }}

<!-- PLAYBOOK SIGNALS -->
{{ if and .Playbook .Playbook.Signals }}
<section>
  <h2>Sinais do Playbook</h2>
  <div class="signals">
    {{ range .Playbook.Signals }}
    <div class="signal {{ signalClass .Level }}">
      <div class="signal-header">
        <span class="signal-badge">{{ upper (string .Level) }}</span>
        {{ if .Class }}<span class="muted" style="font-size:.8rem">{{ .Class }}</span>{{ end }}
      </div>
      <div class="signal-message">{{ .Message }}</div>
      <div class="signal-reason">{{ .Reason }}</div>
    </div>
    {{ end }}
  </div>
</section>
{{ end }}

<!-- RECOMMENDATION -->
{{ if and .Recommendation (gt .Contribution 0.0) }}
<section>
  <h2>Próximo Aporte</h2>
  <div class="rec-header">
    <span>Valor disponível:</span>
    <span class="rec-total">R$ {{ brl .Contribution }}</span>
  </div>
  <div class="rec-actions">
    {{ range .Recommendation.Actions }}
    <div class="rec-action">
      <span class="rec-class">{{ .Class }}</span>
      <span class="rec-amount">R$ {{ brl .Amount }}</span>
    </div>
    {{ end }}
  </div>
  {{ if .Recommendation.Avoid }}
  <div class="rec-avoid">
    Evitar: {{ range .Recommendation.Avoid }}<span>{{ . }}</span> {{ end }}
  </div>
  {{ end }}
  {{ if .Recommendation.Monitor }}
  <div class="rec-monitor">
    Monitorar: {{ range .Recommendation.Monitor }}<span>{{ . }}</span> {{ end }}
  </div>
  {{ end }}
</section>
{{ end }}

<footer>
  Investment Copilot · Gerado em {{ date .GeneratedAt }} · Dados de mercado podem estar defasados. Não constitui recomendação de investimento.
</footer>

</body>
</html>`
