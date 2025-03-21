FROM golang:1.22-alpine

WORKDIR /app

RUN apk add --no-cache postgresql-client

COPY . .

RUN go mod download
RUN go build -o main src/cmd/api/main.go

EXPOSE 8080

COPY wait-for-postgres.sh /wait-for-postgres.sh
RUN chmod +x /wait-for-postgres.sh

CMD ["/wait-for-postgres.sh", "./main"]