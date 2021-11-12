package models

import (
	"encoding/json"

	"github.com/aws/aws-lambda-go/events"
)

type ErrorReponse struct {
	Message string `json:"message"`
}

type GenerateTicketRequest struct {
	Location string `json:"location"`
}

type CompleteTicketRequest struct {
	Location string `json:"location"`
	Id       string `json:"id"`
}

type PayTicketRequest struct {
	Location string `json:"location"`
	Id       string `json:"id"`
	Payment  string `json:"payment"`
}

type APIGResponse events.APIGatewayProxyResponse

type PostgresCredentials struct {
	Username             string `json:"username"`
	Password             string `json:"password"`
	Port                 int    `json:"port"`
	Host                 string `json:"host"`
	DBInstanceIdentifier string `json:"dbInstanceIdentifier"`
}

func GenerateResponse(i interface{}, status int) (APIGResponse, error) {

	body, err := json.Marshal(i)
	if err != nil {
		return APIGResponse{StatusCode: 500}, err
	}
	resp := APIGResponse{
		StatusCode:      status,
		IsBase64Encoded: false,
		Body:            string(body),
	}

	return resp, nil
}

func GenerateError(message string, status int) (APIGResponse, error) {

	response := ErrorReponse{
		Message: message,
	}

	body, err := json.Marshal(response)
	if err != nil {
		return APIGResponse{StatusCode: 500}, err
	}
	resp := APIGResponse{
		StatusCode:      status,
		IsBase64Encoded: false,
		Body:            string(body),
	}

	return resp, nil
}
