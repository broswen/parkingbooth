package repository

import (
	"database/sql"
	"fmt"
	"os"

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

func NewRepository(t string) (TicketRepository, error) {
	switch t {
	case "map":
		return MapRepository{
			m: make(map[string]map[string]Ticket, 0),
		}, nil
	case "postgres":
		host := os.Getenv("POSTGRES_HOST")
		port := os.Getenv("POSTGRES_PORT")
		user := os.Getenv("POSTGRES_USER")
		password := os.Getenv("POSTGRES_PASS")
		dbname := os.Getenv("POSTGRES_DB")
		connString := fmt.Sprintf("host=%s port=%s user=%s "+
			"password=%s dbname=%s sslmode=disable",
			host, port, user, password, dbname)

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

	default:
		return nil, fmt.Errorf("unknown repository type: %v\n", t)
	}
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
