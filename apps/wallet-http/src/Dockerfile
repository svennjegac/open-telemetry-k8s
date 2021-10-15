FROM golang:1.17 AS build

WORKDIR /app

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .

RUN go build -o /app-exe ./cmd/api

EXPOSE 8112

ENTRYPOINT ["/app-exe"]
