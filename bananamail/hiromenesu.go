package main

import (
	"log"
)

type Hiromenesu struct{}

func (h *Hiromenesu) Run(args []string) int {
	if len(args) != 1 {
		log.Fatalf("Please set text '*********' ")
	}
	text := ParseMsg(args[0])

	if err := GmailSend("ヒロメネス", text); err != nil {
		log.Fatal(err)
	}
	return 0
}

func (h *Hiromenesu) Synopsis() string {
	return "This is 'hiromenesu|ヒロメネス' of Subject"
}

func (h *Hiromenesu) Help() string {
	return "subject is hiromenesu"
}
