package main

import (
	"reflect"
	"testing"
)

func TestParseEmail(t *testing.T) {
	rawTestEmail := `Date: Fri, 15 Dec 2000 08:20:00 -0800 (PST)
From: michelle.cash@enron.com
To: molly.magee@enron.com
Subject: Andy Becnel
Mime-Version: 1.0
Content-Type: text/plain; charset=us-ascii
Content-Transfer-Encoding: 7bit
X-From: Michelle Cash
X-To: Molly Magee
X-cc:
X-bcc:
X-Folder: \Michelle_Cash_Nov2001\Notes Folders\Sent
X-Origin: Cash-M
X-FileName: mcash_NonPriv.nsf

Molly,

Let's schedule interviews with Bechel with Mark Haedicke, Elizabeth Sager,
Sheila Tweed, Dan Lyons, and Lisa Mellencamp.

Thanks a lot.

Michelle`

	expectedEmail := Email{
		To:   []string{"molly.magee@enron.com"},
		From: "michelle.cash@enron.com",
		Body: `Molly,

Let's schedule interviews with Bechel with Mark Haedicke, Elizabeth Sager,
Sheila Tweed, Dan Lyons, and Lisa Mellencamp.

Thanks a lot.

Michelle`,
	}

	result := ParseEmail(rawTestEmail)

	if expectedEmail.From != result.From {
		t.Errorf("Expected From %v, got %v.", expectedEmail.From, result.From)
	}

	if !reflect.DeepEqual(expectedEmail.To, result.To) {
		t.Errorf("Expected To %v, got %v", expectedEmail.To, result.To)
	}

	if expectedEmail.Body != result.Body {
		t.Errorf("Expected Body %v, got %v", expectedEmail.Body, result.Body)
	}
}

func TestWordcount(t *testing.T) {

	testEmail := Email{
		To:   []string{"molly.magee@enron.com"},
		From: "michelle.cash@enron.com",
		Body: `Molly,

Let's schedule interviews with Bechel with Mark Haedicke, Elizabeth Sager,
Sheila Tweed, Dan Lyons, and Lisa Mellencamp.

Thanks a lot.

Michelle`,
	}

	testEmailChan := make(chan Email)
	testAnswerChan := make(chan int)

	go wordCount(testEmailChan, testAnswerChan)

	testEmailChan <- testEmail

	close(testEmailChan)
	answer := <-testAnswerChan

	if answer != 22 {
		t.Errorf("Expected 22 words, got %v", answer)
	}

}
