package account

import (
	"context"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type Account struct {
	Id    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

type AccountEventType string

const (
	InEvent  AccountEventType = "IN"
	OutEvent AccountEventType = "OUT"
)

type AccountEvent struct {
	Id        string           `json:"id"`
	AccountId string           `json:"accountId"`
	Type      AccountEventType `json:"type"`
	Location  string           `json:"location"`
	Time      int64            `json:"time"`
}

type AccountRepository interface {
	GetAccount(id string) (Account, error)
	AddEvent(e AccountEvent) error
	CreateAccount(a Account) (Account, error)
	UpdateAccount(a Account) (Account, error)
	DeleteAccount(id string) error
}

type MapRepository struct {
	m map[string]Account
	e map[string][]AccountEvent
}

type DynamoDBRepository struct {
	TableName string
	ddb       *dynamodb.Client
}

func NewMap() (AccountRepository, error) {
	return MapRepository{
		m: make(map[string]Account, 0),
		e: make(map[string][]AccountEvent, 0),
	}, nil
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

func (mr MapRepository) CreateAccount(a Account) (Account, error) {
	_, ok := mr.m[a.Id]
	if ok {
		return Account{}, fmt.Errorf("account with id already exists")
	}
	mr.m[a.Id] = a
	mr.e[a.Id] = make([]AccountEvent, 0)
	return a, nil
}

func (mr MapRepository) GetAccount(id string) (Account, error) {
	a, ok := mr.m[id]
	if !ok {
		return Account{}, fmt.Errorf("account doesn't exist")
	}
	return a, nil
}

func (mr MapRepository) UpdateAccount(a Account) (Account, error) {
	_, ok := mr.m[a.Id]
	if !ok {
		return Account{}, fmt.Errorf("account doesn't exist")
	}
	mr.m[a.Id] = a
	return a, nil
}

func (mr MapRepository) DeleteAccount(id string) error {
	delete(mr.m, id)
	delete(mr.e, id)
	return nil
}

func (mr MapRepository) AddEvent(e AccountEvent) error {
	events, ok := mr.e[e.AccountId]
	if !ok {
		return fmt.Errorf("account doesn't exist")
	}
	mr.e[e.AccountId] = append(events, e)
	return nil
}

func (dr DynamoDBRepository) CreateAccount(a Account) (Account, error) {
	putItemInput := &dynamodb.PutItemInput{
		TableName: &dr.TableName,
		Item: map[string]types.AttributeValue{
			"PK":    &types.AttributeValueMemberS{Value: fmt.Sprintf("A#%s", a.Id)},
			"SK":    &types.AttributeValueMemberS{Value: fmt.Sprintf("A#%s", a.Id)},
			"id":    &types.AttributeValueMemberS{Value: a.Id},
			"name":  &types.AttributeValueMemberS{Value: a.Name},
			"email": &types.AttributeValueMemberS{Value: a.Email},
		},
		ConditionExpression: aws.String("attribute_not_exist(PK)"),
	}

	_, err := dr.ddb.PutItem(context.Background(), putItemInput)
	if err != nil {
		return Account{}, err
	}

	return a, nil
}

func (dr DynamoDBRepository) GetAccount(id string) (Account, error) {
	getItemInput := &dynamodb.GetItemInput{
		TableName: &dr.TableName,
		Key: map[string]types.AttributeValue{
			"PK": &types.AttributeValueMemberS{Value: fmt.Sprintf("A#%s", id)},
			"SK": &types.AttributeValueMemberS{Value: fmt.Sprintf("A#%s", id)},
		},
	}

	getItemResult, err := dr.ddb.GetItem(context.Background(), getItemInput)
	if err != nil {
		return Account{}, err
	}
	if getItemResult.Item == nil {
		return Account{}, fmt.Errorf("account not found")
	}

	a := Account{
		Id:    getItemResult.Item["id"].(*types.AttributeValueMemberS).Value,
		Name:  getItemResult.Item["name"].(*types.AttributeValueMemberS).Value,
		Email: getItemResult.Item["email"].(*types.AttributeValueMemberS).Value,
	}

	return a, nil
}

func (dr DynamoDBRepository) UpdateAccount(a Account) (Account, error) {
	updateItemInput := &dynamodb.UpdateItemInput{
		TableName: &dr.TableName,
		Key: map[string]types.AttributeValue{
			"PK": &types.AttributeValueMemberS{Value: fmt.Sprintf("A#%s", a.Id)},
			"SK": &types.AttributeValueMemberS{Value: fmt.Sprintf("A#%s", a.Id)},
		},
		UpdateExpression: aws.String("SET #n = :n, #e = :e"),
		ExpressionAttributeNames: map[string]string{
			"#n": "name",
			"#e": "email",
		},
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":n": &types.AttributeValueMemberS{Value: a.Name},
			":e": &types.AttributeValueMemberS{Value: a.Email},
		},
		ConditionExpression: aws.String("attribute_exists(PK)"),
	}

	_, err := dr.ddb.UpdateItem(context.Background(), updateItemInput)
	if err != nil {
		return Account{}, err
	}

	return a, nil
}

func (dr DynamoDBRepository) DeleteAccount(id string) error {
	deleteItemInput := &dynamodb.DeleteItemInput{
		TableName: &dr.TableName,
		Key: map[string]types.AttributeValue{
			"PK": &types.AttributeValueMemberS{Value: fmt.Sprintf("A#%s", id)},
			"SK": &types.AttributeValueMemberS{Value: fmt.Sprintf("A#%s", id)},
		},
	}

	_, err := dr.ddb.DeleteItem(context.Background(), deleteItemInput)
	if err != nil {
		return err
	}
	return nil
}

func (dr DynamoDBRepository) AddEvent(e AccountEvent) error {
	putItemInput := &dynamodb.PutItemInput{
		TableName: &dr.TableName,
		Item: map[string]types.AttributeValue{
			"PK":        &types.AttributeValueMemberS{Value: fmt.Sprintf("A#%s", e.AccountId)},
			"SK":        &types.AttributeValueMemberS{Value: fmt.Sprintf("E#%s", e.Id)},
			"id":        &types.AttributeValueMemberS{Value: e.Id},
			"accountid": &types.AttributeValueMemberS{Value: e.AccountId},
			"type":      &types.AttributeValueMemberS{Value: string(e.Type)},
			"location":  &types.AttributeValueMemberS{Value: e.Location},
			"time":      &types.AttributeValueMemberN{Value: fmt.Sprintf("%d", e.Time)},
		},
	}

	_, err := dr.ddb.PutItem(context.Background(), putItemInput)
	if err != nil {
		return err
	}

	return nil
}
