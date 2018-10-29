package main

import (
	"errors"
	"fmt"
	"net/http"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-lambda-go/lambda/events"
	sendgrid "github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

func main() {
	lambda.Start(Handler)
}

func Handler(req events.APIGatewayProxyRequest) (m events.APIGatewayProxyResponse, err error) {
	name := req.QueryStringParameters["name"]
	if name == "" {
		return clientError(http.StatusBadRequest)
	}
	to := req.QueryStringParameters["to"]
	if to == "" {
		return clientError(http.StatusBadRequest)
	}
	subject := req.QueryStringParameters["subject"]
	if to == "" {
		return clientError(http.StatusBadRequest)
	}
	msg := req.QueryStringParameters["message"]
	if to == "" {
		return clientError(http.StatusBadRequest)
	}
	err = SendMail(name, to, subject, msg)
	if err != nil {
		return clientError(http.StatusBadRequest)
	}
	return events.APIGatewayProxyResponse{StatusCode: http.StatusOK}, err
}

func SendMail(userName, userEmail, subject, msg string) (err error) {
	from := mail.NewEmail("sendgrid example", "example@sendgrid.com")
	to := mail.NewEmail(userName, userEmail)
	message := mail.NewSingleEmail(from, subject, to, msg, msg)
	client := sendgrid.NewSendClient(os.Getenv("SENDGRID_API_KEY"))
	response, err := client.Send(message)
	if err != nil {
		return
	}
	if response.StatusCode != http.StatusOK {
		return errors.New(
			fmt.Sprintf("Email send failure, status code: %v | body: %s", response.StatusCode, response.Body)
		)
	}
	return
}

func clientError(status int) (events.APIGatewayProxyResponse, error) {
	return events.APIGatewayProxyResponse{
		StatusCode: status,
		Body:       http.StatusText(status),
	}, nil
}
