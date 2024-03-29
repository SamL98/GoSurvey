package main

import (
	"database/sql"
	"log"
	"time"

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

func (m *dbmanager) CloseConnection() {
	if err := m.db.Close(); err != nil {
		log.Println("Error closing postgres ", err)
	}
}

func (m *dbmanager) CheckConnection() (success bool, err error) {
	if err := m.db.Ping(); err != nil {
		return false, err
	}
	return true, nil
}

func (m *dbmanager) GetRandomResponse(r *Response, used bool) error {
	rows, err := m.db.Query("SELECT wave, id, condition FROM Responses WHERE wave=$1 AND used=$2 ORDER BY random() LIMIT 1", r.wave, used)
	if err != nil {
		return err
	}
	defer rows.Close()

	if !rows.Next() {
		if used {
			return errors.New("zero rows from random response query")
		}
		log.Println("All conditions are used, querying for used one.")
		return m.GetRandomResponse(r, true)
	}

	if err := rows.Scan(&r.wave, &r.id, &r.condition); err != nil {
		return err
	}

	log.Println("Wave:", r.wave, "Id:", r.id, "Condition:", r.condition, time.Now())
	if r.condition == 0 {
		return errors.New("response condition should not be zero")
	}

	rows, err = m.db.Query("SELECT s1, s2 FROM Questions WHERE response=$1", r.id)
	if err != nil {
		return err
	}
	defer rows.Close()

	questions := []Question{}
	i := 1
	for rows.Next() {
		q := Question{number: int8(i), distractor: false}
		if err := rows.Scan(&q.s1, &q.s2); err != nil {
			log.Println("Error scanning question row ", err)
		}
		questions = append(questions, q)
		i++
	}
	r.targets = questions

	return nil
}

func (m *dbmanager) MarkResponseAsUsed(id int) error {
	rows, err := m.db.Query("UPDATE Responses SET used = true WHERE id=$1", id)
	if err != nil {
		return err
	}
	if rows != nil {
		rows.Close()
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

	for i := 0; i < len(r.targets); i++ {
		q := r.targets[i]
		insRows, err := m.db.Query("INSERT INTO Questions (response, time, number, interest) VALUES ($1, $2, $3, $4)", id, q.time, int(q.number)-1, q.interest)
		if insRows != nil {
			insRows.Close()
		}
		if err != nil {
			return err
		}
	}
	return nil
}
