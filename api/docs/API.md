# API — referência

Contrato HTTP da API em Go (Gin), conferido contra o código em `internal/requests/routes/`, `internal/requests/router.go` e `internal/controllers/`. Os testes de contrato estão em `internal/apitest/`.

## Convenções

- **Base:** pelo nginx, as rotas da API ficam em `/api/...` e o redirect público em `/r/...`. Os exemplos usam `http://localhost` (`PROXY_PORT=80`).
- **Formato:** JSON. Erros sempre como `{"error": "<mensagem>"}`; sucesso sem entidade como `{"message": "<mensagem>"}`.
- **Exceções ao formato:** o 401 do middleware de autenticação e o 429 do rate limiter respondem **sem corpo**.
- **Autenticação:** JWT HS256 no header `Authorization: Bearer <token>`, e só nele (token na query string é ignorado). O token expira em `TOKEN_LIFE_SPAN` minutos. Os claims `exp` e `user_id` são obrigatórios.
- **Não há logout no servidor:** o token vale até expirar. Para sair, o cliente descarta o token.
- **Rate limit:** janela fixa por IP de cliente, com até `REQUISITION_LIMIT` requisições a cada `MINUTES_OF_DURATION` minutos. Vale para **todas** as rotas, incluindo `/r/` e `/healthz`. Acima do limite, a resposta é `429`, sem corpo. Se o Redis estiver fora, a requisição passa (fail open, com log). O IP do cliente vem do `X-Forwarded-For` só quando a conexão chega de um IP em `TRUSTED_PROXIES`.
- **Erro interno:** `500 {"error":"erro interno"}`. O detalhe vai só para o log do servidor.
- **CORS:** a API não envia cabeçalhos CORS, porque o front é servido pela mesma origem.

## Resumo das rotas

| Método | Caminho | Auth | Sucesso | O que faz |
|---|---|---|---|---|
| GET | `/api/ping` | — | 200 | Verifica se o processo responde |
| POST | `/api/user` | — | 201 | Cadastra usuário e devolve token |
| POST | `/api/login` | — | 200 | Autentica e devolve token |
| PUT | `/api/user/:id` | Bearer | 200 | Atualiza nome e e-mail do próprio usuário |
| DELETE | `/api/user/:id` | Bearer | 200 | Exclui o próprio usuário e os links dele |
| GET | `/api/links` | Bearer | 200 | Lista os links do usuário |
| POST | `/api/links` | Bearer | 201 | Cria link curto |
| GET | `/api/links/:slug` | Bearer | 200 | Detalhe de um link do usuário |
| POST | `/api/links/:slug/deactivate` | Bearer | 200 | Desativa um link do usuário (sem volta) |
| GET | `/r/:slug` | — | 302 | Redirect público para o destino |
| GET | `/healthz` | — | 200 | Healthcheck do container (não passa pelo nginx) |

Nas rotas com Bearer, token ausente, inválido ou expirado dá `401`, sem corpo (middleware).

---

## Saúde

### `GET /api/ping`
`200 {"message":"pong"}`. Não toca no banco.

### `GET /healthz`
Usada pelo healthcheck do Docker, só dentro da rede do compose. Pelo nginx, `/healthz` cai no `index.html` do front.
- `200 {"status":"ok","db":"ok"}`: o ping no Postgres respondeu em até 1 s.
- `503 {"status":"indisponível","db":"erro"}`: o Postgres não respondeu.
- O Redis não entra na conta, porque a API funciona sem ele.

---

## Usuários

### `POST /api/user` — cadastro
```json
{ "full_name": "Ana", "email": "ana@exemplo.com", "password": "senha-forte-123" }
```
| Status | Quando |
|---|---|
| 201 | `{"message":"User created successfully","token":"<jwt>"}` |
| 409 | `{"error":"e-mail já cadastrado"}` |
| 422 | campo ausente, e-mail inválido, senha com menos de 8 caracteres ou mais de 72 **bytes** (limite do bcrypt) |
| 400 | `{"error":"Error generating token"}`: falha ao assinar o token. O usuário já foi criado. Inconsistência conhecida: deveria ser 500 |

### `POST /api/login`
```json
{ "email": "ana@exemplo.com", "password": "senha-forte-123" }
```
| Status | Quando |
|---|---|
| 200 | `{"token":"<jwt>"}` |
| 401 | `{"error":"Invalid email or password"}`, igual para senha errada e e-mail inexistente. O tempo de resposta também é igual: o e-mail inexistente compara a senha com um hash fictício |
| 400 | corpo inválido. Inconsistência conhecida: as outras rotas usam 422 |
| 400 | `{"error":"Error generating token"}`. Inconsistência conhecida: deveria ser 500 |

### `PUT /api/user/:id` — atualizar (Bearer)
```json
{ "full_name": "Ana Souza", "email": "ana.souza@exemplo.com" }
```
Os dois campos são obrigatórios. A posse é checada **antes** do corpo.

| Status | Quando |
|---|---|
| 200 | `{"message":"usuário atualizado"}` |
| 403 | `:id` não é o usuário do token, qualquer que seja o corpo |
| 422 | `:id` inválido ou corpo inválido |
| 409 | e-mail já usado por outro usuário |
| 404 | usuário não existe mais |

### `DELETE /api/user/:id` — excluir (Bearer)
| Status | Quando |
|---|---|
| 200 | `{"message":"usuário excluído"}`. Os links do usuário são apagados em cascata (`ON DELETE CASCADE`) e marcados como inexistentes no cache |
| 403 | `:id` não é o usuário do token |
| 404 | usuário não existe |
| 422 | `:id` inválido |

Um token de usuário excluído continua passando no middleware até expirar. As rotas respondem 404 ou lista vazia nesse caso.

---

## Links

Objeto `link` nas respostas:
```json
{
  "slug": "aZ3k9Qp",
  "url": "http://localhost/r/aZ3k9Qp",
  "target_url": "https://exemplo.com/pagina",
  "is_active": true,
  "last_change": "2026-11-20T13:47:08.123-03:00"
}
```
`url` é montada a partir de `PUBLIC_BASE_URL` (configuração), nunca do header `Host`.

### `POST /api/links` — criar (Bearer)
```json
{ "target_url": "https://exemplo.com/pagina" }
```
- **`target_url`:** só `http` ou `https`, com host, sem credenciais embutidas (`user@host`) e com no máximo 2048 bytes.
- **Código curto:** 7 caracteres base62 aleatórios (`crypto/rand`). Cada chamada cria um link novo, mesmo para um destino repetido.
- **Colisão:** se o código gerado já existir, a API tenta outro, até 5 vezes.

| Status | Quando |
|---|---|
| 201 | `{"id":1,"slug":"aZ3k9Qp","url":"http://localhost/r/aZ3k9Qp","owner":1}` |
| 422 | URL ausente ou inválida |
| 500 | 5 colisões seguidas ou falha do gerador |

### `GET /api/links` — listar (Bearer)
`200 {"links":[<link>, ...]}`, só com os links do usuário do token (lista vazia se não houver).

### `GET /api/links/:slug` — detalhe (Bearer)
- `200 {"link": <link>}`.
- `404 {"error":"link não encontrado"}`: também para link de **outro usuário**, para não revelar que o código existe.

### `POST /api/links/:slug/deactivate` — desativar (Bearer)
- `200 {"message":"link desativado"}`.
- `404`: código inexistente ou de outro usuário.
- Não existe rota de reativação.

---

## Redirect público

### `GET /r/:slug`
| Status | Quando |
|---|---|
| 302 | link ativo, com `Location: <target_url>`. É 302 e não 301 porque o navegador guarda o 301 sem prazo, e a desativação não teria efeito para quem já visitou |
| 404 | `{"error":"link não encontrado"}` |
| 410 | `{"error":"link desativado"}` |

**Cache:** a resposta usa cache-aside no Redis, com a chave `link:<slug>`.
- **TTL:** `LINK_CACHE_TTL_SECONDS` (padrão 600) para link existente e `LINK_CACHE_NEGATIVE_TTL_SECONDS` (padrão 30) para código inexistente.
- **Atualização:** criar, desativar e excluir o dono atualizam a entrada na hora.
- **Redis fora ou lento:** acima de 100 ms por operação, a consulta vai ao banco.

---

## Exemplo com curl

```bash
BASE=http://localhost
TOKEN=$(curl -s -X POST $BASE/api/user -H 'Content-Type: application/json' \
  -d '{"full_name":"Ana","email":"ana@exemplo.com","password":"senha-forte-123"}' | jq -r .token)
SLUG=$(curl -s -X POST $BASE/api/links -H 'Content-Type: application/json' \
  -H "Authorization: Bearer $TOKEN" -d '{"target_url":"https://exemplo.com"}' | jq -r .slug)
curl -si $BASE/r/$SLUG | grep -iE '^(HTTP|location)'
```
