FROM golang:1.14.4-alpine3.12

RUN mkdir -p app

ADD . /app

WORKDIR /app

RUN go build main.go

CMD /app/main
