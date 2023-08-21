FROM golang:1.21 as build

WORKDIR /go/src/app
ADD . /go/src/app

RUN go get -d -v ./...

RUN go build -tags aws,postgres -o /go/bin/app

FROM ubuntu:latest
COPY --from=build /go/bin/app /
CMD ["/app"]
