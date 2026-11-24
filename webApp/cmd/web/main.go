package main

import (
	"encoding/gob"
	"flag"
	"log"
	"net/http"
	"os"
	"strings"
	"webApp/pkg/data"
	"webApp/pkg/db"

	"github.com/alexedwards/scs/v2"
)

type application struct {
	DSN     string
	DB      db.PostgresConn
	Session *scs.SessionManager
}

func main() {
	gob.Register(data.User{})

	// set up app config
	app := application{}

	env, err := os.ReadFile(".env")
	if err != nil && !os.IsNotExist(err) {
		log.Fatal(err)
	}

	var dsn string
	for line := range strings.SplitSeq(string(env), "\n") {
		if value, ok := strings.CutPrefix(strings.TrimSpace(line), "DB_DSN="); ok {
			dsn = value
			break
		}
	}

	flag.StringVar(&app.DSN, "dsn", dsn, "Postgres connection")
	flag.Parse()
	if app.DSN == "" {
		log.Fatal("set DB_DSN in .env or pass -dsn")
	}

	conn, err := app.connectToDB()
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	app.DB = db.PostgresConn{DB: conn}

	// get a session manager
	app.Session = getSession()

	// print out a message
	log.Println("Starting server on port 8082...")

	// start the server
	err = http.ListenAndServe(":8082", app.routes())
	if err != nil {
		log.Fatal(err)
	}
}
