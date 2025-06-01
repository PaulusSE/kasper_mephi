FROM golang:1.21.5-alpine3.18 AS builder

# Установка зависимостей для сборки
RUN apk add --update --no-cache \
    gcc \
    musl-dev \
    python3 \
    py3-pip \
    && python3 -m ensurepip \
    && rm -r /usr/lib/python*/ensurepip \
    && pip3 install --no-cache --upgrade pip setuptools

# Установка python-pptx с очисткой кэша
RUN pip3 install --no-cache python-pptx

WORKDIR /usr/src/app

# Копируем только файлы зависимостей Go
COPY go.mod go.sum ./
RUN go mod download && go mod verify

# Копируем остальные файлы
COPY . .

# Сборка приложения
RUN go build -o ./bin/server ./cmd/kasper/main.go

# Финальный образ
FROM alpine:3.18

# Установка только необходимых runtime-зависимостей
RUN apk add --update --no-cache python3 libstdc++

# Копируем python-pptx из builder-этапа
COPY --from=builder /usr/lib/python3.11/site-packages/ /usr/lib/python3.11/site-packages/
COPY --from=builder /usr/bin/pptx /usr/bin/pptx*

# Копируем собранное приложение
COPY --from=builder /usr/src/app/bin/server /app/server

# Очистка ненужных файлов
RUN rm -rf /var/cache/apk/* /tmp/*

WORKDIR /app
EXPOSE 8080

CMD ["./bin/server"]
