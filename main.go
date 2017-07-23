package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"text/template"
)

type templateHandler struct {
	once     sync.Once
	filename string
	text     string
	s1       string
	s2       string
	templ    *template.Template
}

func (t *templateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	t.once.Do(func() {
		t.templ = template.Must(template.ParseFiles(filepath.Join("templates", t.filename)))
	})
	data := map[string]interface{}{
		"Text": t.text,
		"S1":   t.s1,
		"S2":   t.s2,
	}
	t.templ.Execute(w, data)
}

func main() {
	items := make(map[string]string)
	for _, item := range os.Environ() {
		splits := strings.Split(item, "=")
		key := splits[0]
		val := splits[1]
		items[key] = val
	}
	dbUrl := items["postgres_url"]

	db, err := sql.Open("postgres", dbUrl)
	if err != nil {
		log.Fatal("Error connecting to postgres", err)
	}

	wave := 2
	rows, err := db.Query("SELECT * FROM Responses WHERE wave=$1 ORDER BY random() LIMIT 1", wave)
	if err != nil {
		log.Fatal("Error querying random response from postgres", err)
	}
	defer rows.Close()

	var id int
	if err := rows.Scan(&id); err != nil {
		log.Fatal("Error getting id from row", err)
	}

	if _, err := db.Query("DELETE * FROM Responses WHERE id=$1", id); err != nil {
		log.Fatal("Error removing id from db", id, err)
	}

	var addr = flag.String("addr", ":8080", "The addr of the application.")
	flag.Parse()

	http.Handle("/", &templateHandler{
		filename: "intro.html",
		text:     "",
		s1:       "",
		s2:       "",
	})

	texts := [4]string{"Text1. S1: {{.S1}}. S2: {{.S2}}",
		"Text2. S1: {{.S1}}. S2: {{.S2}}",
		"Text3. S1: {{.S1}}. S2: {{.S2}}",
		"Text4. S1: {{.S1}}. S2: {{.S2}}"}

	var questions [4]question
	for i := 1; i <= 4; i++ {
		for j := 1; i <= 2; j++ {
			col := fmt.Sprintf("q%ds%d", i, j)
			if err := rows.Scan(&col); err != nil {
				log.Fatal("Error getting column from row", col, err)
			}
			questions[i] = question{text: texts[i]}
			if j == 1 {
				questions[i].s1 = col
			} else {
				questions[i].s2 = col
			}
		}
	}

	http.Handle("/q1", &templateHandler{
		filename: "question.html",
		text:     questions[0].text,
		s1:       questions[0].s1,
		s2:       questions[0].s2,
	})

	log.Println("Starting web server on", *addr)
	if err := http.ListenAndServe(*addr, nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
