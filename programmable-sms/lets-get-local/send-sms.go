package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

var accountSID, authToken string

func redirectPolicyFunc(req *http.Request, via []*http.Request) error {
	req.SetBasicAuth(accountSID, authToken)
	return nil
}

func main() {
	accountSID = os.Getenv("TWILIOQUEST_ACCOUNT_SID")
	authToken = os.Getenv("TWILIOQUEST_AUTH_TOKEN")
	toPhoneNum := os.Getenv("TWILIOQUEST_TO_PHONE_NUM")
	fromPhoneNum := os.Getenv("TWILIOQUEST_FROM_PHONE_NUM")
	smsSendURL := os.Getenv("TWILIOQUEST_SEND_SMS_ENDPOINT")
	smsBody := fmt.Sprintf("Greetings! The current time is: %s B2LO9VLWL2UT9XX", time.Now().Format("150405"))

	client := &http.Client{
		CheckRedirect: redirectPolicyFunc,
	}

	httpVerb := "POST"
	data := url.Values{}
	data.Set("To", toPhoneNum)
	data.Set("From", fromPhoneNum)
	data.Set("Body", smsBody)

	req, err := http.NewRequest(httpVerb, smsSendURL, strings.NewReader(data.Encode()))
	if err != nil {
		log.Fatal(err)
	}
	req.SetBasicAuth(accountSID, authToken)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	log.Printf("Executing %s %s\n", httpVerb, smsSendURL)
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	log.Printf("Response status: %s\n", resp.Status)
	log.Printf("Response body: %s", string(body))
}
