package main

import (
	"log"
	"os"
	"strings"
	"testing"

	"webApp/pkg/repository/dbrepo"
)

var app application
var expiredToken string

func TestMain(m *testing.M) {
	env, err := os.ReadFile("../../.env")
	if err != nil {
		log.Fatal(err)
	}

	var jwtSecret string
	for line := range strings.SplitSeq(string(env), "\n") {
		key, value, ok := strings.Cut(strings.TrimSpace(line), "=")
		if !ok {
			continue
		}

		switch key {
		case "JWT_SECRET":
			jwtSecret = value
		case "EXPIRED_TOKEN":
			expiredToken = value
		}
	}
	if strings.TrimSpace(jwtSecret) == "" {
		log.Fatal("set JWT_SECRET in ../../.env")
	}

	if strings.TrimSpace(expiredToken) == "" {
		log.Fatal("set EXPIRED_TOKEN in ../../.env")
	}

	app = application{
		DB:        &dbrepo.TestDBRepo{},
		Domain:    "example.com",
		JWTSecret: jwtSecret,
	}

	m.Run()
}
