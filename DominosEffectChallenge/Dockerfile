
  

# Build stage (Go required)
FROM golang:1.23 AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o app ./server/cmd

# Runtime stage (Go NOT required)
FROM gcr.io/distroless/base-debian12
WORKDIR /app
COPY --from=builder /app/app .
CMD ["./app"]
