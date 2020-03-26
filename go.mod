// +heroku goVersion go1.14
// +heroku install ./cmd/heroku/server

module github.com/niktheblak/ruuvitag-cloud-api

go 1.14

require (
	github.com/aws/aws-lambda-go v1.15.0
	github.com/aws/aws-sdk-go v1.29.32
	github.com/jmespath/go-jmespath v0.3.0 // indirect
	github.com/julienschmidt/httprouter v1.3.0
	github.com/lib/pq v1.3.0
	github.com/niktheblak/ruuvitag-gollector v0.0.0-20200311123532-8b4933f1b309
)
