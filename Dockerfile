FROM golang:1.14 as base

WORKDIR /go/src/sammy

COPY . .

RUN go get -d -v ./...
RUN go build main.go

CMD ["./main"]
