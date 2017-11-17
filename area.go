package main

import (
	"net/http"

	"golang.org/x/net/html"
)

const (
	areaURL = "http://radiko.jp/area"
)

func GetAreaID() (string, error) {
	var areaID string
	resp, err := http.Get(areaURL)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	doc, err := html.Parse(resp.Body)
	if err != nil {
		return "", err
	}

	var search func(*html.Node)
	search = func(doc *html.Node) {
		if doc.Type == html.ElementNode && doc.Data == "span" && len(doc.Attr) > 0 {
			areaID = doc.Attr[0].Val
		}
		for fchild := doc.FirstChild; fchild != nil; fchild = fchild.NextSibling {
			search(fchild)
		}
	}
	search(doc)
	return areaID, nil
}
