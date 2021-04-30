package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/smtp"
	"os"

	"github.com/yuin/goldmark"
	meta "github.com/yuin/goldmark-meta"
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

	msg := "Subject: Hello there o/ \n\n" +
		body

	err := smtp.SendMail("smtp.gmail.com:587",
		smtp.PlainAuth("", from, pass, "smtp.gmail.com"),
		from, to, []byte(msg))

	if err != nil {
		log.Printf("smtp error: %s", err)
		return
	}

	log.Print("sent")
}
