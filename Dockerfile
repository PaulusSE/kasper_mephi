# СТЕЙДЖ 1: Build Go binary
FROM golang:1.21.5-alpine3.18 AS go-builder

WORKDIR /usr/src/app
COPY go.mod go.sum ./
RUN go mod download && go mod verify
COPY . .
ENV CGO_ENABLED=0
RUN go build -o /bin/server ./cmd/kasper/main.go

# СТЕЙДЖ 2: Python-окружение с оптимизацией пространства
FROM python:3.11-slim AS python-builder

ENV PYTHONDONTWRITEBYTECODE=1
ENV PYTHONUNBUFFERED=1

# Установка зависимостей с очисткой кэша
RUN apt-get update && \
    apt-get install -y --no-install-recommends \
    build-essential \
    python3-dev \
    libgomp1 \
    libopenblas-dev && \
    rm -rf /var/lib/apt/lists/* && \
    apt-get clean all

# Копируем и устанавливаем зависимости
COPY requirements.txt /tmp/requirements.txt

# Установка всех зависимостей и очистка в одном RUN
RUN pip install --no-cache-dir torch==2.1.2+cpu -f https://download.pytorch.org/whl/torch_stable.html && \
    pip install --no-cache-dir -r /tmp/requirements.txt && \
    python -m spacy download ru_core_news_sm && \
    apt-get purge -y --auto-remove build-essential python3-dev && \
    rm -rf /root/.cache /tmp/*

# СТЕЙДЖ 3: Финальный образ (минимальный)
FROM python:3.11-slim

# Установка runtime зависимостей
RUN apt-get update && \
    apt-get install -y --no-install-recommends \
    libgomp1 \
    libopenblas0 && \
    rm -rf /var/lib/apt/lists/* && \
    apt-get clean all

# Копируем артефакты из предыдущих стейджей
COPY --from=python-builder /usr/local/lib/python3.11/site-packages /usr/local/lib/python3.11/site-packages
COPY --from=python-builder /usr/local/bin /usr/local/bin
COPY --from=go-builder /bin/server /app/bin/server

# Копируем необходимые файлы
WORKDIR /app
COPY --from=go-builder /usr/src/app/internal/app/reco_model_2.py ./internal/app/
COPY --from=go-builder /usr/src/app/internal/app/parsed_articles.pkl ./internal/app/
COPY --from=go-builder /usr/src/app/internal/pkg/service/presentation/generate_presentation.py ./internal/pkg/service/presentation/

# Настройки
ENV PYTHONPATH=/app
ENV LD_PRELOAD=/usr/lib/x86_64-linux-gnu/libgomp.so.1

EXPOSE 8080
CMD ["/app/bin/server"]
