FROM golang:1.20.5-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o app

FROM alpine:3.15
RUN apk --no-cache add ca-certificates
WORKDIR /app
COPY --from=builder /app/app .
ENTRYPOINT ["/app/app"]