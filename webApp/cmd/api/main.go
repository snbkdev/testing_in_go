package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"webApp/pkg/repository"
	"webApp/pkg/repository/dbrepo"
)

const port = 8090

type application struct {
	DSN       string
	DB        repository.DatabaseRepo
	Domain    string
	JWTSecret string
}

func main() {
	var app application
	flag.StringVar(&app.Domain, "domain", "example.com", "Domain for application, eg company.com")

	env, err := os.ReadFile(".env")
	if err != nil && !os.IsNotExist(err) {
		log.Fatal(err)
	}

	values := map[string]string{}
	for line := range strings.SplitSeq(string(env), "\n") {
		key, value, ok := strings.Cut(strings.TrimSpace(line), "=")
		if !ok {
			continue
		}
		if _, exists := values[key]; !exists {
			values[key] = value
		}
	}

	flag.StringVar(&app.DSN, "dsn", values["DB_DSN"], "Postgres connection")
	flag.StringVar(&app.JWTSecret, "jwt-secret", "", "signing secret")
	flag.Parse()
	if app.DSN == "" {
		log.Fatal("set DB_DSN in .env or pass -dsn")
	}

	if app.JWTSecret == "" {
		app.JWTSecret = values["JWT_SECRET"]
	}
	if strings.TrimSpace(app.JWTSecret) == "" {
		log.Fatal("set JWT_SECRET in .env or pass -jwt-secret")
	}

	flag.Parse()

	conn, err := app.connectToDB()
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	app.DB = &dbrepo.PostgresDBRepo{DB: conn}

	log.Printf("Starting api on port %d\n", port)

	err = http.ListenAndServe(fmt.Sprintf(":%d", port), app.routes())
	if err != nil {
		log.Fatal(err)
	}
}
