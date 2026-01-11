FROM node:20-alpine

ARG PORT
ENV PORT=$APP_PORT

RUN npm install -g http-server

WORKDIR /web

COPY package*.json ./

RUN npm install

COPY . .

RUN npm run build

EXPOSE "$PORT"

CMD ["http-server", "dist"]