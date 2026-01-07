FROM golang:1.24-alpine

WORKDIR /app

# зависимости
COPY go.mod go.sum ./
RUN go mod download

# код
COPY . .

# билд бинарника
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o api ./cmd/api

EXPOSE 8080

CMD ["./api"]