package main

import (
	"net/http"

	"github.com/aws/aws-lambda-go/events"
)

func OKResponse(body ...string) events.APIGatewayProxyResponse {
	return Response(http.StatusOK, body...)
}

func BadRequest(body ...string) events.APIGatewayProxyResponse {
	return Response(http.StatusBadRequest, body...)
}

func InternalServerError(body ...string) events.APIGatewayProxyResponse {
	return Response(http.StatusInternalServerError, body...)
}

func Response(status int, body ...string) events.APIGatewayProxyResponse {
	var respBody string
	if len(body) > 0 {
		respBody = body[0]
	}
	return events.APIGatewayProxyResponse{
		StatusCode: status,
		Body:       respBody,
	}
}
