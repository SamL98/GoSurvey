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

func populateDemoData() {
	demoData = make(map[string]map[string]interface{})
	demoData["sex"] = map[string]interface{}{
		"Question": "What is your sex?",
		"Choices":  []string{"Male", "Female"},
		"NextType": "text",
		"NextName": "age",
		"Name":     "sex",
	}
	demoData["age"] = map[string]interface{}{
		"Question": "How old are you?",
		"NextType": "text",
		"NextName": "race",
		"Name":     "age",
	}
	demoData["race"] = map[string]interface{}{
		"Question": "What is your race?",
		"NextType": "text",
		"NextName": "major",
		"Name":     "race",
	}
	demoData["major"] = map[string]interface{}{
		"Question": "What is/was your major in college?",
		"NextType": "text",
		"NextName": "year",
		"Name":     "major",
	}
	demoData["year"] = map[string]interface{}{
		"Question": "What year are you in college?",
		"NextType": "mc",
		"NextName": "affiliation",
		"Name":     "year",
	}
	demoData["affiliation"] = map[string]interface{}{
		"Question": "Generally speaking,  do you usually think of yourself as a Republican, a Democrat, an Independent or what?",
		"Choices":  []string{"Strong Democrat", "Democrat", "Lean Democrat", "Lean Republican", "Republican", "Strong Republican"},
		"NextType": "mc",
		"NextName": "ideology",
		"Name":     "affiliation",
	}
	demoData["ideology"] = map[string]interface{}{
		"Question": "Generally speaking, do you consider yourself to be:",
		"Choices":  []string{"extemely conservative", "conservative", "slightly conservative", "moderate", "slightly liberal", "liberal", "extremely liberal"},
		"NextType": "mc",
		"NextName": "interest",
		"Name":     "ideology",
	}
	demoData["interest"] = map[string]interface{}{
		"Question": "How interested are you in what’s going on in American government and politics?",
		"Choices":  []string{"Very Interested", "Interested", "Somewhat Interested", "Not interested at all"},
		"NextType": "knowledge",
		"NextName": "majorities",
		"Name":     "interest",
	}
}

func populateKnowledgeData() {
	knowledgeData["majorities"] = map[string]interface{}{
		"Question": "How much of a majority is required for the U.S. Senate and House to override a presidential veto?",
		"Choices": []string{
			"3/4 in the House and 2/3 in the Senate",
			"2/3 in the House and 3/4 in the Senate",
			"3/4 in the House and 3/4 in the Senate",
			"2/3 in the House and 2/3 in the Senate",
		},
		"Name": "majorities",
		"Next": "majority",
	}
	knowledgeData["majority"] = map[string]interface{}{
		"Question": "Which party currently has the most members in the House of Representatives in Washington?",
		"Choices":  []string{"Democrats", "Republicans"},
		"Name":     "majority",
		"Next":     "war",
	}
	knowledgeData["war"] = map[string]interface{}{
		"Question": "Which branch of government has the official power to declare war?",
		"Choices":  []string{"President", "Congress", "Supreme Court"},
		"Name":     "war",
		"Next":     "pence",
	}
	knowledgeData["pence"] = map[string]interface{}{
		"Question": "What job or political office does Mike Pence now hold?",
		"Choices":  []string{"Secretary of State", "Vice President", "Speaker of the House", "Senate Minority Whip"},
		"Name":     "pence",
		"Next":     "majWhip",
	}
	knowledgeData["majWhip"] = map[string]interface{}{
		"Question": "Who is the current House Majority Whip?",
		"Choices":  []string{"Steve Scalise", "Roy Blunt", "Kevin McCarthy", "Jim Clyburn", "Eric Cantor"},
		"Name":     "majWhip",
		"Next":     "minWhip",
	}
	knowledgeData["minWhip"] = map[string]interface{}{
		"Question": "Who is the current House Minority Whipe?",
		"Choices":  []string{"Eric Cantor", "Roy Blunt", "Jim Clyburn", "Ted Stevens", "Steny Hoyer"},
		"Name":     "minWhip",
		"Next":     "pm",
	}
	knowledgeData["pm"] = map[string]interface{}{
		"Question": "Who is the current Prime Minister of Great Britain",
		"Choices":  []string{"Theresa May", "Gordon Brown", "Andrea Leadsom", "David Cameron"},
		"Name":     "pm",
		"Next":     "cons",
	}
	knowledgeDatap["cons"] = map[string]interface{}{
		"Question": "Which of the two political parties is more conservative?",
		"Choices":  []string{"Democrats", "Republicans"},
		"Name":     "cons",
		"Next":     "none",
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

func KnowledgeInstructions(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	t := templateHandler{filename: "knowledge_instructions.html"}
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
		"Number": int(index),
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

func KnwowledgeQuestion(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	q := templateHandler{filename: "knowledge_question.html"}
	q.once.Do(func() {
		q.templ = template.Must(template.ParseFiles(filepath.Join("templates", q.filename)))
	})

	name := ps.ByName("name")
	q.templ.Execute(w, knowledgeData[name])
}

func CompletionHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	t := templateHandler{filename: "complete.html"}
	t.once.Do(func() {
		t.templ = template.Must(template.ParseFiles(filepath.Join("templates", t.filename)))
	})
	t.templ.Execute(w, make(map[string]interface{}))
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

	log.Println(q.time)
	currRes.questions[int(index)-1] = q
}

func ParseStatistic(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	index, err := strconv.ParseInt(ps.ByName("q"), 10, 64)
	if err != nil {
		log.Fatal("Error parsing number for statistic", err)
	}
	no := int(index) - 1
	qno := no / 2
	sno := 0
	if no%2 != 0 {
		sno = 1
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatal("Error reading body", err)
	}
	defer r.Body.Close()
	bodyStr := string(body)

	q := currRes.questions[qno]
	if sno == 0 {
		q.s1 = bodyStr
	} else {
		q.s2 = bodyStr
	}
}

func ParseRecallStrength(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	index, err := strconv.ParseInt(ps.ByName("q"), 10, 64)
	if err != nil {
		log.Fatal("Error parsing number for recall strength", err)
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatal("Error reading body", err)
	}
	defer r.Body.Close()

	bodyStr := string(body)
	strength, err := strconv.ParseInt(bodyStr, 10, 64)
	if err != nil {
		log.Fatal("Error parsing strength", err)
	}

	qno := (int(index) - 1) / 2
	q := currRes.questions[qno]
	q.strength = int8(strength)
}

func ParseMCDemographic(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatal("Error reading body", err)
	}
	defer r.Body.Close()
	bodyStr := string(body)

	name := ps.ByName("name")
	switch name {
	case "sex":
		if bodyStr == "Male" {
			currRes.demographic.sex = false
		} else {
			currRes.demographic.sex = true
		}
	case "affiliation":
		if index := indexOf(bodyStr, demoData[name]["Choices"]); index >= 0 {
			currRes.demographic.affiliation = int8(index)
		}
	case "ideology":
		if index := indexOf(bodyStr, demoData[name]["Choices"]); index >= 0 {
			currRes.demographic.ideology = int8(index)
		}
	case "interest":
		if index := indexOf(bodyStr, demoData[name]["Choices"]); index >= 0 {
			currRes.demographic.interest = int8(index)
		}
	default:
		break
	}
}

func ParseTextDemographic(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatal("Error reading body", err)
	}
	defer r.Body.Close()
	bodyStr := string(body)

	name := ps.ByName("name")
	switch name {
	case "age":
		age, err := parseInt(bodyStr)
		if err != nil {
			log.Fatal("Error parsing age from body.")
		}
		currRes.demographic.age = int8(age)
	case "race":
		currRes.demographic.race = bodyStr
	case "major":
		currRes.demographic.major = bodyStr
	case "year":
		currRes.demographic.year = bodyStr
	default:
		break
	}
}

func ParseKnowledge(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatal("Error reading body", err)
	}
	defer r.Body.Close()
	bodyStr := string(body)

	name := ps.ByName("name")
	switch name {

	}
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
