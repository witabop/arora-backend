package main

import (
	"arora-search-brain/body"
	"fmt"

	"github.com/aws/aws-lambda-go/lambda"
)

func main() {
	lambda.Start(handler)
}

func handler() {
	fmt.Println("temp")
	// var requestData RequestData
	// if err := json.Unmarshal([]byte(request.Body), &requestData); err != nil {
	// 	return events.APIGatewayProxyResponse{
	// 		StatusCode: 400,
	// 		Body:       "Invalid request data",
	// 	}, nil
	// }

	body.Test()
}
