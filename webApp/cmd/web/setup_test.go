package main

import (
	"flag"
	"log"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"webApp/pkg/db"
)

var app application

func TestMain(m *testing.M) {
	pathToTemplates = "./../../templates/"
	app.Session = getSession()

	root, _ := os.Getwd()
	for {
		if _, err := os.Stat(filepath.Join(root, "go.mod")); err == nil {
			break
		}
		root = filepath.Dir(root)
	}

	env, err := os.ReadFile(filepath.Join(root, ".env"))

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

	os.Exit(m.Run())
}
