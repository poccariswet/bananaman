package main

import (
	"log"
)

type Theme struct{}

func (t *Theme) Run(args []string) int {
	if len(args) != 1 {
		log.Fatalf("Please set text '*********' ")
	}
	text := args[0]

	if err := GmailSend("テーマ", text); err != nil {
		log.Fatal(err)
	}
	return 0
}

func (c *Theme) Synopsis() string {
	return "This is 'Theme|テーマ' of Subject"
}

func (t *Theme) Help() string {
	return "subject is theme"
}
