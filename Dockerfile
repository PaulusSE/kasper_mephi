FROM golang:1.21.5-alpine3.18

# Установка Python и pip
RUN apk add --update --no-cache python3 py3-pip && \
    python3 -m ensurepip && \
    rm -r /usr/lib/python*/ensurepip && \
    pip3 install --no-cache --upgrade pip setuptools

# Установка библиотеки python-pptx
RUN pip3 install python-pptx

WORKDIR /usr/src/app

# pre-copy/cache go.mod for pre-downloading dependencies and only redownloading them in subsequent builds if they change
COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .
RUN go build -o ./bin/server ./cmd/kasper/main.go

EXPOSE 8080

CMD ["./bin/server"]