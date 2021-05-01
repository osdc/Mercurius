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

// Struct to store subscribers list array
type subscriberList struct {
	Subscribers []string `json:"subscribers"`
}

func main() {

	log.SetFlags(log.LstdFlags | log.Lshortfile)

	// Init goldmark configuration
	markdown := goldmark.New(
		goldmark.WithExtensions(
			meta.Meta,
		),
	)

	// Open emails file for reading emails (Will be replaced in future)
	emailFile, err := os.Open("emails.json")
	if err != nil {
		log.Print(err)
	}
	defer emailFile.Close()

	// Read emails from the file into a buffer
	byteValue, _ := ioutil.ReadAll(emailFile)

	// Transferring email data to a var
	var subsList subscriberList
	json.Unmarshal(byteValue, &subsList)

	// Reading content to be sent
	// TODO: Replace this with logic to read file containing the latest post
	content, _ := ioutil.ReadFile("../../content/post/example.md")

	// Convert the markdown into a compatible HTML format
	var htmlContent bytes.Buffer
	if err := markdown.Convert(content, &htmlContent); err != nil {
		log.Fatalln(err)
	}
	// Convert to string for using and tweaking it everywhere
	contentString := htmlContent.String()

	// Iterate over emails list and start sending
	for _, email := range subsList.Subscribers {

		// Get the unsubscribe hash string
		completeContent := addUnsubscribeLink(contentString, email)

		// TODO here: Add image to the top of content if needed

		// Send the email to the recipient
		err := send(completeContent, email)
		if err != nil {
			log.Println("Error while sending mail to subscriber", email, "\nError : ", err)
		}
	}

}

// Utility to get the unsubscribe hash
func addUnsubscribeLink(contentString string, email string) string {

	// Get the encrypted hash to be sent
	unsubString := encryptUnsubscribeString(email)

	// Prepare the HTML template of the unsubscribe option (raw as of now)
	unSubscribeTemplate := fmt.Sprintf("<a href=\"/unsubscribe?uniqhash=%s\">UnSubscribe</a>", unsubString)

	// Append the hash to the email body
	newContent := fmt.Sprintf("%s\n\n%s", contentString, unSubscribeTemplate)

	return newContent
}

// Helper function to above to encrypt user email in order to find
// unsubscribe link hash
func encryptUnsubscribeString(email string) string {

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

// Utility to send emails using gomail
func send(body string, to string) error {
	from := os.Getenv("MAIL_ID")
	pass := os.Getenv("MAIL_PASSWORD")

	dialer := gomail.NewDialer("smtp.gmail.com", 587, from, pass)
	s, err := dialer.Dial()
	if err != nil {
		log.Println(err)
		return err
	}

	m := gomail.NewMessage()
	m.SetHeader("From", from)
	m.SetAddressHeader("To", to, to)
	m.SetHeader("Subject", "Newsletter Test")
	m.SetBody("text/html", body)

	if err := gomail.Send(s, m); err != nil {
		log.Printf("Could not send email to %q: %v", body, err)
		return err
	}
	m.Reset()
	return nil
}
