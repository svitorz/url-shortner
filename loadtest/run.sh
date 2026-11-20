#!/usr/bin/env bash
# Roda o teste de carga do redirect pelo nginx, N rodadas, e guarda os JSON.
#
# Uso (na raiz do repo, stack no ar):
#   PROXY_PORT=8088 RUN_LABEL=baseline loadtest/run.sh
# Variáveis: RUNS (3), VUS (50), WARMUP (10s), DURATION (60s), PROXY_PORT (80).
#
# O que faz:
#   1. sobe a API com o override de rate limit alto (loadtest/compose.loadtest.yml);
#   2. cria um usuário e um link pelo nginx do host (setup, fora da medição);
#   3. roda o k6 num container na rede do compose, batendo em http://proxy
#      (o nginx), sem passar pelo port-forward do Docker no host;
#   4. grava loadtest/results/<RUN_LABEL>-<n>.json;
#   5. devolve a API à config normal (sem o override).
# Requer: bash, curl, jq, docker compose. Não lê o .env.
set -euo pipefail

PROXY_PORT=${PROXY_PORT:-80}
HOST_URL="http://localhost:${PROXY_PORT}"
RUNS=${RUNS:-3}
VUS=${VUS:-50}
WARMUP=${WARMUP:-10s}
DURATION=${DURATION:-60s}
RUN_LABEL=${RUN_LABEL:?defina RUN_LABEL (ex.: baseline, cache)}
K6_IMAGE=grafana/k6:2.4.0
NETWORK=url-shortner_url-shortner-network
OUT_DIR="$(pwd)/loadtest/results"
COMPOSE=(docker compose -f docker-compose.yml -f loadtest/compose.loadtest.yml)

restore() { docker compose up -d --wait api >/dev/null 2>&1 || true; }
trap restore EXIT

echo "== subindo a API com rate limit alto"
"${COMPOSE[@]}" up -d --wait api >/dev/null

echo "== setup: usuário e link de teste"
email="load-$(date +%s%N)@loadtest.test"
token=$(curl -sf -X POST "$HOST_URL/api/user" -H 'Content-Type: application/json' \
  -d "{\"full_name\":\"Load\",\"email\":\"$email\",\"password\":\"senha-forte-123\"}" | jq -r .token)
slug=$(curl -sf -X POST "$HOST_URL/api/links" -H 'Content-Type: application/json' \
  -H "Authorization: Bearer $token" -d '{"target_url":"https://example.com/loadtest"}' | jq -r .slug)
code=$(curl -s -o /dev/null -w '%{http_code}' "$HOST_URL/r/$slug")
[[ $code == 302 ]] || { echo "redirect do slug de teste deu $code"; exit 1; }
echo "   slug=$slug"

mkdir -p "$OUT_DIR"
for i in $(seq 1 "$RUNS"); do
  echo "== rodada $i/$RUNS ($VUS VUs, aquecimento $WARMUP, medição $DURATION)"
  docker run --rm --network "$NETWORK" -u "$(id -u):$(id -g)" \
    -v "$(pwd)/loadtest:/scripts:ro" -v "$OUT_DIR:/out" \
    -e BASE_URL=http://proxy -e SLUG="$slug" -e VUS="$VUS" -e WARMUP="$WARMUP" \
    -e DURATION="$DURATION" -e RUN_NAME="$RUN_LABEL-$i" \
    "$K6_IMAGE" run --quiet /scripts/redirect.js
done
echo "== resultados em loadtest/results/$RUN_LABEL-*.json"
