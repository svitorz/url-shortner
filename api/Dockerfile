FROM golang:1.24.4

ARG PORT
ENV PORT=$APP_PORT

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY * ./ 

RUN CGO_ENABLED=0 GOOS=linux go build -o /url-shortner

EXPOSE "$PORT" 

CMD ["/url-shortner"]