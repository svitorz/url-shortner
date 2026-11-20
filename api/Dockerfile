# syntax=docker/dockerfile:1

# --- build: compila um binário estático ---
FROM golang:1.24.4-alpine AS build
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
# -trimpath tira caminhos locais do binário; -s -w remove símbolos de debug.
# timetzdata embute a base de fusos (a DSN usa America/Sao_Paulo e o alpine
# não traz tzdata): o binário não depende da imagem de runtime para isso.
RUN CGO_ENABLED=0 GOOS=linux go build -tags timetzdata -trimpath -ldflags="-s -w" -o /out/url-shortner .

# --- runtime: só o binário, rodando sem root ---
FROM alpine:3.22
# Usuário sem shell de login, sem home; UID fixo facilita auditar.
RUN adduser -D -H -u 10001 -s /sbin/nologin app
COPY --from=build /out/url-shortner /usr/local/bin/url-shortner
ENV GIN_MODE=release
USER app
# Documentação apenas: a porta real vem de APP_PORT.
EXPOSE 8000
ENTRYPOINT ["/usr/local/bin/url-shortner"]
