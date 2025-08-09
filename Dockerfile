# Use official Golang image as base
FROM golang:1.24-alpine

WORKDIR /PASSMAN

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o main .

EXPOSE 1019

CMD go test ./... -v && ./main
