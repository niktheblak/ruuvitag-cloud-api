// +heroku goVersion go1.17
// +heroku install ./cmd/heroku/server

module github.com/niktheblak/ruuvitag-cloud-api

go 1.17

require (
	github.com/aws/aws-lambda-go v1.27.0
	github.com/aws/aws-sdk-go v1.41.4
	github.com/julienschmidt/httprouter v1.3.0
	github.com/lib/pq v1.10.3
	github.com/niktheblak/ruuvitag-gollector v0.0.0-20210821201042-c7b52442eaf9
)

require github.com/jmespath/go-jmespath v0.4.0 // indirect
