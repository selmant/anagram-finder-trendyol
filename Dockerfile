FROM golang:1.20.5-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build

FROM alpine:3.15
RUN apk --no-cache add ca-certificates
WORKDIR /app
COPY --from=builder /app/anagram-finder-trendyol .
ENTRYPOINT ["/app/anagram-finder-trendyol"]