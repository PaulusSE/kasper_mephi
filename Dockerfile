# СТЕЙДЖ 1: Build Go binary на любой удобной платформе
FROM golang:1.21.5-alpine3.18 AS go-builder

WORKDIR /usr/src/app

# Кэшируем модули заранее
COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .
RUN go build -o /bin/server ./cmd/kasper/main.go

FROM python:3.11-slim

# ———————
# 1. Установим все build-зависимости для pip и ML
RUN apt-get update && apt-get install -y --no-install-recommends \
    build-essential \
    python3-dev \
    git \
    && rm -rf /var/lib/apt/lists/*

# 2. Установим Python-зависимости (PyTorch подтянет wheel)
COPY requirements.txt /tmp/requirements.txt
RUN pip install --no-cache-dir --upgrade pip \
    && pip install --no-cache-dir -r /tmp/requirements.txt \
    && python -m spacy download ru_core_news_sm

WORKDIR /usr/src/app

# 3. Копируем Go-бинарь из предыдущего стейджа
COPY --from=go-builder /bin/server /usr/src/app/bin/server

# 4. Копируем python-скрипты и данные (убери этот COPY, если они уже есть после COPY . .)
COPY internal/app/reco_model_2.py /usr/src/app/reco_model_2.py
COPY internal/app/parsed_articles.pkl /usr/src/app/parsed_articles.pkl

# 5. Копируем остальной проект (если нужен)
COPY . .

EXPOSE 8080

CMD ["./bin/server"]