package ticket

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"strconv"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/broswen/parkingbooth/models"
	_ "github.com/lib/pq"
)

type Ticket struct {
	Id       string `json:"id"`
	Location string `json:"location"`
	Start    int64  `json:"start"`
	Stop     int64  `json:"stop,omitempty"`
	Duration int64  `json:"duration,omitempty"`
	Payment  string `json:"payment,omitempty"`
}

type TicketRepository interface {
	GetTicket(location, id string) (Ticket, error)
	SaveTicket(t Ticket) (Ticket, error)
}

type MapRepository struct {
	m map[string]map[string]Ticket
}

type PostgresRepository struct {
	db *sql.DB
}

type DynamoDBRepository struct {
	TableName string
	ddb       *dynamodb.Client
}

func NewPostgres(creds models.PostgresCredentials) (TicketRepository, error) {
	connString := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		creds.Host, creds.Port, creds.Username, creds.Password, "postgres")

	db, err := sql.Open("postgres", connString)
	if err != nil {
		return nil, fmt.Errorf("sql open: %v\n", err)
	}

	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("db ping: %v\n", err)
	}

	return PostgresRepository{
		db: db,
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

func NewMap() (TicketRepository, error) {
	return MapRepository{
		m: make(map[string]map[string]Ticket, 0),
	}, nil
}

func (mr MapRepository) GetTicket(location, id string) (Ticket, error) {
	l, ok := mr.m[location]
	if !ok {
		return Ticket{}, fmt.Errorf("no tickets for location: %v\n", location)
	}
	t, ok := l[id]
	if !ok {
		return Ticket{}, fmt.Errorf("no tickets for id: %v\n", id)
	}
	return t, nil
}

func (mr MapRepository) SaveTicket(t Ticket) (Ticket, error) {
	l, ok := mr.m[t.Location]
	if !ok {
		mr.m[t.Location] = make(map[string]Ticket, 0)
	}
	l = mr.m[t.Location]
	l[t.Id] = t
	return t, nil
}

func (pr PostgresRepository) GetTicket(location, id string) (Ticket, error) {
	getStatement := `
	SELECT * FROM tickets WHERE location_id=$1 AND id=$2`
	var ticket Ticket
	err := pr.db.QueryRow(getStatement, location, id).Scan(&ticket.Id, &ticket.Location, &ticket.Start, &ticket.Stop, &ticket.Duration, &ticket.Payment)
	if err != nil {
		return Ticket{}, err
	}
	return ticket, nil
}

func (pr PostgresRepository) SaveTicket(t Ticket) (Ticket, error) {
	insertStatement := `
	INSERT INTO tickets(id, location_id, start_epoch, stop_epoch, duration_seconds, payment_id) VALUES($1, $2, $3, $4, $5, $6) ON CONFLICT ON CONSTRAINT tickets_id_location_id_key
	DO UPDATE SET start_epoch=$3, stop_epoch=$4, duration_seconds=$5, payment_id=$6`
	_, err := pr.db.Exec(insertStatement, t.Id, t.Location, t.Start, t.Stop, t.Duration, t.Payment)
	if err != nil {
		return Ticket{}, err
	}
	return t, nil
}

func (dr DynamoDBRepository) GetTicket(location, id string) (Ticket, error) {
	getItemInput := &dynamodb.GetItemInput{
		TableName: &dr.TableName,
		Key: map[string]types.AttributeValue{
			"PK": &types.AttributeValueMemberS{Value: fmt.Sprintf("T#%s", id)},
			"SK": &types.AttributeValueMemberS{Value: fmt.Sprintf("T#%s", id)},
		},
	}

	getItemResult, err := dr.ddb.GetItem(context.Background(), getItemInput)
	if err != nil {
		return Ticket{}, err
	}
	if getItemResult.Item == nil {
		return Ticket{}, fmt.Errorf("ticket not found")
	}

	start, err := strconv.ParseInt(getItemResult.Item["start"].(*types.AttributeValueMemberN).Value, 10, 64)
	if err != nil {
		return Ticket{}, err
	}
	stop, err := strconv.ParseInt(getItemResult.Item["stop"].(*types.AttributeValueMemberN).Value, 10, 64)
	if err != nil {
		return Ticket{}, err
	}
	duration, err := strconv.ParseInt(getItemResult.Item["duration"].(*types.AttributeValueMemberN).Value, 10, 64)
	if err != nil {
		return Ticket{}, err
	}
	ticket := Ticket{
		Id:       getItemResult.Item["id"].(*types.AttributeValueMemberS).Value,
		Location: getItemResult.Item["location"].(*types.AttributeValueMemberS).Value,
		Start:    start,
		Stop:     stop,
		Duration: duration,
		Payment:  getItemResult.Item["payment"].(*types.AttributeValueMemberS).Value,
	}

	return ticket, nil
}

func (dr DynamoDBRepository) SaveTicket(t Ticket) (Ticket, error) {
	putItemInput := &dynamodb.PutItemInput{
		TableName: &dr.TableName,
		Item: map[string]types.AttributeValue{
			"PK":       &types.AttributeValueMemberS{Value: fmt.Sprintf("T#%s", t.Id)},
			"SK":       &types.AttributeValueMemberS{Value: fmt.Sprintf("T#%s", t.Id)},
			"id":       &types.AttributeValueMemberS{Value: t.Id},
			"location": &types.AttributeValueMemberS{Value: t.Location},
			"start":    &types.AttributeValueMemberN{Value: fmt.Sprintf("%d", t.Start)},
			"stop":     &types.AttributeValueMemberN{Value: fmt.Sprintf("%d", t.Stop)},
			"duration": &types.AttributeValueMemberN{Value: fmt.Sprintf("%d", t.Duration)},
			"payment":  &types.AttributeValueMemberS{Value: t.Payment},
		},
	}
	_, err := dr.ddb.PutItem(context.Background(), putItemInput)
	if err != nil {
		return Ticket{}, err
	}
	return t, nil
}
