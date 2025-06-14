FROM golang:1.24.4-alpine3.22 AS builder

WORKDIR /src
COPY go.mod go.sum .
RUN go mod download

COPY . .
RUN go mod tidy
RUN go test ./... && go build -o users cmd/users/main.go

FROM alpine:3.22

WORKDIR /app
COPY --from=builder /src/users .

ENTRYPOINT ["/app/users"]
