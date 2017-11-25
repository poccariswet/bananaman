package main

import (
	"log"
)

type Sengen struct{}

func (s *Sengen) Run(args []string) int {
	if len(args) != 1 {
		log.Fatalf("Please set text '*********' ")
	}
	text := args[0]

	if err := GmailSend("宣言", text); err != nil {
		log.Fatal(err)
	}
	return 0
}

func (s *Sengen) Synopsis() string {
	return "This is 'sengen|宣言' of Subject"
}

func (s *Sengen) Help() string {
	return "subject is sengen"
}
