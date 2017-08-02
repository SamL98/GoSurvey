package main

import (
	"log"
	"net/http"
	"path/filepath"
	"strconv"
	"sync"
	"text/template"

	"github.com/julienschmidt/httprouter"
)

type templateHandler struct {
	once     sync.Once
	filename string
	templ    *template.Template
}

func populateDemoData() {
	demoData = make(map[string]map[string]interface{})
	demoData["sex"] = map[string]interface{}{
		"Question": "What is your sex?",
		"Choices":  []string{"Male", "Female"},
		"NextType": "text",
		"NextName": "age",
	}
	demoData["age"] = map[string]interface{}{
		"Question": "How old are you?",
		"NextType": "text",
		"NextName": "race",
	}
	demoData["race"] = map[string]interface{}{
		"Question": "What is your race?",
		"NextType": "text",
		"NextName": "major",
	}
	demoData["major"] = map[string]interface{}{
		"Question": "What is/was your major in college?",
		"NextType": "text",
		"NextName": "year",
	}
	demoData["year"] = map[string]interface{}{
		"Question": "What year are you in college?",
		"NextType": "mc",
		"NextName": "affiliation",
	}
	demoData["affiliation"] = map[string]interface{}{
		"Question": "Generally speaking,  do you usually think of yourself as a Republican, a Democrat, an Independent or what?",
		"Choices":  []string{"Strong Democrat", "Democrat", "Lean Democrat", "Lean Republican", "Republican", "Strong Republican"},
		"NextType": "mc",
		"NextName": "ideology",
	}
	demoData["ideology"] = map[string]interface{}{
		"Question": "Generally speaking, do you consider yourself to be:",
		"Choices":  []string{"extemely conservative", "conservative", "slightly conservative", "moderate", "slightly liberal", "liberal", "extremely liberal"},
		"NextType": "mc",
		"NextName": "interest",
	}
	demoData["interest"] = map[string]interface{}{
		"Question": "How interested are you in what’s going on in American government and politics?",
		"Choices":  []string{"Very Interested", "Interested", "Somewhat Interested", "Not interested at all"},
		"NextType": "knowledge",
		"NextName": "majorities",
	}
}

func Intro(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	ip := getIPAdress(r)
	log.Println(ip)

	filename := "intro.html"
	if ipAddressExists(ip) {
		filename = "error.html"
	}

	t := templateHandler{filename: filename}
	t.once.Do(func() {
		t.templ = template.Must(template.ParseFiles(filepath.Join("templates", t.filename)))
	})
	t.templ.Execute(w, make(map[string]interface{}))
}

func Instructions(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	t := templateHandler{filename: "instructions.html"}
	t.once.Do(func() {
		t.templ = template.Must(template.ParseFiles(filepath.Join("templates", t.filename)))
	})
	t.templ.Execute(w, make(map[string]interface{}))
}

func RecallInstructions(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	t := templateHandler{filename: "recall_instructions.html"}
	t.once.Do(func() {
		t.templ = template.Must(template.ParseFiles(filepath.Join("templates", t.filename)))
	})
	t.templ.Execute(w, make(map[string]interface{}))
}

func QuestionHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	q := templateHandler{filename: "question.html"}
	q.once.Do(func() {
		q.templ = template.Must(template.ParseFiles(filepath.Join("templates", q.filename)))
	})

	index, err := strconv.ParseInt(ps.ByName("q"), 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("500 - Something went wrong. Sorry."))
		log.Fatal("Error converting string to int ", ps.ByName("q"), err)
	}
	question := res.questions[int(index)-1]

	data := map[string]interface{}{
		"Text":   question.text,
		"S1":     question.s1,
		"S2":     question.s2,
		"Number": index,
	}
	q.templ.Execute(w, data)
}

func RecallQuestion(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	q := templateHandler{filename: "recall_question.html"}
	q.once.Do(func() {
		q.templ = template.Must(template.ParseFiles(filepath.Join("templates", q.filename)))
	})

	index, err := strconv.ParseInt(ps.ByName("q"), 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("500 - Something went wrong. Sorry."))
		log.Fatal("Error converting string to int ", ps.ByName("q"), err)
	}
	text := recallTexts[int(index)-1]

	data := map[string]interface{}{
		"Text":   text,
		"Number": index,
	}
	q.templ.Execute(w, data)
}

func RecallStrength(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	t := templateHandler{filename: "recall_strength.html"}
	t.once.Do(func() {
		t.templ = template.Must(template.ParseFiles(filepath.Join("templates", t.filename)))
	})
	index, err := strconv.ParseInt(ps.ByName("q"), 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("500 - Something went wrong. Sorry."))
		log.Fatal("Error converting string to int ", ps.ByName("q"), err)
	}

	data := map[string]interface{}{
		"Number": index,
	}
	t.templ.Execute(w, data)
}

func DemographicsQuestion(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	qtype := ps.ByName("type")
	name := ps.ByName("q")

	var filename string
	if qtype == "mc" {
		filename = "mc_demographic.html"
	} else {
		filename = "text_demographics.html"
	}

	t := templateHandler{filename: filename}
	t.once.Do(func() {
		t.templ = template.Must(template.ParseFiles(filepath.Join("templates", t.filename)))
	})
	t.templ.Execute(w, demoData[name])
}
