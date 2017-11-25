package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/smtp"
	"os"
	"path/filepath"

	"github.com/mitchellh/cli"
)

type Mail struct {
	From  string `json:"from"`
	Pass  string `json:"pass"`
	To    string `json:"to"`
	Rname string `json:"rname"`
}

var (
	root     string
	fileroot string
)

func init() {
	homepath := os.Getenv("HOME")
	root = filepath.Join(homepath, ".gmail")
	_, err := os.Stat(root)
	if err != nil {
		os.Mkdir(root, 0777)
		fmt.Println("Made a Dir at" + root)
	}
	fileroot = filepath.Join(root, "gmail.json")
}

func GmailSend(sub, text string) error {
	m := Mail{}
	file, err := ioutil.ReadFile(fileroot)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(file, &m); err != nil {
		return err
	}

	auth := smtp.PlainAuth(
		"",
		m.From,
		m.Pass,
		"smtp.gmail.com",
	)
	msg := "From:" + m.From + "\n" +
		"To:" + m.To + "\n" +
		"Subject: " + sub + "\n\n" +
		"ラジオネーム" + "  " + m.Rname + "\n\n" + text

	err = smtp.SendMail("smtp.gmail.com:587", auth, m.From, []string{m.To}, []byte(msg))
	if err != nil {
		return err
	}
	return nil
}

func main() {
	c := cli.NewCLI("Gmail to bananamoon cli", "1.0.0")
	c.Args = os.Args[1:]
	c.Commands = map[string]cli.CommandFactory{
		"init": func() (cli.Command, error) {
			return &Init{}, nil
		},
		"theme": func() (cli.Command, error) {
			return &Theme{}, nil
		},
		"hiromenesu": func() (cli.Command, error) {
			return &Hiromenesu{}, nil
		},
		"henken": func() (cli.Command, error) {
			return &Henken{}, nil
		},
		"sengen": func() (cli.Command, error) {
			return &Sengen{}, nil
		},
		"ensyutu": func() (cli.Command, error) {
			return &Ensyutu{}, nil
		},
	}

	exitStatus, err := c.Run()
	if err != nil {
		log.Println(err)
	}

	os.Exit(exitStatus)
}
