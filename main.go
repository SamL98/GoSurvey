package main

import (
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
	_ "github.com/lib/pq"
)

var host string

var recallTexts []string

var envvars map[string]string

func main() {
	envvars = getEnv()
	dbURL := envvars["postgres_url"]
	envPort := envvars["PORT"]

	pgManager = dbmanager{url: dbURL}
	if success := pgManager.OpenConnection(); !success {
		log.Fatal("Could not open connection to postgres")
	}
	defer pgManager.db.Close()

	if success, err := pgManager.CheckConnection(); !success || err != nil {
		log.Fatal("Error pinging postgres ", err)
	}
	log.Println("Successfully connected to Postgres.")

	currWave = 2
	responses = make(map[int][]Response)
	sessionID = 0

	addr := ":" + envPort
	if envPort == "" || envPort == "8080" {
		addr = ":8080"
		host = "http://localhost" + addr
	} else {
		host = "https://sotrapp.herokuapp.com"
	}

	r := httprouter.New()
	r.GET("/", Intro)
	r.GET("/instructions", Instructions)
	r.GET("/question/:q", QuestionHandler)
	r.GET("/completed", CompletionHandler)

	r.POST("/report_interest/:q", ParseInterest)

	log.Println("Starting web server on", addr)
	if err := http.ListenAndServe(addr, r); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
