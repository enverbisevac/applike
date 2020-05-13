FROM golang:1.12

ENV GO111MODULE=on

WORKDIR /go/src/app

COPY cmd cmd
COPY pkg pkg
COPY migrations migrations
COPY go.mod go.mod


RUN go mod download

RUN go build -o applike ./cmd/

EXPOSE 3000

CMD ["./applike"]