FROM golang:1.22 as build

WORKDIR /go/src/app

COPY go.mod go.sum ./
RUN go mod download && go mod verify
ADD . .
RUN go build -v -tags aws,postgres -o /go/bin/app

FROM ubuntu:latest
COPY --from=build /go/bin/app /
CMD ["/app"]
