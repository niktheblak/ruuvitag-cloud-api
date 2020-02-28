package main

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

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

func ParseLimit(limitStr string) int {
	var limit int64
	if limitStr != "" {
		limit, _ = strconv.ParseInt(limitStr, 10, 64)
	}
	if limit <= 0 {
		limit = 20
	}
	return int(limit)
}

func ParseTimeRange(fromStr, toStr string) (from time.Time, to time.Time, err error) {
	if fromStr != "" {
		from, err = time.Parse("2006-01-02", fromStr)
	}
	if err != nil {
		return
	}
	if toStr != "" {
		to, err = time.Parse("2006-01-02", toStr)
	}
	if err != nil {
		return
	}
	if !from.IsZero() && !to.IsZero() && from == to {
		to = to.AddDate(0, 0, 1)
	}
	if to.IsZero() || to.After(time.Now()) {
		to = time.Now().UTC()
	}
	if from.After(to) {
		err = fmt.Errorf("from timestamp cannot be after to timestamp")
	}
	return
}
