package main

import (
	"arora-search-finger/body"
	"arora-search-finger/layer"
	"encoding/json"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func main() {
	lambda.Start(handler)
}

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var requestData layer.RequestData
	if err := json.Unmarshal([]byte(request.Body), &requestData); err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 400,
			Body:       "Invalid request data",
		}, nil
	}

	finger := body.Finger(requestData.MaxID, requestData.SearchCriteria, request.RequestContext.Stage)
	validUniverses := finger.Feel()

	responseData := layer.ResponseData{
		ValidUniverses: validUniverses,
	}

	bytesResponse, err := json.Marshal(responseData)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       "Failed to marshal response data",
		}, nil
	}

	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       string(bytesResponse),
	}, nil
}
