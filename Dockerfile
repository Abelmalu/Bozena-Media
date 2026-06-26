# -------- Build stage --------
FROM golang:1.25-alpine AS builder

WORKDIR /app

ARG SERVICE_NAME
ARG MAIN_PATH=cmd/main.go

COPY go.mod go.sum ./

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -o myapp ./${SERVICE_NAME}/${MAIN_PATH}

# -------- Runtime stage --------
FROM alpine:latest

WORKDIR /app


COPY --from=builder /app/myapp .


CMD ["./myapp"]
