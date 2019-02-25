FROM golang:alpine

RUN apk update && apk add --no-cache git gcc bash make

ADD . /go/src/github.com/plusspeed/${SERVICE}
WORKDIR /go/src/github.com/plusspeed/${SERVICE}
RUN make clean install
CMD ["/app/main"]
