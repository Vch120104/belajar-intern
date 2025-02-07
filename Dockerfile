FROM golang:alpine

WORKDIR /aftersales

COPY . .

RUN go mod tidy

RUN swag init

CMD ["go","run","main.go"]
