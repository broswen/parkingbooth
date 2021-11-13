package location

import (
	"context"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type Location struct {
	Id          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type LocationRepository interface {
	GetLocation(id string) (Location, error)
	SaveLocation(l Location) (Location, error)
}

type DynamoDBRepository struct {
	TableName string
	ddb       *dynamodb.Client
}

func NewDynamoDB(table string) (DynamoDBRepository, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
	}
	ddbClient := dynamodb.NewFromConfig(cfg)

	return DynamoDBRepository{
		ddb:       ddbClient,
		TableName: table,
	}, nil
}

func NewMap() (LocationRepository, error) {
	return MapRepository{
		m: make(map[string]Location, 0),
	}, nil
}

type MapRepository struct {
	m map[string]Location
}

func (mr MapRepository) GetLocation(id string) (Location, error) {
	l, ok := mr.m[id]
	if !ok {
		return Location{}, fmt.Errorf("location not found")
	}
	return l, nil
}

func (mr MapRepository) SaveLocation(l Location) (Location, error) {
	mr.m[l.Id] = l
	return l, nil
}

func (dr DynamoDBRepository) GetLocation(id string) (Location, error) {
	getItemInput := &dynamodb.GetItemInput{
		TableName: &dr.TableName,
		Key: map[string]types.AttributeValue{
			"PK": &types.AttributeValueMemberS{Value: fmt.Sprintf("L#%s", id)},
		},
	}

	getItemResult, err := dr.ddb.GetItem(context.Background(), getItemInput)
	if err != nil {
		return Location{}, err
	}
	if getItemResult.Item == nil {
		return Location{}, fmt.Errorf("location not found")
	}

	l := Location{
		Id:          getItemResult.Item["id"].(*types.AttributeValueMemberS).Value,
		Name:        getItemResult.Item["name"].(*types.AttributeValueMemberS).Value,
		Description: getItemResult.Item["description"].(*types.AttributeValueMemberS).Value,
	}

	return l, nil
}

func (dr DynamoDBRepository) SaveLocation(l Location) (Location, error) {
	putItemInput := &dynamodb.PutItemInput{
		TableName: &dr.TableName,
		Item: map[string]types.AttributeValue{
			"PK":          &types.AttributeValueMemberS{Value: fmt.Sprintf("L#%s", l.Id)},
			"id":          &types.AttributeValueMemberS{Value: l.Id},
			"name":        &types.AttributeValueMemberS{Value: l.Name},
			"description": &types.AttributeValueMemberS{Value: l.Description},
		},
	}
	_, err := dr.ddb.PutItem(context.Background(), putItemInput)
	if err != nil {
		return Location{}, err
	}
	return l, nil
}
