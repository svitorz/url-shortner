// Teste de carga do redirect público GET /r/<slug> (Fase 4).
//
// Modelo fechado: VUS usuários virtuais em loop, sem pausa. Um cenário de
// aquecimento (não medido) seguido do cenário medido; os números do relatório
// vêm só das métricas com a tag scenario:measure.
//
// Não segue o redirect (redirects: 0): mede a nossa resposta 302, não o
// site de destino.
//
// Variáveis: BASE_URL (ex.: http://proxy), SLUG (obrigatória), VUS, WARMUP, DURATION.
import http from 'k6/http';
import { check } from 'k6';

const BASE_URL = __ENV.BASE_URL || 'http://proxy';
const SLUG = __ENV.SLUG;
const VUS = parseInt(__ENV.VUS || '50', 10);
const WARMUP = __ENV.WARMUP || '10s';
const DURATION = __ENV.DURATION || '60s';

if (!SLUG) {
  throw new Error('defina SLUG (código de um link ativo)');
}

export const options = {
  discardResponseBodies: true,
  summaryTrendStats: ['avg', 'min', 'med', 'p(90)', 'p(95)', 'p(99)', 'max'],
  scenarios: {
    warmup: { executor: 'constant-vus', vus: VUS, duration: WARMUP, gracefulStop: '0s' },
    measure: { executor: 'constant-vus', vus: VUS, duration: DURATION, startTime: WARMUP, gracefulStop: '5s' },
  },
  // Thresholds nas submétricas do cenário medido fazem o k6 incluí-las no resumo.
  thresholds: {
    'http_req_duration{scenario:measure}': ['max>=0'],
    'http_reqs{scenario:measure}': ['count>=0'],
    'http_req_failed{scenario:measure}': ['rate>=0'],
    'checks{scenario:measure}': ['rate>=0'],
  },
};

const url = `${BASE_URL}/r/${SLUG}`;
const params = { redirects: 0 };

export default function () {
  const res = http.get(url, params);
  check(res, {
    'status 302': (r) => r.status === 302,
    'tem Location': (r) => !!r.headers['Location'],
  });
}

export function handleSummary(data) {
  const m = (name) => data.metrics[name] && data.metrics[name].values;
  const d = m('http_req_duration{scenario:measure}');
  const reqs = m('http_reqs{scenario:measure}');
  const failed = m('http_req_failed{scenario:measure}');
  const checks = m('checks{scenario:measure}');
  const seconds = parseFloat(DURATION);
  const out = {
    vus: VUS,
    duration: DURATION,
    requests: reqs.count,
    rps: reqs.count / seconds,
    latency_ms: { p50: d.med, p95: d['p(95)'], p99: d['p(99)'], avg: d.avg, max: d.max },
    http_req_failed_rate: failed.rate,
    checks_ok_rate: checks.rate,
  };
  return {
    stdout: JSON.stringify(out, null, 2) + '\n',
    [`/out/${__ENV.RUN_NAME || 'run'}.json`]: JSON.stringify(out, null, 2),
  };
}
