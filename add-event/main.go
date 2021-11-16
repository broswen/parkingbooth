package main

import (
	"context"
	"encoding/json"
	"log"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/broswen/parkingbooth/account"
	"github.com/broswen/parkingbooth/location"
	"github.com/broswen/parkingbooth/models"
)

var accountService *account.Service

func Handler(ctx context.Context, event events.APIGatewayProxyRequest) (models.APIGResponse, error) {
	var addEventRequest models.AddEventRequest
	err := json.Unmarshal([]byte(event.Body), &addEventRequest)
	if err != nil {
		return models.GenerateError(err.Error(), 400)
	}
	e := account.AccountEvent{
		AccountId: addEventRequest.AccountId,
		Location:  addEventRequest.Location,
		Type:      addEventRequest.Type,
	}
	err = accountService.AddEvent(e)
	if err != nil {
		return models.GenerateError(err.Error(), 500)
	}

	return models.GenerateResponse(map[string]string{"message": "OK"}, 200)
}

func init() {
	accountRepo, err := account.NewDynamoDB(os.Getenv("ACCOUNTTABLE"))
	if err != nil {
		log.Fatalf("new account repository: %v\n", err)
	}
	locationRepo, err := location.NewDynamoDB(os.Getenv("LOCATIONTABLE"))
	if err != nil {
		log.Fatalf("new location repository: %v\n", err)
	}
	accountService, err = account.NewService(accountRepo, locationRepo)
	if err != nil {
		log.Fatalf("new account service: %v\n", err)
	}
}

func main() {
	lambda.Start(Handler)
}
