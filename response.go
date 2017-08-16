package main

type Response struct {
	wave        int
	id          int
	condition   int
	targets     []Question
	distractors []Question
	questions   []Question
}
