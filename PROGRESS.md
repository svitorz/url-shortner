# PROGRESS — url-shortner

> Memória do trabalho. Ler ANTES de qualquer coisa em uma sessão nova ou após compactação de contexto.

## Fase atual

**Fase 2 — correções da API: concluída.** Aguardando "continuar" para a Fase 3.
Branch: `fix/api-stabilize` (criado a partir de `main` @ `f15ea2a`), com as fases 0–2. Nada enviado ao remoto.

## Checklist

| Fase | Item | Status |
|---|---|---|
| 0 | Ler api/, web/, proxy/, compose, docs | ✅ |
| 0 | Verificar H1–H22 (estático + execução contra Postgres/Redis descartáveis) | ✅ |
| 0 | Branch `fix/api-stabilize` + `PROGRESS.md` | ✅ |
| 1 | Corrigir H1 com `AutoMigrate` (`9fa93a3`) | ✅ provado |
| 1 | Volume persistente do Postgres (`394530c`) | ✅ provado |
| 1 | Rede: web/proxy na rede da API (N1, `a8d7dc1`) | ✅ |
| 1 | Postgres/Redis sem porta no host, `PROXY_PORT`, healthcheck do Redis (`c7b0553`) | ✅ |
| 1 | `.env.example` sem comentários inline (N2), `TOKEN_LIFE_SPAN=60` (`7166991`) | ✅ |
| 1 | Script de fumaça `scripts/smoke.sh` (`d539500`) | ✅ |
| 1 | Config/DB/Redis montados no `main`, sem `init()` (`f9d3a97`) | ✅ fumaça idêntica |
| 2 | H2 código curto aleatório + dono correto (`b805f4d`) | ✅ fumaça + `shortcode` tests |
| 2 | H4 validação de URL (`57b43d9`) | ✅ `TestValidateTargetURL` (17 casos) + fumaça |
| 2 | H3 `/r/:slug` público + `PUBLIC_BASE_URL` + nginx `location /r/` (`c10c1e5`) | ✅ fumaça pelo nginx |
| 2 | H7 tags `binding` (`8f844a8`) | ✅ `go vet` limpo |
| 2 | H5 redirect 302 / 404 / 410 (`073986d`) | ✅ fumaça |
| 2 | H6 detalhe/desativar 200/404 (`8e84ce6`) | ✅ fumaça |
| 2 | H8 `ON DELETE CASCADE` + ajuste idempotente da FK (`20103c6`) | ✅ teste com Postgres real + fumaça |
| 2 | H9 409 em e-mail duplicado, update 200 (`9589ddf`) | ✅ fumaça |
| 2 | N3 login sem enumeração (`7b3ac38`) | ✅ fumaça + tempo medido |
| 2 | N4 senha > 72 bytes → 422 (`65c1c19`) | ✅ fumaça |
| 2 | H10 janela em minutos (`b65e1bc`) | ✅ teste + TTL=58 no Redis real |
| 2 | H11 TTL só na 1ª requisição + H12 prefixo (`539772d`) | ✅ 6 testes miniredis |
| 2 | H12 `TRUSTED_PROXIES`, nginx sobrescreve XFF, API em 127.0.0.1 (`f00ccdd`) | ✅ 3 testes + verificação no compose |
| 2 | H14 logs sem token/hash/e-mail (`a01daf8`) | ✅ grep nos logs após fumaça: 0 |
| 2 | H14 JWT só Bearer, HS256, `exp` obrigatório (`e4e2419`) | ✅ 9 casos em `TestRejectsInvalidTokens` |
| 2 | H14 `API_KEY` ≥ 32 bytes (`e10dc8b`) | ✅ teste |
| 2 | H15 remover logout (`3cd6479`), `go mod tidy` (`5fe4388`) | ✅ fumaça, build/vet/test |
| 2 | Teste do loop de retry do código curto (colisão forçada) | ⏳ Fase 5 (precisa de Postgres de teste no CI) |
| 3–7 | ver plano | ⏳ |

## Como subir (estado atual)

```bash
cp .env.example .env   # preencher DB_PASSWORD, REDIS_PASSWORD, API_KEY (openssl rand -hex 32), PUBLIC_BASE_URL
docker compose up -d --build
PROXY_PORT=80 scripts/smoke.sh   # use o mesmo PROXY_PORT do .env
```

Testes Go (de `api/`): `go test ./...`. O teste de migração usa Postgres real e é pulado sem
`TEST_DATABASE_DSN` (ex.: container descartável `postgres:18` em `127.0.0.1:55432`).

Na máquina de desenvolvimento as portas 80 e 8080 estão ocupadas por outros projetos; o `.env` local usa `PROXY_PORT=8088`.
O front em `/` ainda dá **502** (H16, Fase 3).

## Como a Fase 0 foi verificada

- `go vet ./...` e `go build` com Go 1.24.4 (máquina local, Debian, kernel 6.12).
- Binário da API rodado no host contra containers descartáveis `postgres:18` e `redis` (portas 127.0.0.1:55432/56379, removidos ao final). O stack do projeto **não** foi subido (subir a API apagaria qualquer banco existente — H1).
- Script de sondagem com curl (scratchpad, não versionado) exercitando todas as rotas.
- `docker compose config` rodado numa cópia do compose com `.env` = `.env.example` (no scratchpad), só para ver a interpolação.
- Nenhum `.env` existe no repositório; histórico do git só contém `.env.example` (sem segredos reais).

## Resultado das rotas na Fase 0 (código de `main`, limiter alto) — histórico

| Caso | Resultado |
|---|---|
| GET /api/ping | 200 |
| POST /api/user (novo / e-mail duplicado / senha de 100 chars) | 201 / 400 / **500** |
| POST /api/login (senha errada / e-mail inexistente) | 401 / **400** (status difere → enumeração de usuários) |
| POST /api/links (URL nova) | 201, `url` = `http://host/r/<slug>` → **404** |
| POST /api/links (mesma URL, outro usuário) | **201 com o link do usuário A** (`owner: 1`) |
| POST /api/links `javascript://…`, `ftp://…` | **201** |
| POST /api/links URL de 220 chars | **400** com erro cru do Postgres (`varchar(64)`) |
| GET /api/links/:slug (dono) | **302** com JSON |
| GET /api/links/:slug (outro usuário / inexistente) | **401** / **401** |
| GET /api/r/:slug (ativo / inexistente / inativo) | **301** / **422** / 422 |
| GET /api/links?token=… | 200 (token na query aceito) |
| POST /api/links/:slug/deactivate (outro usuário / inexistente / dono) | 401 / 401 / 200 |
| PUT /api/user/:id (e-mail de outro / ok / id de outro) | **500** / **202 string solta** / 403 |
| POST /api/logout | 200 vazio (não faz nada) |
| DELETE /api/user/:id com links | **500** expondo struct do erro do Postgres |
| Reinício da API | **todos os dados apagados** (2 usuários → 0) |

## Script de fumaça — resultado atual

`PROXY_PORT=8088 scripts/smoke.sh`, pelo nginx, stack do compose, 2026-11-20, fim da Fase 2:
**39 passaram, 0 falharam** (o caminho de logout removido tem 1 verificação a menos que o de logout existente).
`go vet ./...` limpo; `go test ./...` com 7 pacotes ok (com `TEST_DATABASE_DSN` definido).

### Histórico: linha de base da Fase 1

**17 passaram, 23 falharam.** Igual antes e depois do refactor (diff normalizado vazio).

Passam: ping, cadastros, login, senha errada 401, sem token/token inválido 401, criar link, B cria com mesma URL (201), listar, URL vazia, URL gigante (*por acidente*: 400 com erro cru do `varchar(64)`), desativar como dono, atualizar outro usuário 403, logout 200, **persistência após reinício (H1)**.

Falham (todas esperadas, item → fase 2): H9 duplicado 400≠409; N4 senha longa 500; N3 e-mail inexistente 400; H14 `?token=` 200; H3 URL curta (e sem porta: nginx repassa `Host $host`); H2 B recebe slug de A, lista de B vazia; H4 `javascript:`/`ftp:` 201; H6 detalhe 302/401/401; H5 redirect 301 / inexistente 422 / inativo 422; desativar alheio/inexistente 401; H9 update duplicado 500 e sucesso 202; H15 token válido após logout; H8 delete 500 e login do "deletado" 200.

Persistência do volume: container `db` recriado com `--force-recreate`, `users` 4 → 4.

## Problemas verificados (resumo; detalhes no relatório da Fase 0)

Críticos: H1 (DROP em todo boot), segredo JWT forjável (`API_KEY` vazio aceito; com o `.env.example` via compose o segredo vira o texto do comentário), H2 (link de outro dono), H8 (delete 500), proxy fora da rede da API.
Ver relatório completo na conversa da Fase 0; cada item será tratado na fase correspondente.

## Log de decisões

| Data | Decisão | Porquê |
|---|---|---|
| 2026-11-19 | Não subir o stack do projeto na Fase 0 | H1 apagaria dados de um banco existente; verificação feita em containers descartáveis |
| 2026-11-19 | **Migrations: `AutoMigrate`** (decisão do usuário) | Simples, nunca apaga, zero dependência nova. Descartado: goose/golang-migrate (histórico e rollback, mas +1 dependência e schema duplicado struct/SQL). Trade-off aceito: AutoMigrate não remove/renomeia colunas; mudança destrutiva de schema exige SQL manual. Migrar para goose depois é barato |
| 2026-11-19 | Falha de migration encerra o processo | Seguir com schema em estado desconhecido gera erros difíceis de diagnosticar |
| 2026-11-19 | Postgres/Redis sem porta publicada no host (antecipado da Fase 3) | Necessário para subir nesta máquina (6379 ocupada) e não há motivo para expor. Para acesso local: `ports: ["127.0.0.1:5432:5432"]` |
| 2026-11-19 | `PROXY_PORT` com padrão 80 | Porta 80 ocupada nesta máquina; mantém o padrão para quem não tem conflito |
| 2026-11-19 | Healthcheck do Redis + `depends_on: service_healthy` (antecipado da Fase 3) | A API passou a falhar cedo se o Redis não responder; sem isso a ordem de boot seria uma corrida. `REDISCLI_AUTH` evita senha na linha de comando do healthcheck |
| 2026-11-19 | `TOKEN_LIFE_SPAN=60` no exemplo | O valor é lido em minutos; 3600 dava tokens de 60 h |
| 2026-11-19 | Script de fumaça em bash+curl+jq afirmando o contrato ALVO | Cada falha atual vira a prova de um conserto da Fase 2. Onde a decisão ainda está aberta, aceita um conjunto (`302\|307`, `403\|404`). Descartado: atualizar `app.http.example` (manual, não dá pass/fail) |
| 2026-11-19 | **Injeção explícita de dependências** em vez de globais | `config.Load` valida tudo e lista o que falta; `main.run` conecta e monta `Controller{DB,Tokens,HashCost}` e `RateLimiter`. Alternativa descartada: manter `repository.DB`/`RDB` globais, só inicializados pelo `main` (diff menor, mas testes paralelos compartilhariam estado global e a dependência continuaria implícita). Sem interfaces: os testes da Fase 2/5 vão usar Postgres/Redis reais |
| 2026-11-19 | `API_KEY` vazia agora impede o boot | Consequência direta do "falhar cedo se algo faltar"; fecha a forja de JWT da Fase 0. A validação de força da chave fica para a Fase 2 (H14) |
| 2026-11-19 | Bug H10 mantido no refactor, com comentário no `main` | O refactor não pode mudar comportamento; corrigido com teste na Fase 2 |
| 2026-11-19 | `REDIS_PASSWORD` continua opcional | Exigir é item da Fase 3 |
| 2026-11-19 | **Código curto: aleatório base62, 7 chars, `crypto/rand`, retry em violação de unicidade (até 5)** | Não enumerável. 62^7 ≈ 3,5×10¹². Descartado: ID sequencial em base62 (sem retry, mas `/r/1`, `/r/2`… expõe todos os links) |
| 2026-11-19 | **Mesma URL duas vezes: sempre cria novo** | Sem índice extra; cada link é independente; ninguém recebe link de outro por construção. Custo: duplicatas na lista |
| 2026-11-19 | **Redirect 302, rota `GET /r/:slug` fora de `/api`, nginx `location /r/`, URL curta a partir de `PUBLIC_BASE_URL`** | 301 é cacheado para sempre (desativar não teria efeito). 307 só difere em métodos não-GET. `Host` é controlado pelo cliente e perde a porta atrás do nginx |
| 2026-11-19 | **Rate limiter fail open + testes com miniredis** | Redis fora não derruba redirects; o limiter é proteção, não regra de negócio. miniredis: `FastForward` testa a janela sem sleep, roda no CI sem infra. Custo: menos fiel que Redis real |
| 2026-11-19 | **Logout: remover a rota** | JWT é stateless; a rota fingia revogar. Front só descarta o token; TTL de 60 min limita a janela. Descartado: denylist no Redis (+1 GET por requisição autenticada, e com fail open o Redis fora reativaria tokens) |
| 2026-11-19 | **Deletar usuário: `ON DELETE CASCADE` no banco** | Integridade garantida pelo Postgres. Como o AutoMigrate não altera FK existente, `LoadMigrations` ganha um ajuste idempotente da FK (primeira exceção ao AutoMigrate puro) |
| 2026-11-20 | Link de outro usuário responde **404** (não 403) | Não revela que o código existe. Rotas `/user/:id` continuam 403 para id de outro (o id não é segredo) |
| 2026-11-20 | Link desativado no redirect: **410 Gone**; inexistente: 404 | Distingue "existiu e foi desativado" de "nunca existiu", útil para o front |
| 2026-11-20 | Validação → **422**; erro interno → 500 `{"error":"erro interno"}` com log no servidor | Nunca mais mensagem crua do Postgres no corpo |
| 2026-11-20 | Coluna `slug` continua `varchar(64)` | Códigos têm 7 chars; reduzir o tamanho exigiria ALTER que falharia com slugs antigos longos |
| 2026-11-20 | Unicidade reconhecida por SQLSTATE `23505` (`pgconn.PgError`) | Sem `TranslateError` global do GORM, que mudaria a mensagem de outros erros |
| 2026-11-20 | Login: e-mail inexistente compara com hash fictício de mesmo custo | Mesmo status **e** mesmo tempo que senha errada (medido: 1,002 s vs 1,024 s) |
| 2026-11-20 | **`TRUSTED_PROXIES` = faixas privadas, API publicada só em 127.0.0.1, nginx sobrescreve `X-Forwarded-For` com `$remote_addr`** (decisão do usuário) | IP do nginx muda e a sub-rede do compose não é fixa. Anexar o XFF permitiria a um cliente da LAN (IP privado = "confiável") injetar IP. Descartado: sub-rede fixa (colide entre máquinas); não publicar a API (perde depuração direta). Se um CDN/LB ficar na frente do nginx, trocar por `real_ip` |
| 2026-11-20 | `EXPIRE NX` em vez de Lua para a janela fixa | Atômico no MULTI, sem script, e cura chave sem TTL. Exige Redis ≥ 7 (compose usa `redis:latest`) |
| 2026-11-20 | Log de acesso sem query string; GORM com `ParameterizedQueries` | O `gin.Logger` e o logger padrão do GORM vazavam token (query) e e-mail/hash (valores do SQL) |
| 2026-11-20 | `API_KEY` mínima de 32 bytes | HS256 com chave curta é quebrável offline a partir de um token |
| 2026-11-20 | Mensagens de erro novas em português | Projeto e front em PT; mensagens antigas do login (`Invalid email or password`) mantidas para não quebrar o front |

## Contrato da API (mudanças)

Formato de erro: sempre `{"error": "<mensagem>"}`. Sucesso sem entidade: `{"message": "<mensagem>"}`.

| Rota | Antes | Agora |
|---|---|---|
| `GET /api/r/:slug` | existia | **removida** |
| `GET /r/:slug` (novo, fora de `/api`) | — | **302** + `Location`; **404** inexistente; **410** desativado |
| `POST /api/links` | slug = `slug.Make(url)`, reaproveitava link de outro dono; aceitava qualquer esquema; `url` = `scheme://Host/r/…` (quebrada) | código base62 de 7 chars, sempre novo link do usuário; `target_url` só http/https com host, sem credenciais, ≤ 2048 → senão **422**; `url` = `PUBLIC_BASE_URL/r/<código>`; token inválido 401 (era 500 `{}`) |
| `GET /api/links/:slug` | 302 com JSON; 401 para alheio/inexistente | **200** `{link}`; **404** para alheio/inexistente |
| `POST /api/links/:slug/deactivate` | 200 `{success}`; 401 alheio/inexistente | 200 `{message}`; **404** alheio/inexistente |
| `POST /api/user` | duplicado 400; senha > 72 bytes 500 | duplicado **409**; senha > 72 bytes **422** |
| `PUT /api/user/:id` | duplicado 500; sucesso **202** string | duplicado **409**; sucesso **200** `{message}`; 404 se não existe; 403 `{error}` (era `{unauthorized}`) |
| `DELETE /api/user/:id` | 500 se tinha links | **200** `{message}`, links apagados em cascata; 404 se não existe; 403 `{error}` |
| `POST /api/login` | e-mail inexistente 400 | **401**, mesma mensagem da senha errada |
| `POST /api/logout` | 200 sem efeito | **removida (404)** |
| Autenticação | aceitava `?token=` e qualquer esquema no `Authorization`; token sem `exp` nunca expirava | só `Authorization: Bearer`; só HS256; `exp` e `user_id` obrigatórios |
| Rate limit | janela quebrada (1 s renovável), chave `<ip>`, XFF de qualquer um | janela fixa `MINUTES_OF_DURATION` min, chave `ratelimit:<ip>`, XFF só de `TRUSTED_PROXIES`; fail open com log |

Novas variáveis: `PUBLIC_BASE_URL` (obrigatória), `TRUSTED_PROXIES` (opcional). `API_KEY` ≥ 32 bytes; `REQUISITION_LIMIT` e `MINUTES_OF_DURATION` > 0.

## Medições

Nenhuma ainda.

## Achados fora de escopo

- `TOKEN_LIFE_SPAN` é interpretado em **minutos** (`auth/token.go`); o exemplo foi corrigido para 60, mas o nome da variável não diz a unidade.
- O healthcheck do db ainda usa `$${DB_USER}`/`$${DB_NAME}` (vazias no container; passa por acaso). Fase 3.
- `.env` é montado em `/app/.env` além do `env_file` (redundante). Fase 3.
- Tokens não verificam se o usuário ainda existe: após `DELETE /api/user/:id`, o token segue válido até expirar (rotas respondem vazio/404, sem dano, mas é estranho).
- Links `javascript:`/`ftp:` criados antes do H4 continuam no banco (só no ambiente local). Sem migração de limpeza.
- O rate limit é global: conta também `/api/ping` e o redirect `/r/`. Na Fase 4 o teste de carga precisa de `REQUISITION_LIMIT` alto (anotar).
- O token é parseado duas vezes por requisição autenticada (middleware e controller). Poderia guardar o `user_id` no `gin.Context`.
- `LogoutButton.vue` ainda chama `POST /api/logout` (agora 404); o logout funciona porque o token é limpo no `finally`, mas aparece um erro rápido. Fase 6.
- Front ainda usa `/api/r/` em `MyLinksView` e `RedirectTesterView` (rota removida). Fase 6.
- `HASH_COST=14` torna cadastro/login lentos (~1 s por hash); relevante para testes de carga em rotas autenticadas.
- Gin em modo debug (sem `GIN_MODE=release`).
- `LinkDetailsView` lê `res.data.target_url`, mas a API responde `{link:{...}}`; e como a API responde 302, o axios rejeita a resposta → a tela de detalhe sempre mostra erro.
- `AccountView` exige que o usuário digite o próprio ID (não há `GET /me`).
- `RegisterView` ignora o token devolvido no cadastro; `LoginView` grava o token direto no `localStorage` sem passar por `useAuth().login` (menu não atualiza).
- `VITE_API_BASE` passado como `environment` em runtime no compose não tem efeito: variáveis Vite são de build.
- Sem `.dockerignore` em `api/` e `web/`.
- DSN fixa `TimeZone=America/Sao_Paulo` (`config/database.go:21`).

## Ideias futuras

- (vazio)

## Próximo passo

Fase 3 (branch `chore/infra-nginx`): Dockerfile multi-stage da API (não-root), front servido por nginx com `try_files` e gzip (resolve o 502 do H16), proxy com `/api` e `/r/`, healthchecks (db com `POSTGRES_*`, api com endpoint de health), exigir `REDIS_PASSWORD`, documentar segurança/TLS.
Decidir com o usuário antes: o branch `chore/infra-nginx` sai de `fix/api-stabilize` (recomendado: a `main` ainda não tem as fases 0–2) ou `fix/api-stabilize` é mergeado na `main` primeiro.
