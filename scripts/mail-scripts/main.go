package main

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
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

// Assuming the HTML rendered string is sent as contentString with the email
func addUnsubscribeLink(contentString string, email string) string {

	// Get the encrypted hash to be sent
	unsubString := encryptUnsubscribeString(contentString, email)
	return unsubString
}

func encryptUnsubscribeString(plainSecret string, email string) string {

	encKey := os.Getenv("EMAIL_ENC_KEY")
	//Since the key is in string, we need to convert decode it to bytes
	key, err := hex.DecodeString(encKey)
	if err != nil {
		log.Println(err)
	}

	// convert the string to encrypt to bytes
	byteSecret := []byte(email)
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err.Error())
	}

	// Creating a new GCM
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		panic(err.Error())
	}

	// Creating a nonce. Nonce should be from GCM
	nonce := make([]byte, aesGCM.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		panic(err.Error())
	}

	ciphertext := aesGCM.Seal(nonce, nonce, byteSecret, nil)
	return fmt.Sprintf("%x", ciphertext)
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
