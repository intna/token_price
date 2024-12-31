package actions

import (
	"fmt"
	"log"

	"token_price/src/utils"

	"github.com/resend/resend-go/v2"
)

func Send(price float64, token string) {
	emailContent := fmt.Sprintf(
		"[%s] The latest price is <span style=\"color:red;\">%.9f</span>",
		token, price)
	subject := fmt.Sprintf("[%s]Current Price Update", token)
	apiKey := utils.GetValue("RESEND_KEY", "")
	client := resend.NewClient(apiKey)

	gmail := utils.GetValue("GMAIL", "")
	params := &resend.SendEmailRequest{
		From:    "Cris <tokenPrice@resend.dev>",
		To:      []string{gmail},
		Html:    emailContent,
		Subject: subject,
	}

	_, err := client.Emails.Send(params)
	if err != nil {
		log.Println(err.Error())
		return
	}
	log.Printf("Success send email to %s", params.To[0])
}
