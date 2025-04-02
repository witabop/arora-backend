package main

import (
	"arora-search-brain/body"
	"arora-search-brain/layer"
	"encoding/json"
	"log"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func main() {
	lambda.Start(handler)
}

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var (
		requestData layer.RequestData
		headers     = map[string]string{
			"Access-Control-Allow-Origin":  "*",
			"Access-Control-Allow-Methods": "POST",
		}
	)

	if err := json.Unmarshal([]byte(request.Body), &requestData); err != nil {
		log.Println(err.Error())
		return events.APIGatewayProxyResponse{
			StatusCode: 400,
			Headers:    headers,
			Body:       "Invalid request data",
		}, nil
	}

	brain := body.Brain(requestData.NumGames, requestData.SearchCriteria)
	validUniverses := brain.Think()

	responseData := layer.ResponseData{
		Data: validUniverses,
	}

	bytesResponse, _ := json.Marshal(responseData)
	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Headers:    headers,
		Body:       string(bytesResponse),
	}, nil
}
