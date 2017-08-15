package main

import (
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
	_ "github.com/lib/pq"
)

var host string

var currRes Response

var recallTexts []string
var res Response

func main() {
	envvars := getEnv()
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

	currWave := 2
	res = Response{wave: currWave}
	currRes = Response{
		wave:      currWave + 1,
		questions: []Question{Question{}, Question{}, Question{}, Question{}},
	}

	if err := pgManager.GetRandomResponse(&res); err != nil {
		log.Println("Error querying random row from postgres ", err)
	}

	/*if err := pgManager.MarkResponseAsUsed(res.id); err != nil {
		log.Println("Error marking id as used from postgres ", res.id, err)
	}*/

	addr := ":" + envPort
	if envPort == "" || envPort == "8080" {
		addr = ":8080"
		host = "http://localhost" + addr
	} else {
		host = "https://social-transmission.herokuapp.com"
	}

	texts := [4]string{"According to a recent poll, <span id=\"s1\"></span>% of Americans say that they would prefer working under a male boss. <span id=\"s2\"></span>% of Americans would prefer to work under a female boss.",
		"Same sex marriage is a contested topic among Americans. In a poll conducted by the Pew Research Center, <span id=\"s1\"></span>% of respondents reported favoring same-sex marriage. <span id=\"s2\"></span>% reported opposing same-sex marriage.",
		"According to a database compiled by The Washington Post, 963 individuals where killed by police in 2016. Of those shot, <span id=\"s1\"></span> individuals were white and <span id=\"s2\"></span> were black.",
		"In 2007, <span id=\"s1\"></span> Mexican immigrants lived in the United States. In 2014, <span id=\"s2\"></span> Mexican immigrants lived in the US. Mexican immigrants have been at the center of one of the largest mass migrations in modern history."}

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
