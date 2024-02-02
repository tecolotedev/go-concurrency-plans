package main

import (
	"context"
	"final-project/cmd/web/db"
	"final-project/cmd/web/session"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"
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
	}

	appConfig.serve()

	// set up mail

	// listen for web connections

}
