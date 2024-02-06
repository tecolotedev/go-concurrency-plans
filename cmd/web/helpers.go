package main

import "fmt"

func (appConfig *Config) sendEmail(msg Message) {
	fmt.Println("here sendEmail 1")
	appConfig.Wait.Add(1)
	fmt.Println("here sendEmail 2")

	appConfig.Mailer.MailerChan <- msg
	fmt.Println("here sendEmail 3")

}
