package models

import (
	"database/sql"
	"example.com/rest-api/db"
	"time"
)

type Event struct {
	ID          int64     `json:"id" `
	Name        string    `json:"name" binding:"required"`
	Description string    `json:"description" binding:"required"`
	Location    string    `json:"location" binding:"required"`
	DateTime    time.Time `json:"date_time" binding:"required"`
	UserID      int64     `json:"user_id"`
}

func (e *Event) Save() error {
	query := `INSERT INTO events(name, description, location, date_time, user_id)
VALUES (?,?,?,?,?) 
    `
	statement, err := db.DB.Prepare(query)
	defer func(statement *sql.Stmt) {
		err := statement.Close()
		if err != nil {
			println(err.Error())
		}
	}(statement)
	result, err := statement.Exec(e.Name, e.Description, e.Location, e.DateTime, e.UserID)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	e.ID = id
	return err
}

func GetAllEvents() ([]Event, error) {
	query := "SELECT * FROM events"
	rows, err := db.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var events []Event
	for rows.Next() {
		var event Event
		rows.Scan(&event.ID, &event.Name, &event.Description, &event.Location, &event.DateTime, &event.UserID)
		events = append(events, event)
	}
	return events, nil
}

func GetEventById(id int64) (*Event, error) {
	query := "SELECT * FROM events WHERE id = ?"
	row := db.DB.QueryRow(query, id)
	var event Event
	err := row.Scan(&event.ID, &event.Name, &event.Description, &event.Location, &event.DateTime, &event.UserID)
	if err != nil {
		return nil, err
	}

	return &event, nil
}

func (e Event) Update() error {
	query := `UPDATE events
SET name=?, description=?, location=?, date_time=?
WHERE id = ?`

	statement, err := db.DB.Prepare(query)
	if err != nil {
		return err
	}
	defer statement.Close()

	_, err = statement.Exec(e.Name, e.Description, e.Location, e.DateTime, e.ID)
	return err
}

func (e Event) Delete() error {
	query := "DELETE FROM events WHERE id = ?"
	statement, err := db.DB.Prepare(query)
	if err != nil {
		return err
	}
	defer statement.Close()
	_, err = statement.Exec(e.ID)
	return err
}

func (e Event) Register(userId int64) error {
	query := "INSERT INTO registrations(event_id, user_id) VALUES (?,?)"
	statement, err := db.DB.Prepare(query)

	if err != nil {
		return err
	}

	defer statement.Close()

	_, err = statement.Exec(e.ID, userId)

	return err
}

func (e Event) Unregister(userId int64) error {
	query := "DELETE FROM registrations WHERE event_id = ? AND user_id = ?"
	statement, err := db.DB.Prepare(query)

	if err != nil {
		return err
	}

	defer statement.Close()

	_, err = statement.Exec(e.ID, userId)
	return err
}
