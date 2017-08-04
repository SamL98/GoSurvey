package main

type Response struct {
	wave        int
	id          int
	ip          string
	start       int
	questions   []Question
	demographic Demographics
	knowledge   Knowledge
}
