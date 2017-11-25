package main

import (
	"log"
)

type Henken struct{}

func (h *Henken) Run(args []string) int {
	if len(args) != 1 {
		log.Fatalf("Please set text '*********' ")
	}
	text := args[0]

	if err := GmailSend("偏見", text); err != nil {
		log.Fatal(err)
	}
	return 0
}

func (h *Henken) Synopsis() string {
	return "This is 'henken|偏見' of Subject"
}

func (h *Henken) Help() string {
	return "subject is henken"
}
