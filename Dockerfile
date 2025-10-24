# ---------- BUILD STAGE ----------
# Можно собрать под любую платформу: docker buildx ...
ARG GO_VERSION=1.24.8
FROM --platform=$BUILDPLATFORM golang:${GO_VERSION}-alpine AS builder

ENV GOTOOLCHAIN=auto \
    CGO_ENABLED=0

# Включим кэш модулей и базовые утилиты
RUN apk add --no-cache git build-base

# Имя сервисной папки передаём аргументом
# Пример: --build-arg SERVICE=auth_service
ARG SERVICE
ARG TARGETOS
ARG TARGETARCH

# Рабочая директория = корень репо
WORKDIR /src

# Сначала тянем только мод-файлы, чтобы кэшировались зависимости
# Структура ожидается: <repo>/<SERVICE>/go.mod, go.sum, cmd/main.go, internal/...
COPY ${SERVICE}/go.mod ${SERVICE}/go.sum ${SERVICE}/
WORKDIR /src/${SERVICE}
RUN go mod download

# Теперь копируем исходники сервиса и собираем статический бинарь
COPY ${SERVICE}/ ./

ENV CGO_ENABLED=0
RUN GOOS=${TARGETOS:-linux} GOARCH=${TARGETARCH:-amd64} \
    go build -ldflags "-s -w" -o /out/app ./cmd

# ---------- RUNTIME STAGE ----------
FROM alpine:3.20

# Полезные пакеты в рантайме (логирование времени, TLS, healthcheck через wget/curl)
RUN apk add --no-cache ca-certificates tzdata curl

# Кладём бинарь в фиксированный путь
COPY --from=builder /out/app /usr/local/bin/app

# Непривилегированный пользователь
USER 10001:10001
WORKDIR /app

# Жёстко заданный entrypoint (без переменных!)
ENTRYPOINT ["/usr/local/bin/app"]
