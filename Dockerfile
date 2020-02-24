FROM golang:1 as build-env

WORKDIR /go/src/app
ADD . /go/src/app

RUN go build -o /go/bin/app cmd/api/*.go

FROM gcr.io/distroless/base
COPY --from=build-env /go/bin/app /
ENTRYPOINT ["/app"]
