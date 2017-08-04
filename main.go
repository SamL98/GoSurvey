package main

import (
	"flag"
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
var demoData map[string]map[string]interface{}
var knowledgeData map[string]map[string]interface{}

func main() {
	populateDemoData()
	populateKnowledgeData()

	items := make(map[string]string)
	for _, item := range os.Environ() {
		splits := strings.Split(item, "=")
		key := splits[0]
		val := splits[1]
		items[key] = val
	}
	dbURL := items["postgres_url"]

	pgManager = dbmanager{url: dbURL}
	pgManager.OpenConnection()
	defer pgManager.db.Close()

	if success, err := pgManager.CheckConnection(); !success || err != nil {
		log.Fatal("Error pinging postgres ", err)
	}
	log.Println("Successfully connected to Postgres.")

	if err := pgManager.GetAllIPAddresses(); err != nil {
		log.Fatal("Error getting IP addresses from postres ", err)
	}

	currWave := 2
	res = Response{wave: currWave}
	currRes = Response{
		wave: currWave + 1,
		questions: []Question{
			Question{}, Question{},
			Question{}, Question{},
		},
		demographic: Demographics{},
	}

	if err := pgManager.GetRandomResponse(&res); err != nil {
		log.Fatal("Error querying random row from postgres ", err)
	}

	/*if err := pgManager.DeleteRow(res.id); err != nil {
		log.Fatal("Error deleting id from postgres ", res.id, err)
	}*/

	var addr = flag.String("addr", ":8080", "The addr of the application.")
	flag.Parse()

	texts := [4]string{"Text1. S1: <span id=\"s1\"></span>. S2: <span id=\"s2\"></span>",
		"Text2. S1: <span id=\"s1\"></span>. S2: <span id=\"s2\"></span>",
		"Text3. S1: <span id=\"s1\"></span>. S2: <span id=\"s2\"></span>",
		"Text4. S1: <span id=\"s1\"></span>. S2: <span id=\"s2\"></span>"}

	for i := range res.questions {
		res.questions[i].text = texts[i]
	}

	recallTexts = []string{"what percentage of Americans prefer to have a <u>female</u> boss?",
		"what percentage of Americans prefer to have a <u>male</u> boss?",
		"what percentage of Americans polled in 2015 reported <u>favoring</u>&nbsp;allowing gays and lesbians to marry?",
		"According to one of the blog posts, what percentage of Americans polled in 2015 reported <u>opposing</u>&nbsp;allowing gays and lesbians to marry?",
		"According to one of the blog posts, how many millions of Mexican immigrants lived in the United States in <u>2007</u>?",
		"According to one of the blog posts, how many millions of Mexican immigrants lived in the United States in <u>2014</u>?",
		"According to one of the blog posts, in 2016, 963 individuals were shot and killed by police. Of those shot, how many individuals were <u>white</u>?",
		"According to one of the blog posts, in 2016, 963 individuals were shot and killed by police. Of those shot, how many individuals were <u>black</u>?"}

	r := httprouter.New()
	r.GET("/", Intro)
	r.GET("/instructions", Instructions)
	r.GET("/question/:q", QuestionHandler)
	r.GET("/recall_instructions", RecallInstructions)
	r.GET("/recall_question/:q", RecallQuestion)
	r.GET("/recall_strength/:q", RecallStrength)
	r.GET("/demographics/:type/:q", DemographicsQuestion)
	r.GET("/knowledge_instructions", KnowledgeInstructions)
	r.GET("/knowledge_question/:name", KnowledgeQuestion)
	r.GET("/completed", CompletionHandler)

	r.POST("/report_interest/:q", ParseInterest)
	r.POST("/report_recall_statistic/:q", ParseStatistic)
	r.POST("/report_recall_strength/:q", ParseRecallStrength)
	r.POST("/report_mcdemographic/:name", ParseMCDemographic)
	r.POST("report_textdemographic/:name", ParseTextDemographic)
	r.POST("/report_knowledge/:name", ParseKnowledge)

	log.Println("Starting web server on", *addr)
	if err := http.ListenAndServe(*addr, r); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
