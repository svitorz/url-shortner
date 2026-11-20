#!/usr/bin/env bash
# Roda gofmt, go vet e go test da API com os testes de integração ligados.
#
# Sobe um Postgres 18 descartável (container --rm, porta só em 127.0.0.1),
# separado do banco do compose: os dados de desenvolvimento não são tocados.
# Cada teste cria e apaga o próprio banco dentro dele (internal/testdb).
#
# Uso: scripts/test.sh [args extras do go test, ex.: -run Redirect -v]
# Variáveis: TEST_PG_PORT (55432).
set -euo pipefail

PORT=${TEST_PG_PORT:-55432}
NAME=url-shortner-testdb
cd "$(dirname "$0")/../api"

docker run -d --rm --name "$NAME" -e POSTGRES_PASSWORD=test -e POSTGRES_DB=test \
  -p "127.0.0.1:$PORT:5432" postgres:18 >/dev/null
trap 'docker stop "$NAME" >/dev/null' EXIT

for _ in $(seq 1 30); do
  docker exec "$NAME" pg_isready -U postgres -d test >/dev/null 2>&1 && break
  sleep 1
done

unformatted=$(gofmt -l .)
if [[ -n $unformatted ]]; then
  echo "gofmt: arquivos fora do formato:"; echo "$unformatted"; exit 1
fi
go vet ./...
TEST_DATABASE_DSN="host=127.0.0.1 port=$PORT user=postgres password=test dbname=test sslmode=disable" \
  go test -count=1 "$@" ./...
