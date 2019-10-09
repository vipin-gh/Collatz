FROM golang:latest

LABEL maintainer="AshishM <mahamuni.ashish@gmail.com>"
LABEL description="producer image, which pushes the numbers into RabbitMQ"

WORKDIR /app

COPY go.mod go.sum main.go ./

RUN go mod download

RUN go build -o main .

CMD ["./main"]
