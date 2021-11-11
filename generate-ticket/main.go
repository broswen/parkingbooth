package main

import (
	"context"
	"encoding/json"
	"log"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	"github.com/broswen/parkingbooth/models"
	"github.com/broswen/parkingbooth/repository"
	"github.com/broswen/parkingbooth/ticket"
)

var ticketService *ticket.Service
var smClient *secretsmanager.Client

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
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
	}

	// Using the Config value, create the DynamoDB client
	smClient = secretsmanager.NewFromConfig(cfg)

	response, err := smClient.GetSecretValue(context.TODO(), &secretsmanager.GetSecretValueInput{
		SecretId: aws.String(os.Getenv("DB_SECRET_ARN")),
	})
	if err != nil {
		log.Fatalf("get secret value, %v", err)
	}
	var postgresCreds models.PostgresCredentials
	err = json.Unmarshal([]byte(*response.SecretString), &postgresCreds)
	if err != nil {
		log.Fatalf("marshal secret value, %v", err)
	}

	// TODO pass secret values from above
	repo, err := repository.NewPostgres(postgresCreds)
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
