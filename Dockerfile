# -------- BUILD STAGE --------
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Аргумент для выбора сервиса при сборке
ARG SERVICE

COPY ${SERVICE}/go.mod ${SERVICE}/go.sum ./
RUN go mod download

COPY ${SERVICE}/ ./

RUN go build -o /bin/${SERVICE} ./cmd/main.go

# -------- RUNTIME STAGE --------
FROM alpine:3.20

WORKDIR /app

ARG SERVICE
COPY --from=builder /bin/${SERVICE} /usr/local/bin/${SERVICE}
RUN apk add --no-cache tzdata

ENV GIN_MODE=release

CMD ["/bin/sh", "-c", "${SERVICE}"]
