# ---------- Stage 1: Build ----------
FROM golang:1.24-alpine AS builder

ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64

WORKDIR /server

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o server ./cmd

# ---------- Stage 2: Run ----------
FROM alpine:latest

WORKDIR /server

COPY --from=builder /server/server .

EXPOSE 8080

ENTRYPOINT ["./server"]
