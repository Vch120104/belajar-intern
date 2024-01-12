FROM golang:alpine

WORKDIR /aftersales

COPY . .

RUN go mod tidy

CMD ["go","run","main.go"]