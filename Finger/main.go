package main

import (
	"encoding/json"
	"roblox-universe-finger/body"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

type RequestData struct {
	MaxID *int64 `json:"maxID"`
}

type ResponseData struct {
	Success  int8    `json:"success"`
	ValidIDs []int64 `json:"validIDs"`
}

func main() {
	lambda.Start(handler)
}

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var requestData RequestData
	if err := json.Unmarshal([]byte(request.Body), &requestData); err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 400,
			Body:       "Invalid request data",
		}, nil
	}

	finger := body.Finger(*requestData.MaxID)
	validIDs, suc := finger.Feel()

	responseData := ResponseData{
		Success:  suc,
		ValidIDs: validIDs,
	}

	// test
	// test 2
	//test 3
	//test 4789
	//test 5
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
