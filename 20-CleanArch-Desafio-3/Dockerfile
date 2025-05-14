FROM golang:1.21

WORKDIR /app

COPY . .

RUN GOPROXY=direct go mod download

RUN CGO_ENABLED=1 GOOS=linux go build -o main ./cmd/ordersystem

EXPOSE 8000 50051 8080

CMD ["./main"]