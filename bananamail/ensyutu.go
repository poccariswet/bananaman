package main

import (
	"log"
)

type Ensyutu struct{}

func (e *Ensyutu) Run(args []string) int {
	if len(args) != 1 {
		log.Fatalf("Please set text '*********' ")
	}
	text := args[0]

	if err := GmailSend("演出", text); err != nil {
		log.Fatal(err)
	}
	return 0
}

func (e *Ensyutu) Synopsis() string {
	return "This is 'ensyutu|演出' of Subject"
}

func (e *Ensyutu) Help() string {
	return "subject is ensyutu"
}
