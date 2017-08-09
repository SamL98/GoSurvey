package main

import (
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/julienschmidt/httprouter"
	_ "github.com/lib/pq"
)

var currRes Response

var recallTexts []string
var res Response

func main() {
	items := make(map[string]string)
	for _, item := range os.Environ() {
		splits := strings.Split(item, "=")
		key := splits[0]
		val := splits[1]
		items[key] = val
	}
	dbURL := items["postgres_url"]
	envPort := items["PORT"]

	pgManager = dbmanager{url: dbURL}
	pgManager.OpenConnection()
	defer pgManager.db.Close()

	if success, err := pgManager.CheckConnection(); !success || err != nil {
		log.Fatal("Error pinging postgres ", err)
	}
	log.Println("Successfully connected to Postgres.")

	currWave := 2
	res = Response{wave: currWave}
	currRes = Response{
		wave:      currWave + 1,
		questions: []Question{Question{}, Question{}, Question{}, Question{}},
	}

	if err := pgManager.GetRandomResponse(&res); err != nil {
		log.Fatal("Error querying random row from postgres ", err)
	}

	/*if err := pgManager.MarkResponseAsUsed(res.id); err != nil {
		log.Fatal("Error marking id as used from postgres ", res.id, err)
	}*/

	addr := ":" + envPort
	if envPort == "" {
		addr = ":8080"
	}

	texts := [4]string{"Text1. S1: <span id=\"s1\"></span>. S2: <span id=\"s2\"></span>",
		"Text2. S1: <span id=\"s1\"></span>. S2: <span id=\"s2\"></span>",
		"Text3. S1: <span id=\"s1\"></span>. S2: <span id=\"s2\"></span>",
		"Text4. S1: <span id=\"s1\"></span>. S2: <span id=\"s2\"></span>"}

	for i := range res.questions {
		res.questions[i].text = texts[i]
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
