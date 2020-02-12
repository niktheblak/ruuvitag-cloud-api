package main

import (
	"net/http"

	"github.com/aws/aws-lambda-go/events"
)

func okResponse(body ...string) events.APIGatewayProxyResponse {
	return resp(http.StatusOK, body...)
}

func badRequest(body ...string) events.APIGatewayProxyResponse {
	return resp(http.StatusBadRequest, body...)
}

func internalServerError(body ...string) events.APIGatewayProxyResponse {
	return resp(http.StatusInternalServerError, body...)
}

func forbidden(body ...string) events.APIGatewayProxyResponse {
	return resp(http.StatusForbidden, body...)
}

func resp(status int, body ...string) events.APIGatewayProxyResponse {
	var respBody string
	if len(body) > 0 {
		respBody = body[0]
	}
	return events.APIGatewayProxyResponse{
		StatusCode: status,
		Body:       respBody,
	}
}
