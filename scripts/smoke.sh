#!/usr/bin/env bash
# Teste de fumaça: exercita todas as rotas da API pelo nginx.
#
# Uso:  scripts/smoke.sh                 (BASE_URL padrão: http://localhost:$PROXY_PORT)
#       BASE_URL=http://localhost:8000 scripts/smoke.sh   (direto na API)
#       SMOKE_RESTART=0 scripts/smoke.sh                  (pula o teste de reinício)
#
# Cada verificação afirma o comportamento CORRETO esperado (o contrato alvo),
# com a hipótese do PROGRESS.md entre colchetes. Onde a decisão ainda não foi
# tomada, aceita-se um conjunto de status (ex.: "302|307").
# Requer: bash, curl, jq. Sai com código 1 se alguma verificação falhar.
set -uo pipefail

PROXY_PORT=${PROXY_PORT:-80}
BASE_URL=${BASE_URL:-http://localhost:${PROXY_PORT}}
API="$BASE_URL/api"
RUN=$(date +%s%N)
H='Content-Type: application/json'
BODY=$(mktemp)
trap 'rm -f "$BODY"' EXIT

pass=0
fail=0

# req METHOD URL [args curl...] -> define $code e deixa o corpo em $BODY
req() {
  local m=$1 u=$2
  shift 2
  code=$(curl -s -o "$BODY" -w '%{http_code}' -X "$m" "$u" "$@")
}

# check "descrição" "esperado (regex, ex.: 200|201)" "obtido"
check() {
  if [[ $3 =~ ^($2)$ ]]; then
    pass=$((pass + 1))
    printf 'PASS  %s\n' "$1"
  else
    fail=$((fail + 1))
    printf 'FAIL  %s (esperado %s, obtido %s) %s\n' "$1" "$2" "$3" "$(head -c 160 "$BODY" | tr '\n' ' ')"
  fi
}

json() { jq -r "$1 // empty" "$BODY" 2>/dev/null; }
auth() { printf 'Authorization: Bearer %s' "$1"; }

register() { # email -> token
  req POST "$API/user" -H "$H" -d "{\"full_name\":\"Smoke\",\"email\":\"$1\",\"password\":\"senha-forte-123\"}"
}
login() { # email senha
  req POST "$API/login" -H "$H" -d "{\"email\":\"$1\",\"password\":\"$2\"}"
}
create_link() { # token url
  req POST "$API/links" -H "$H" -H "$(auth "$1")" -d "{\"target_url\":\"$2\"}"
}

EA="a-$RUN@smoke.test"
EB="b-$RUN@smoke.test"
TARGET="https://example.com/smoke/$RUN"

echo "== base: $API"

# ---------- saúde ----------
req GET "$API/ping"
check "ping" "200" "$code"

# ---------- front (nginx serve o build do Vite) ----------
ctype() { curl -s -o "$BODY" -w '%{http_code} %{content_type}' "$@"; }
r=$(ctype "$BASE_URL/")
check "[H16] / serve o index.html do front" "200 text/html.*" "$r"
r=$(ctype "$BASE_URL/links")
check "[H16] rota da SPA (/links) cai no index.html (try_files)" "200 text/html.*" "$r"
asset=$(grep -o '/assets/[^"]*\.js' "$BODY" | head -1)
r=$(curl -s -o /dev/null -H 'Accept-Encoding: gzip' -w '%{http_code} %header{content-encoding}' "$BASE_URL${asset:-/assets/nao-achou.js}")
check "[H18] JS do front com gzip" "200 gzip" "$r"
r=$(ctype "$BASE_URL/assets/nao-existe.js")
check "[H16] asset inexistente -> 404 (não o index.html)" "404 .*" "$r"

# ---------- usuários ----------
register "$EA"
check "cadastro A" "201" "$code"
register "$EA"
check "[H9] cadastro com e-mail duplicado -> 409" "409" "$code"
register "$EB"
check "cadastro B" "201" "$code"
req POST "$API/user" -H "$H" -d "{\"full_name\":\"X\",\"email\":\"c-$RUN@smoke.test\",\"password\":\"$(printf 'a%.0s' {1..100})\"}"
check "[N4] senha > 72 bytes -> 400/422" "400|422" "$code"

login "$EA" "senha-forte-123"
check "login A" "200" "$code"
TA=$(json .token)
login "$EB" "senha-forte-123"
TB=$(json .token)
login "$EA" "senha-errada-123"
check "login com senha errada -> 401" "401" "$code"
login "nao-existe-$RUN@smoke.test" "senha-errada-123"
check "[N3] login com e-mail inexistente -> 401 (igual à senha errada)" "401" "$code"

# ---------- auth ----------
req GET "$API/links"
check "sem token -> 401" "401" "$code"
req GET "$API/links" -H "Authorization: Bearer token-invalido"
check "token inválido -> 401" "401" "$code"
req GET "$API/links?token=$TA"
check "[H14] token na query string -> 401" "401" "$code"

# ---------- links ----------
create_link "$TA" "$TARGET"
check "A cria link" "201" "$code"
SLUG_A=$(json .slug)
SHORT_A=$(json .url)

if [[ -n $SHORT_A ]]; then
  c=$(curl -s -o /dev/null -w '%{http_code}' "$SHORT_A")
  loc=$(curl -s -o /dev/null -w '%{redirect_url}' "$SHORT_A")
else
  c=none loc=
fi
check "[H3] URL curta devolvida redireciona" "302|307" "$c"
check "[H3] URL curta aponta para o destino" "$(printf '%s' "$TARGET" | sed 's/[.]/[.]/g')" "$loc"

create_link "$TB" "$TARGET"
check "[H2] B cria link com a mesma URL" "201" "$code"
SLUG_B=$(json .slug)
[[ -n $SLUG_B && $SLUG_B != "$SLUG_A" ]] && same=diferente || same=igual
check "[H2] B não recebe o slug de A" "diferente" "$same"
req GET "$API/links" -H "$(auth "$TB")"
check "[H2] B lista os próprios links" "200" "$code"
check "[H2] lista de B contém 1 link" "1" "$(json '.links | length')"

LONG="https://example.com/$(head -c 2100 /dev/zero | tr '\0' a)"
create_link "$TA" "$LONG"
check "[H2/H4] URL gigante (2100+) -> 400/422" "400|422" "$code"
create_link "$TA" "javascript://example.com/%0aalert(1)"
check "[H4] esquema javascript: -> 400/422" "400|422" "$code"
create_link "$TA" "ftp://example.com/arquivo"
check "[H4] esquema ftp: -> 400/422" "400|422" "$code"
create_link "$TA" ""
check "[H4] URL vazia -> 400/422" "400|422" "$code"

req GET "$API/links" -H "$(auth "$TA")"
check "A lista links" "200" "$code"

req GET "$API/links/$SLUG_A" -H "$(auth "$TA")"
check "[H6] detalhe pelo dono -> 200" "200" "$code"
req GET "$API/links/$SLUG_A" -H "$(auth "$TB")"
check "[H6] detalhe por outro usuário -> 403/404" "403|404" "$code"
req GET "$API/links/nao-existe-$RUN" -H "$(auth "$TA")"
check "[H6] detalhe de slug inexistente -> 404" "404" "$code"

req GET "$BASE_URL/r/$SLUG_A"
check "[H5] redirect de slug ativo -> 302/307" "302|307" "$code"
req GET "$BASE_URL/r/nao-existe-$RUN"
check "[H5] redirect de slug inexistente -> 404" "404" "$code"

req POST "$API/links/$SLUG_A/deactivate" -H "$(auth "$TB")"
check "desativar link de outro usuário -> 403/404" "403|404" "$code"
req POST "$API/links/nao-existe-$RUN/deactivate" -H "$(auth "$TA")"
check "desativar slug inexistente -> 404" "404" "$code"
req POST "$API/links/$SLUG_A/deactivate" -H "$(auth "$TA")"
check "dono desativa link -> 200" "200" "$code"
req GET "$BASE_URL/r/$SLUG_A"
check "[H5] redirect de link inativo -> 404/410" "404|410" "$code"

# ---------- atualizar usuário ----------
# o id do usuário vem do JWT (claim user_id)
uid() {
  local p
  p=$(cut -d. -f2 <<<"$1" | tr '_-' '/+')
  while (( ${#p} % 4 )); do p+='='; done
  base64 -d <<<"$p" 2>/dev/null | jq -r .user_id 2>/dev/null
}
UA=$(uid "$TA")
UB=$(uid "$TB")
req PUT "$API/user/$UA" -H "$H" -H "$(auth "$TA")" -d "{\"full_name\":\"A\",\"email\":\"$EB\"}"
check "[H9] atualizar para e-mail de outro -> 409" "409" "$code"
req PUT "$API/user/$UA" -H "$H" -H "$(auth "$TA")" -d "{\"full_name\":\"A atualizado\",\"email\":\"$EA\"}"
check "[H9] atualizar usuário -> 200" "200" "$code"
req PUT "$API/user/$UB" -H "$H" -H "$(auth "$TA")" -d "{\"full_name\":\"hack\",\"email\":\"$EA\"}"
check "atualizar outro usuário -> 403" "403" "$code"

# ---------- logout ----------
login "$EA" "senha-forte-123"
TL=$(json .token)
req POST "$API/logout" -H "$(auth "$TL")"
if [[ $code == 404 ]]; then
  check "[H15] logout removido (404)" "404" "$code"
else
  check "[H15] logout responde 200/204" "200|204" "$code"
  req GET "$API/links" -H "$(auth "$TL")"
  check "[H15] token após logout -> 401" "401" "$code"
fi

# ---------- deletar usuário ----------
req DELETE "$API/user/$UA" -H "$(auth "$TA")"
check "[H8] deletar usuário que tem links -> 200/204" "200|204" "$code"
login "$EA" "senha-forte-123"
check "[H8] login do usuário deletado -> 401" "401" "$code"

# ---------- persistência ----------
if [[ ${SMOKE_RESTART:-1} != 0 ]]; then
  docker compose restart api >/dev/null 2>&1
  for _ in $(seq 1 30); do
    [[ $(curl -s -o /dev/null -w '%{http_code}' "$API/ping") == 200 ]] && break
    sleep 1
  done
  login "$EB" "senha-forte-123"
  check "[H1] dados sobrevivem ao reinício da API (login de B)" "200" "$code"
fi

echo "== $pass passaram, $fail falharam"
[[ $fail -eq 0 ]]
