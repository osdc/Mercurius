package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"os"

	"github.com/yuin/goldmark"
	meta "github.com/yuin/goldmark-meta"
	"gopkg.in/gomail.v2"
)

type subscriberList struct {
	Subscribers []string `json:"subscribers"`
}

func main() {
	markdown := goldmark.New(
		goldmark.WithExtensions(
			meta.Meta,
		),
	)

	jsonFile, err := os.Open("emails.json")

	if err != nil {
		log.Print(err)
	}

	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)

	var list subscriberList

	json.Unmarshal(byteValue, &list)

	content, _ := ioutil.ReadFile("../../content/post/example.md")

	var buf bytes.Buffer

	if err := markdown.Convert(content, &buf); err != nil {
		panic(err)
	}

	log.Print(string(buf.String()))

	send(string(buf.String()), list.Subscribers)
}

func send(body string, to []string) {
	from := os.Getenv("MAIL_ID")
	pass := os.Getenv("MAIL_PASSWORD")

	d := gomail.NewDialer("smtp.gmail.com", 587, from, pass)
	s, err := d.Dial()
	if err != nil {
		panic(err)
	}

	m := gomail.NewMessage()
	for _, r := range to {
		m.SetHeader("From", from)
		m.SetAddressHeader("To", r, r)
		m.SetHeader("Subject", "Newsletter Test")
		m.SetBody("text/html", body)

		if err := gomail.Send(s, m); err != nil {
			log.Printf("Could not send email to %q: %v", r, err)
		}
		m.Reset()
	}
}
