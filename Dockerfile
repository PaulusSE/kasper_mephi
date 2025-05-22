# СТЕЙДЖ 1: Build Go binary
FROM golang:1.21.5-alpine3.18 AS go-builder

WORKDIR /usr/src/app
COPY go.mod go.sum ./
RUN go mod download && go mod verify
COPY . .
ENV CGO_ENABLED=0
RUN go build -o /bin/server ./cmd/kasper/main.go

# СТЕЙДЖ 2: Python-окружение
FROM python:3.11-slim

ENV PYTHONDONTWRITEBYTECODE=1
ENV PYTHONUNBUFFERED=1

RUN apt-get update && apt-get install -y --no-install-recommends \
    build-essential \
    python3-dev \
    git \
    libgomp1 \
    libopenblas-dev \
    && rm -rf /var/lib/apt/lists/*

COPY requirements.txt /tmp/requirements.txt
RUN pip install --no-cache-dir torch==2.1.2+cpu -f https://download.pytorch.org/whl/torch_stable.html
RUN pip install --no-cache-dir -r /tmp/requirements.txt
RUN python -m spacy download ru_core_news_sm@3.7.0 || { echo "Error installing dependencies"; exit 1; }


WORKDIR /usr/src/app
COPY --from=go-builder /bin/server ./bin/server
COPY internal/app/reco_model_2.py .
COPY internal/app/parsed_articles.pkl .
COPY internal/pkg/service/presentation/generate_presentation.py .

EXPOSE 8080
CMD ["./bin/server"]
