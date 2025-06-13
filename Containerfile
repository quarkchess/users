FROM golang:1.24.4-alpine3.22 AS builder

WORKDIR /src
COPY go.mod go.sum .
RUN go mod tidy

COPY . .
RUN go test ./... && go build -o auth cmd/auth/main.go

FROM alpine:3.22

WORKDIR /app
COPY --from=builder /src/auth .

ENV DB_CONNSTRING=""
ENV LISTEN_ADDRESS=":8000"

ENTRYPOINT ["/app/auth"]
