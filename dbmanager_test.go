package main

import (
	"testing"
)

func TestOpenConnectionWithoutURL(t *testing.T) {
	db := dbmanager{}
	if success := db.OpenConnection(); success == false {
		t.Error("OpenConnection should not set the database manager's db without a url.")
	}
}

func TestOpenConnectionWithURL(t *testing.T) {
	db := dbmanager{url: getEnv()["postgres_url"]}
	if success := db.OpenConnection(); !success {
		t.Error("OpenConnection did not complete successfully.")
	} else if db.db == nil {
		t.Error("OpenConnection should set the database manager's db with a url.")
	}
}

func BenchmarkOpenConnection(b *testing.B) {
	db := dbmanager{url: getEnv()["postgres_url"]}
	for i := 0; i < b.N; i++ {
		db.OpenConnection()
	}
}

func TestCheckConnectionWithoutURL(t *testing.T) {
	db := dbmanager{}
	db.OpenConnection()
	success, _ := db.CheckConnection()
	if success {
		t.Error("CheckConnection should not be successful without a url.")
	}
}

func TestCheckConnectionWithURL(t *testing.T) {
	db := dbmanager{url: getEnv()["postgres_url"]}
	db.OpenConnection()
	success, err := db.CheckConnection()
	if err != nil {
		t.Error("CheckConnection should not return an error with a valid url.")
	} else if success == false {
		t.Error("CheckConnection should be successfull with a valid url.")
	}
}

func BenchmarkCheckConnection(b *testing.B) {
	db := dbmanager{url: getEnv()["postgres_url"]}
	db.OpenConnection()
	for i := 0; i < b.N; i++ {
		db.CheckConnection()
	}
}

func TestGetRandomResponseNoWave(t *testing.T) {
	db := dbmanager{url: getEnv()["postgres_url"]}
	db.OpenConnection()
	r := Response{}
	if err := db.GetRandomResponse(&r, false); err == nil {
		t.Error("GetRandomResponse error should not be nil without a valid wave.")
	}
}

func TestGetRandomResponseWithNegativeWave(t *testing.T) {
	db := dbmanager{url: getEnv()["postgres_url"]}
	db.OpenConnection()
	r := Response{wave: -1}
	if err := db.GetRandomResponse(&r, false); err == nil {
		t.Error("GetRandomResponse error should not be nil with a negative wave.")
	}
}

func TestGetRandomResponseWithLargeWave(t *testing.T) {
	db := dbmanager{url: getEnv()["postgres_url"]}
	db.OpenConnection()
	r := Response{wave: 1000}
	if err := db.GetRandomResponse(&r, false); err == nil {
		t.Error("GetRandomResponse error should not be nil with too large of a wave.")
	}
}

func TestGetRandomResponseWithValidWave(t *testing.T) {
	db := dbmanager{url: getEnv()["postgres_url"]}
	db.OpenConnection()
	r := Response{wave: 2}
	if err := db.GetRandomResponse(&r, false); err != nil {
		t.Error("GetRandomResponse should not return an error with a valid wave.")
	} else if r.id < 0 {
		t.Error("GetRandomResponse should populate the given response's id with a non-negative integer.")
	} else if r.condition < 0 {
		t.Error("GetRandomResponse should populate the given response's condition with a non-negative integer.")
	} else if len(r.targets) != 4 {
		t.Error("GetRandomResponse should populate the given response with 4 target questions.")
	} else {
		for i := 0; i < len(r.questions); i++ {
			q := r.questions[i]
			if q.number != int8(i+1) {
				t.Errorf("GetRandomResponse should populate response questions with the appropriate number. Question no: %d, i: %d.", int(q.number), i)
			} else if q.s1 == "" || q.s2 == "" {
				t.Error("GetRandomResponse should populate response questions with non-empty statistics.")
			}
		}
	}
}

func BenchmarkGetRandomResponse(b *testing.B) {
	db := dbmanager{url: getEnv()["postgres_url"]}
	db.OpenConnection()
	for i := 0; i < b.N; i++ {
		res := Response{wave: 2}
		db.GetRandomResponse(&res, false)
	}
}
