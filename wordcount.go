package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

//Email represents the sender, recipients, and body of an email.
type Email struct {
	To   []string
	From string
	Body string
}

func main() {
	start := time.Now()
	searchDir := "./maildir/"
	mailFileChan := make(chan string)
	emailChan := make(chan Email)
	answerChan := make(chan int)
	wg := sync.WaitGroup{}

	go wordCount(emailChan, answerChan)

	for i := 0; i < 16; i++ {
		go handleMailFile(mailFileChan, emailChan, &wg)
	}
	err := filepath.Walk(searchDir, func(path string, f os.FileInfo, err error) error {
		wg.Add(1)
		mailFileChan <- path
		return nil
	})
	if err != nil {
		fmt.Print(err.Error())
	}
	wg.Wait()
	close(emailChan)
	words := <-answerChan
	log.Printf("Counted %d words in %v", words, time.Since(start))
}

func wordCount(emailChan chan Email, answerChan chan int) {
	accumulator := 0
	for {
		if email, open := <-emailChan; open {
			bodywords := strings.Fields(email.Body)
			accumulator += len(bodywords)
		} else {
			break
		}
	}
	answerChan <- accumulator
}

func handleMailFile(mailFileChan chan string, emailChan chan Email, wg *sync.WaitGroup) {
	for {
		path, open := <-mailFileChan
		if !open {
			break
		}
		fi, err := os.Lstat(path)
		if err == nil && !fi.IsDir() {
			rawEmail, rawError := readEmailContents(path)
			if rawError != nil {
				log.Printf("Error reading email: %+v", rawError)
				wg.Done()
				continue
			}
			email := ParseEmail(rawEmail)
			emailChan <- email
		}
		wg.Done()
	}
}

func readEmailContents(path string) (string, error) {
	bytestring, err := ioutil.ReadFile(path)
	if err != nil {
		log.Printf("Error reading %s: %+v", path, err)
		return "", err
	}

	return string(bytestring), nil
}

//ParseEmail parses the raw text of the email and finds the To, From,
// and unique body text of that email.
func ParseEmail(rawEmail string) Email {
	emailLines := strings.Split(rawEmail, "\n")
	toReturn := Email{}

	headerEndIndex := 0
	quotedTextBeginIndex := len(emailLines) - 1

	for i := 0; i < len(emailLines); i++ {
		if strings.HasPrefix(emailLines[i], "From: ") && toReturn.From == "" {
			toReturn.From = strings.TrimPrefix(emailLines[i], "From: ")
		}

		if strings.HasPrefix(emailLines[i], "To: ") && toReturn.To == nil {
			toReturn.To = strings.Split(strings.TrimPrefix(emailLines[i], "To: "), ",")
		}

		//The first blank line will be the line that divides the header and the body
		if emailLines[i] == "" && headerEndIndex == 0 {
			headerEndIndex = i
		}

		//The first line that starts with "-----" will *probably* be the end of the body text
		if strings.HasPrefix(emailLines[i], "-----") && quotedTextBeginIndex == len(emailLines)-1 {
			quotedTextBeginIndex = i
		}
	}

	if quotedTextBeginIndex < headerEndIndex+1 {
		quotedTextBeginIndex = headerEndIndex
	}

	toReturn.Body = strings.Join(emailLines[headerEndIndex+1:quotedTextBeginIndex+1], "\n")

	return toReturn
}
