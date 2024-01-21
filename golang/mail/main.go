package main

import (
	"fmt"
	"net/smtp"
	"os"
)

type Config struct {
	Password string
	UserName string
}

const SMTPServer = "smtp.gmail.com"

func Send(dest []string, bodyMessage string) error {
	userName := os.Getenv("SMTP_USERNAME")
	password := os.Getenv("SMTP_PASSWORD")

	var cfg Config
	if userName != "" || password != "" {
		cfg = Config{Password: password, UserName: userName}
	} else {
		return fmt.Errorf("SMTP_USERNAME and SMTP_PASSWORD environment variables are not set")
	}

	auth := smtp.PlainAuth("", cfg.UserName, cfg.Password, SMTPServer)

	msg := []byte("To: " + dest[0] + "\r\n" +
		"Subject: Newsletter\r\n" +
		"\r\n" +
		bodyMessage + "\r\n")

	err := smtp.SendMail(SMTPServer+":587", auth, cfg.UserName, dest, []byte(msg))
	if err != nil {
		return fmt.Errorf("error sending email: %v", err)
	}
	return nil
}

func main() {
	dest := []string{"perebaj@gmail.com"}
	bodyMessage := "Hello, this is a test email from golang"
	err := Send(dest, bodyMessage)
	if err != nil {
		panic(err)
	}
}
