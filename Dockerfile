FROM golang:alpine

WORKDIR /aftersales

COPY . .

RUN go mod tidy

RUN apk add --no-cache git \
    && go install github.com/swaggo/swag/cmd/swag@latest

# Add the Go binary path (for swag) to the PATH environment variable
ENV PATH=$PATH:/root/go/bin

# Initialize the swag documentation
RUN swag init

CMD ["go","run","main.go"]