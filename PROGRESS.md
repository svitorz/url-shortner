# PROGRESS — url-shortner

> Memória do trabalho. Ler ANTES de qualquer coisa em uma sessão nova ou após compactação de contexto.

## Fase atual

**Fase 0 — orientação (somente leitura): concluída.** Aguardando "continuar" para a Fase 1.
Branch: `fix/api-stabilize` (criado a partir de `main` @ `f15ea2a`).

## Checklist

| Fase | Item | Status |
|---|---|---|
| 0 | Ler api/, web/, proxy/, compose, docs | ✅ |
| 0 | Verificar H1–H22 (estático + execução contra Postgres/Redis descartáveis) | ✅ |
| 0 | Branch `fix/api-stabilize` + `PROGRESS.md` | ✅ |
| 1 | Corrigir H1 (migrations) — decisão AutoMigrate vs migrations SQL | ⏳ aguardando decisão |
| 1 | Volume persistente do Postgres | ⏳ |
| 1 | `.env.example` coerente + `docker compose up -d --build` funcionando | ⏳ |
| 1 | Script de fumaça cobrindo todas as rotas | ⏳ |
| 1 | Tirar config/DB/Redis dos `init()` e montar no `main` | ⏳ |
| 2–7 | ver plano | ⏳ |

## Como a Fase 0 foi verificada

- `go vet ./...` e `go build` com Go 1.24.4 (máquina local, Debian, kernel 6.12).
- Binário da API rodado no host contra containers descartáveis `postgres:18` e `redis` (portas 127.0.0.1:55432/56379, removidos ao final). O stack do projeto **não** foi subido (subir a API apagaria qualquer banco existente — H1).
- Script de sondagem com curl (scratchpad, não versionado) exercitando todas as rotas.
- `docker compose config` rodado numa cópia do compose com `.env` = `.env.example` (no scratchpad), só para ver a interpolação.
- Nenhum `.env` existe no repositório; histórico do git só contém `.env.example` (sem segredos reais).

## Resultado atual das rotas (código de `main`, limiter alto)

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

## Problemas verificados (resumo; detalhes no relatório da Fase 0)

Críticos: H1 (DROP em todo boot), segredo JWT forjável (`API_KEY` vazio aceito; com o `.env.example` via compose o segredo vira o texto do comentário), H2 (link de outro dono), H8 (delete 500), proxy fora da rede da API.
Ver relatório completo na conversa da Fase 0; cada item será tratado na fase correspondente.

## Log de decisões

| Data | Decisão | Porquê |
|---|---|---|
| 2026-11-19 | Não subir o stack do projeto na Fase 0 | H1 apagaria dados de um banco existente; verificação feita em containers descartáveis |

## Contrato da API (mudanças)

Nenhuma ainda.

## Medições

Nenhuma ainda.

## Achados fora de escopo

- `TOKEN_LIFE_SPAN` é interpretado em **minutos** (`auth/token.go:28`); o `.env.example` usa 3600 → token de 60 horas.
- Tokens não verificam se o usuário ainda existe.
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

Fase 1: decidir com o usuário AutoMigrate vs migrations SQL versionadas; corrigir H1 primeiro.
