FROM golang:1.23-alpine

WORKDIR /app

COPY . . 

RUN go get -d -v ./...

RUN go build -o crime-manager-api ./cmd/main.go

EXPOSE 8080

CMD ["./crime-manager-api"]