package main

import (
	"context"
	"encoding/json"
	"log"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/broswen/parkingbooth/models"
	"github.com/broswen/parkingbooth/ticket"
)

var ticketService *ticket.Service

func Handler(ctx context.Context, event events.APIGatewayProxyRequest) (models.APIGResponse, error) {
	var generateTicketRequest models.GenerateTicketRequest
	err := json.Unmarshal([]byte(event.Body), &generateTicketRequest)
	if err != nil {
		return models.GenerateError(err.Error(), 400)
	}

	ticket, err := ticketService.GenerateTicket(generateTicketRequest.Location)
	if err != nil {
		return models.GenerateError(err.Error(), 500)
	}

	return models.GenerateResponse(ticket, 200)
}

func init() {
	repo, err := ticket.NewDynamoDB(os.Getenv("TICKETTABLE"))
	if err != nil {
		log.Fatalf("new repository: %v\n", err)
	}
	ticketService, err = ticket.NewService(repo)
	if err != nil {
		log.Fatalf("new ticket service: %v\n", err)
	}
}

func main() {
	lambda.Start(Handler)
}
