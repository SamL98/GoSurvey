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

func (m *dbmanager) OpenConnection() {
	db, err := sql.Open("postgres", m.url)
	if err != nil {
		log.Fatal("Error connecting to postgres ", err)
	}
	m.db = db
}

func (m *dbmanager) CheckConnection() (success bool, err error) {
	if err := m.db.Ping(); err != nil {
		return false, err
	}
	return true, nil
}

func (m *dbmanager) GetAllIPAddresses() error {
	rows, err := m.db.Query("SELECT ip FROM Responses")
	defer rows.Close()

	if err != nil {
		return err
	}

	for rows.Next() {
		var ip string
		if err := rows.Scan(&ip); err != nil {
			log.Fatalln("Error scanning row for ip address.")
		}
		ips = append(ips, ip)
	}
	return nil
}

func (m *dbmanager) GetRandomResponse(r *Response) error {
	rows, err := m.db.Query("SELECT wave, id FROM Responses WHERE wave=$1 ORDER BY random() LIMIT 1", r.wave)
	defer rows.Close()

	if err != nil {
		return err
	}

	if !rows.Next() {
		return errors.New("zero rows from random response query")
	}

	if err := rows.Scan(&r.wave, &r.id); err != nil {
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
			log.Fatalln("Error scanning question row ", err)
		}
		questions = append(questions, q)
		i++
	}
	r.questions = questions

	return nil
}

func (m *dbmanager) DeleteRow(id int) error {
	if _, err := m.db.Query("DELETE FROM Responses WHERE id=$1", id); err != nil {
		return err
	}
	return nil
}
