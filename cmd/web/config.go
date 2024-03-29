package main

import (
	"final-project/cmd/web/data"
	"log"
	"sync"

	"github.com/alexedwards/scs/v2"
	"github.com/jackc/pgx/v5"
)

type Config struct {
	Session       *scs.SessionManager
	Conn          *pgx.Conn
	InfoLog       *log.Logger
	ErrorLog      *log.Logger
	Wait          *sync.WaitGroup
	Models        data.Models
	Mailer        Mail
	ErrorChan     chan error
	ErrorDoneChan chan bool
}
