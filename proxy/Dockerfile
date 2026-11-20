# syntax=docker/dockerfile:1
# Contexto principal: proxy/. O código do front entra pelo contexto extra "web"
# (additional_contexts no docker-compose.yml), sem expor a raiz do repo (.env).

# --- build do front (Vite) ---
FROM node:24-alpine AS front
WORKDIR /web
COPY --from=web package.json package-lock.json ./
RUN npm ci
COPY --from=web . .
# Variáveis VITE_* são lidas no build, não em runtime.
ARG VITE_API_BASE=/api
ENV VITE_API_BASE=$VITE_API_BASE
RUN npm run build

# --- nginx: serve o build e faz proxy para a API ---
FROM nginx:1.30-alpine
COPY templates/ /etc/nginx/templates/
COPY --from=front /web/dist /usr/share/nginx/html
