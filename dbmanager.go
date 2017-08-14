package main

import (
	"database/sql"
	"log"

	"errors"

	_ "github.com/lib/pq"
)

type dbmanager struct {
	url string
	db  *sql.DB
}

var pgManager dbmanager

func (m *dbmanager) OpenConnection() bool {
	db, err := sql.Open("postgres", m.url)
	if err != nil {
		log.Println("Error connecting to postgres ", err)
		return false
	}
	m.db = db
	return true
}

func (m *dbmanager) CheckConnection() (success bool, err error) {
	if err := m.db.Ping(); err != nil {
		return false, err
	}
	return true, nil
}

func (m *dbmanager) GetRandomResponse(r *Response) error {
	rows, err := m.db.Query("SELECT wave, id, condition FROM Responses WHERE wave=$1 AND used=false ORDER BY random() LIMIT 1", r.wave)
	defer rows.Close()

	if err != nil {
		return err
	}

	if !rows.Next() {
		return errors.New("zero rows from random response query")
	}

	if err := rows.Scan(&r.wave, &r.id, &r.condition); err != nil {
		return err
	}

	rows, err = m.db.Query("SELECT s1, s2 FROM Questions WHERE response=$1", r.id)
	defer rows.Close()

	if err != nil {
		return err
	}
	questions := []Question{}
	i := 1
	for rows.Next() {
		q := Question{number: int8(i)}
		if err := rows.Scan(&q.s1, &q.s2); err != nil {
			log.Println("Error scanning question row ", err)
		}
		questions = append(questions, q)
		i++
	}
	r.questions = questions

	return nil
}

func (m *dbmanager) MarkResponseAsUsed(id int) error {
	if _, err := m.db.Query("UPDATE Responses SET used = true WHERE id=$1", id); err != nil {
		return err
	}
	return nil
}

func (m *dbmanager) AddResponse(r *Response, seed int) error {
	rows, err := m.db.Query("INSERT INTO Responses (wave, used, seed) VALUES ($1, false, $2) RETURNING id", r.wave, seed)
	if err != nil {
		return err
	}
	defer rows.Close()
	rows.Next()

	var id int
	if err := rows.Scan(&id); err != nil {
		return err
	}

	for i := 0; i < len(r.questions); i++ {
		q := r.questions[i]
		log.Println(q.number)
		if _, err := m.db.Query("INSERT INTO Questions (response, time, number, interest) VALUES ($1, $2, $3, $4)", id, q.time, int(q.number)-1, q.interest); err != nil {
			return err
		}
	}
	return nil
}
