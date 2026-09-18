# syntax=docker/dockerfile:1

FROM golang:1.26-alpine AS builder
WORKDIR /src

RUN apk add --no-cache git ca-certificates

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /out/fizzbuzz-api ./cmd/api


FROM gcr.io/distroless/static-debian12:nonroot AS runtime

WORKDIR /app
COPY --from=builder /out/fizzbuzz-api /app/fizzbuzz-api

USER nonroot:nonroot
EXPOSE 8080

ENTRYPOINT ["/app/fizzbuzz-api"]
