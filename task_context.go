package main

// TaskContext holds common values useful for task execution
type TaskContext struct {
	startPageURL   string
	startPageTitle string
	pageTitle      string
	pageURL        string
	linkURL        string
	linkText       string
	imgAlt         string
	extension      string

	counter *counter
}
