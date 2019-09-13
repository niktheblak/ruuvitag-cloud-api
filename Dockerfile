FROM golang:1.12 as build-env

WORKDIR /go/src/app
ADD . /go/src/app

# RUN go get cloud.google.com/go/firestore
RUN go get -d ./...
RUN go build -o /go/bin/app appengine/*.go

FROM gcr.io/distroless/base
COPY --from=build-env /go/bin/app /
CMD ["/app"]
