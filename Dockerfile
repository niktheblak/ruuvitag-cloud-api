FROM golang:1.14 as build-env

WORKDIR /go/src/app
ADD . /go/src/app

RUN go get -d -v ./...
RUN go build -o /go/bin/app cmd/api/*.go

FROM gcr.io/distroless/base-debian10
COPY --from=build-env /go/bin/app /
ENTRYPOINT ["/app"]
