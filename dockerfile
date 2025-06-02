FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 go build -o app .

FROM alpine:3.22

WORKDIR /app

COPY --from=builder /app/app /app/
COPY cotacoes/ /app/cotacoes/

EXPOSE 8080

CMD ["./app", "-processar"]
