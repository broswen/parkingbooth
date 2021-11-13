package main

import (
	"context"
	"encoding/json"
	"log"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/broswen/parkingbooth/location"
	"github.com/broswen/parkingbooth/models"
	"github.com/broswen/parkingbooth/ticket"
)

var ticketService *ticket.Service

func Handler(ctx context.Context, event events.APIGatewayProxyRequest) (models.APIGResponse, error) {
	var payTicketRequest models.PayTicketRequest
	err := json.Unmarshal([]byte(event.Body), &payTicketRequest)
	if err != nil {
		return models.GenerateError(err.Error(), 400)
	}

	ticket, err := ticketService.PayTicket(payTicketRequest.Location, payTicketRequest.Id, payTicketRequest.Payment)
	if err != nil {
		return models.GenerateError(err.Error(), 500)
	}

	return models.GenerateResponse(ticket, 200)
}

func init() {
	ticketRepo, err := ticket.NewDynamoDB(os.Getenv("TICKETTABLE"))
	if err != nil {
		log.Fatalf("new repository: %v\n", err)
	}
	locationRepo, err := location.NewDynamoDB(os.Getenv("LOCATIONTABLE"))
	if err != nil {
		log.Fatalf("new location repository: %v\n", err)
	}
	ticketService, err = ticket.NewService(ticketRepo, locationRepo)
	if err != nil {
		log.Fatalf("new ticket service: %v\n", err)
	}
}

func main() {
	lambda.Start(Handler)
}
