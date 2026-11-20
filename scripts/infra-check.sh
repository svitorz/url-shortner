#!/usr/bin/env bash
# Verificações da infra do compose (Fase 3). Rode na raiz do repo com o stack
# no ar (docker compose up -d --build). Não lê nem imprime o .env.
# Requer: bash, docker compose, jq. Sai com código 1 se algo falhar.
set -uo pipefail

pass=0
fail=0
check() { # "descrição" "esperado (regex)" "obtido"
  if [[ $3 =~ ^($2)$ ]]; then
    pass=$((pass + 1)); printf 'PASS  %s\n' "$1"
  else
    fail=$((fail + 1)); printf 'FAIL  %s (esperado %s, obtido %s)\n' "$1" "$2" "$3"
  fi
}

# Portas publicadas no host de cada serviço ("" = nenhuma).
published() {
  docker compose ps --format json "$1" 2>/dev/null |
    jq -rs '[.[].Publishers[]? | select(.PublishedPort > 0) | "\(.URL):\(.PublishedPort)"] | unique | join(",")'
}

check "db sem porta publicada no host" "" "$(published db)"
check "redis sem porta publicada no host" "" "$(published redis)"
check "api publicada só no loopback" "127\\.0\\.0\\.1:[0-9]+" "$(published api)"

for s in db redis api; do
  h=$(docker inspect -f '{{if .State.Health}}{{.State.Health.Status}}{{else}}sem-healthcheck{{end}}' \
    "$(docker compose ps -q "$s")" 2>/dev/null)
  check "$s healthy" "healthy" "$h"
done

check "api roda como usuário não-root" "[1-9][0-9]*" "$(docker compose exec -T api id -u 2>/dev/null | tr -d '\r')"
check "serviço web (http-server) removido" "" "$(docker compose ps -q web 2>/dev/null)"

r=$(docker compose exec -T redis sh -c 'unset REDISCLI_AUTH; redis-cli ping' 2>&1 | tr -d '\r')
check "redis recusa comando sem senha" "NOAUTH.*" "$r"

# Sem a variável, o compose deve recusar subir com mensagem clara.
for v in REDIS_PASSWORD DB_PASSWORD API_KEY; do
  out=$(env "$v=" docker compose config 2>&1 >/dev/null)
  check "compose recusa $v vazia" ".*$v.*obrigat.*" "$out"
done

echo "== $pass passaram, $fail falharam"
[[ $fail -eq 0 ]]
