package repository

import (
	"database/sql"
	"fmt"
	"os"
	"time"
)

type Ticket struct {
	Id       string        `json:"id"`
	Location string        `json:"location"`
	Start    time.Time     `json:"start"`
	Stop     time.Time     `json:"stop,omitempty"`
	Duration time.Duration `json:"duration,omitempty"`
	Payment  string        `json:"payment,omitempty"`
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
	SELECT * FROM tickets WHERE location=$1 AND id =$2`
	var ticket Ticket
	err := pr.db.QueryRow(getStatement, location, id).Scan(&ticket)
	if err != nil {
		return Ticket{}, err
	}
	return ticket, nil
}

func (pr PostgresRepository) SaveTicket(t Ticket) (Ticket, error) {
	insertStatement := `
	INSERT INTO tickets(id, location, start, stop, duration, payment) VALUES ($1, $2, $3, $4, $5, $6) ON CONSTRAINT location_id
	DO UPDATE SET start=$3, stop=$4, duration=$5, payment=$6`
	var ticket Ticket
	err := pr.db.QueryRow(insertStatement, t.Id, t.Location, t.Start, t.Stop, t.Duration, t.Payment).Scan(&ticket)
	if err != nil {
		return Ticket{}, err
	}
	return ticket, nil
}
