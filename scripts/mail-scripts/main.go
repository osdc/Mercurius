package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html"
	"html/template"
	"io/ioutil"
	"log"
	"os"

	"github.com/yuin/goldmark"
	meta "github.com/yuin/goldmark-meta"
	"github.com/yuin/goldmark/parser"
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
	context := parser.NewContext()
	if err := markdown.Convert(content, &buf, parser.WithContext(context)); err != nil {
		panic(err)
	}
	metaData := meta.Get(context)
	title := metaData["title"]
	str := fmt.Sprintf("%v", title)

	t := template.New("template.html")
	t, _ = t.ParseFiles("template.html")

	var body bytes.Buffer

	if err := t.Execute(&body, struct {
		Content string
		Title   string
	}{
		Content: buf.String(),
		Title:   str,
	}); err != nil {
		log.Println(err)
	}
	html := html.UnescapeString(body.String())

	send(html, list.Subscribers)
}

func send(body string, to []string) {
	from := os.Getenv("MAIL_ID")
	pass := os.Getenv("MAIL_PASSWORD")

	d := gomail.NewDialer("smtp.gmail.com", 587, from, pass)
	s, err := d.Dial()
	if err != nil {
		panic(err)
	}

	bodyContent, err := ioutil.ReadFile("email_body.html")
    if err != nil {
        log.Fatal(err)
    }

	m := gomail.NewMessage()
	for _, r := range to {
		fmt.Printf("Sending email to: %s\n", r)
		m.SetHeader("From", from)
		m.SetAddressHeader("To", r, r)
		m.SetHeader("Subject", "Mercurius Test")
		m.SetBody("text/html", string(bodyContent))

		if err := gomail.Send(s, m); err != nil {
			log.Printf("Could not send email to %q: %v", r, err)
		}
		m.Reset()
	}
}
