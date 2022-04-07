// +heroku goVersion go1.17
// +heroku install ./cmd/heroku/server

module github.com/niktheblak/ruuvitag-cloud-api

go 1.7

require (
	github.com/aws/aws-lambda-go v1.29.0
	github.com/aws/aws-sdk-go v1.43.34
	github.com/julienschmidt/httprouter v1.3.0
	github.com/lib/pq v1.10.4
	github.com/niktheblak/ruuvitag-gollector v0.0.0-20210821201042-c7b52442eaf9
)
