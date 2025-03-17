# FROM golang:1.23 AS builder

# WORKDIR /app

# COPY go.mod go.sum ./
# RUN go mod download

# COPY . .

# RUN go build -o codestacker-api ./cmd/main.go

# FROM alpine:latest
# WORKDIR /root/

# COPY --from=builder /app/codestacker-api .

# EXPOSE 8080

# CMD ["./codestacker-api"]
