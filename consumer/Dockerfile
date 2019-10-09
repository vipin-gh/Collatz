FROM golang:latest

LABEL maintainer="AshishM <mahamuni.ashish@gmail.com>"
LABEL description="consumer image, which pushes the numbers into RabbitMQ"

WORKDIR /app

COPY main.go ./

RUN go get -d ./...

RUN go build -o main .

CMD ["./main"]
