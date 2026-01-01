# -------- Build stage --------
FROM golang:1.25.5-alpine AS builder


WORKDIR /app

# Copy dependency files first (cache optimization)
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build a static Linux binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -o myapp


# -------- Runtime stage --------
FROM alpine:latest

WORKDIR /app

# Copy only the binary
COPY --from=builder /app/myapp .

# Expose app port
EXPOSE 8080

# Run the app
CMD ["./myapp"]    


