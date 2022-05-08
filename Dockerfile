FROM golang:1.18.1-buster

WORKDIR /

RUN mkdir /go/src/webex

COPY . /go/src/webex

WORKDIR /go/src/webex

RUN go mod tidy

ARG HOST

ENV HOST $HOST

RUN CGO_ENABLED=1 GOOS=linux GOARCH=amd64 \
    go build -o server .

EXPOSE 3000/tcp

Entrypoint ["./server"]
