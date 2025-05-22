# СТЕЙДЖ 1: Build Go binary
FROM golang:1.21.5-alpine3.18 AS go-builder

WORKDIR /usr/src/app
COPY go.mod go.sum ./
RUN go mod download && go mod verify
COPY . .
RUN go build -o /bin/server ./cmd/kasper/main.go

# СТЕЙДЖ 2: Python-окружение
FROM python:3.11-slim

# 1. Системные зависимости
# hadolint ignore=DL3008
RUN apt-get update && apt-get install -y --no-install-recommends \
    build-essential \
    python3-dev \
    git \
    libgomp1 \
    libopenblas-dev \
    && rm -rf /var/lib/apt/lists/*

# 2. Python-зависимости
# hadolint ignore=DL3008
COPY requirements.txt /tmp/requirements.txt
RUN pip install --no-cache-dir "pip==23.0.1" && \
    pip install --no-cache-dir -r /tmp/requirements.txt && \
    python -m spacy download ru_core_news_sm@3.7.0 || \
    { echo "Error installing dependencies"; exit 1; }  # Явный выход при ошибке

# 3. Копируем Go-бинарь и данные
WORKDIR /usr/src/app
COPY --from=go-builder /bin/server ./bin/server
COPY internal/app/reco_model_2.py .
COPY internal/app/parsed_articles.pkl .
COPY internal/pkg/service/presentation/generate_presentation.py .

EXPOSE 8080
CMD ["./bin/server"]
