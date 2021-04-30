package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/smtp"
	"os"
)

type subscriberList struct {
	Subscribers []string `json:"subscribers"`
}

func main() {

	jsonFile, err := os.Open("emails.json")

	if err != nil {
		log.Print(err)
	}

	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)

	var list subscriberList

	json.Unmarshal(byteValue, &list)

	send("OwO", list.Subscribers)

    content := "Lorem Ipsum"
    addUnsubscribe(content string)
}

// Assuming the HTML rendered string is sent as contentString with the email
func addUnsubscribeLink(contentString string, email string) {

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
