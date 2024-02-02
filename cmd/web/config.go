package main

import (
	"log"
	"sync"

	"github.com/alexedwards/scs/v2"
	"github.com/jackc/pgx/v5"
)

type Config struct {
	Session  *scs.SessionManager
	Conn     *pgx.Conn
	InfoLog  *log.Logger
	ErrorLog *log.Logger
	Wait     *sync.WaitGroup
}
