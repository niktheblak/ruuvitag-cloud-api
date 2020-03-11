// +heroku goVersion go1.14
// +heroku install ./cmd/heroku/server

module github.com/niktheblak/ruuvitag-cloud-api

go 1.14

require (
	github.com/aws/aws-lambda-go v1.14.0
	github.com/aws/aws-sdk-go v1.29.11
	github.com/julienschmidt/httprouter v1.3.0
	github.com/lib/pq v1.3.0
	github.com/niktheblak/ruuvitag-gollector v0.0.0-20200224125718-a9fc8c932fea
)
