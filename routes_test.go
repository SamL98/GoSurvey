package main

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/julienschmidt/httprouter"
)

func TestIntro(t *testing.T) {
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal("Error creating intro request ", err)
	}

	rr := httptest.NewRecorder()
	currWave = 2
	Intro(rr, req, httprouter.Params{})

	if status := rr.Code; status != http.StatusOK {
		t.Fatalf("Handler returned wrong status code: received %v, should be %v",
			status, http.StatusOK)
	}

	buf, err := ioutil.ReadFile("templates/intro.html")
	if err != nil {
		t.Fatalf("Could not read intro template file.")
	}
	expected := strings.TrimSpace(string(buf))
	received := strings.TrimSpace(rr.Body.String())

	if received != expected {
		expSlice := strings.Split(expected, "\n")
		recSlice := strings.Split(expected, "\n")

		for i := 0; i < len(expSlice); i++ {
			if expSlice[i] != recSlice[i] {
				t.Fatalf("Expected: %s, does not match received: %s", expSlice[i], recSlice[i])
				break
			}
		}
	}
}
