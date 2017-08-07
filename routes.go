package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"text/template"

	"github.com/julienschmidt/httprouter"
)

type templateHandler struct {
	once     sync.Once
	filename string
	templ    *template.Template
}

func Intro(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	t := templateHandler{filename: "intro.html"}
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

func CompletionHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	t := templateHandler{filename: "complete.html"}
	t.once.Do(func() {
		t.templ = template.Must(template.ParseFiles(filepath.Join("templates", t.filename)))
	})
	if err := pgManager.AddResponse(&currRes, res.id); err != nil {
		log.Println("Error adding response to postgres", err)
	}
	data := map[string]interface{}{
		"Condition": res.condition,
	}
	t.templ.Execute(w, data)
}

func ParseInterest(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	index, err := strconv.ParseInt(ps.ByName("q"), 10, 64)
	if err != nil {
		log.Fatal("Error parsing question number from params", err)
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatal("Error reading body", err)
	}
	defer r.Body.Close()
	bodyStr := string(body)

	q := Question{number: int8(index)}
	items := strings.Split(bodyStr, ",")
	for i := range items {
		item := items[i]
		splits := strings.Split(item, ":")
		if splits[0] == "time" {
			if timeSpent, err := strconv.ParseInt(splits[1], 10, 64); err != nil {
				q.time = int(timeSpent)
			}
		} else if splits[0] == "interest" {
			if interest, err := strconv.ParseInt(splits[1], 10, 64); err != nil {
				if int(interest) == 0 {
					q.interest = false
				} else {
					q.interest = true
				}
			}
		}
	}
	currRes.questions[int(index)-1] = q
}

func parseInt(str string) (result int, err error) {
	integer, err := strconv.ParseInt(str, 10, 64)
	return int(integer), err
}

func indexOf(str string, strSlice []string) int {
	for i := 0; i < len(strSlice); i++ {
		if strSlice[i] == str {
			return i
		}
	}
	return -1
}
