FROM golang:1.21.5-alpine3.18 AS builder

WORKDIR /usr/src/app

# pre-copy/cache go.mod for pre-downloading dependencies and only redownloading them in subsequent builds if they change
COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .
RUN go build -ldflags="-s -w" -o ./bin/server ./cmd/kasper/main.go

EXPOSE 8080
FROM alpine

WORKDIR /usr/src/app

COPY --from=builder /usr/src/app/bin/server /usr/src/app/bin/server
CMD ["./bin/server"]