package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

type Init struct{}

func (i *Init) Run(args []string) int {
	if len(os.Args) != 5 {
		log.Fatalf("Please set the prams **your address** **pass** **radio name**")
	}

	var mail, passwd, radioname string
	mail = args[0]
	passwd = args[1]
	radioname = args[2]

	m := Mail{
		From:  mail,
		Pass:  passwd,
		To:    "rainmaker0027pockets@ezweb.ne.jp",
		Rname: radioname,
	}
	bytes, _ := json.Marshal(m)
	ioutil.WriteFile(filepath.Join(root, "gmail.json"), bytes, os.ModePerm)

	return 0
}

func (i *Init) Synopsis() string {
	return "This is able to set your address, gmail app pass, your radio name"
}

func (i *Init) Help() string {
	return "init set address, pass, name"
}
