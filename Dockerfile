FROM golang:1.22-alpine as builder
WORKDIR /app
COPY . .

RUN go build  ./cmd/main.go

FROM alpine:3.15


WORKDIR /app
COPY --from=builder /app/.env .
COPY --from=builder /app/config/local.yaml ./config/local.yaml
COPY --from=builder /app/main .


CMD ["/app/main"]