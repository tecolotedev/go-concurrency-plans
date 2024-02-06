package main

import (
	"context"
	"final-project/cmd/web/data"
	"final-project/cmd/web/db"
	"final-project/cmd/web/session"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

const webPort = "80"

func (appConfig *Config) serve() {
	// start http server
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", webPort),
		Handler: appConfig.routes(),
	}

	appConfig.InfoLog.Println("Starting web server!")
	err := srv.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}
}

func main() {
	// connect to the database
	conn := db.InitDB()
	conn.Ping(context.Background())

	// create sessions
	usersSession := session.InitSession()

	// create loggers
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	// create channels

	// create waitgroups
	wg := sync.WaitGroup{}

	// set up the application config
	appConfig := Config{
		Session:  usersSession,
		Conn:     conn,
		Wait:     &wg,
		InfoLog:  infoLog,
		ErrorLog: errorLog,
		Models:   data.New(conn),
	}

	// set up mail
	appConfig.Mailer = appConfig.createMail()
	go appConfig.listenForMail()

	// listen for signals
	go appConfig.listenForShutdown()

	// listen for web connections

	appConfig.serve()

}

func (appConfig *Config) listenForShutdown() {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	appConfig.shutdown()
	os.Exit(0)
}
func (appConfig *Config) shutdown() {
	// perform any cleanup tasks
	appConfig.InfoLog.Println("Would run cleanup tasks...")

	// block until wg is empty
	appConfig.Wait.Wait()
	appConfig.Mailer.DoneChan <- true // wait till we send all emails

	appConfig.InfoLog.Println("closing channels and shutting down applications...")
	close(appConfig.Mailer.MailerChan)
	close(appConfig.Mailer.ErrorChan)
	close(appConfig.Mailer.DoneChan)
}

func (appConfig *Config) createMail() Mail {
	// create channels
	errorChan := make(chan error)
	mailChan := make(chan Message)
	mailerDoneChan := make(chan bool)

	m := Mail{
		Domain:      "localhost",
		Host:        "localhost",
		Port:        1025,
		Encryption:  "none",
		FromName:    "Info",
		FromAddress: "info@mycompany.com",
		Wait:        appConfig.Wait,
		ErrorChan:   errorChan,
		MailerChan:  mailChan,
		DoneChan:    mailerDoneChan,
	}
	return m

}
