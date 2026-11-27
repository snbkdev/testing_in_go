package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type application struct {
	JWTSecret string
	Action    string
}

// This is used to generate a token, so that we can test our api. Run this with go run ./cmd/cli and copy
// the token that is printed out.
// go run ./cmd/cli -action=valid     // will produce a valid token
// go run ./cmd/cli -action=expired   // will produce an expired token

func main() {
	var app application
	flag.StringVar(&app.JWTSecret, "jwt-secret", "", "signing secret (defaults to JWT_SECRET in .env)")
	flag.StringVar(&app.Action, "action", "valid", "action: valid|expired")
	flag.Parse()

	if app.JWTSecret == "" {
		env, err := os.ReadFile(".env")
		if err != nil && !os.IsNotExist(err) {
			log.Fatal(err)
		}
		for line := range strings.SplitSeq(string(env), "\n") {
			key, value, ok := strings.Cut(strings.TrimSpace(line), "=")
			if ok && key == "JWT_SECRET" {
				app.JWTSecret = value
				break
			}
		}
	}
	if strings.TrimSpace(app.JWTSecret) == "" {
		log.Fatal("set JWT_SECRET in .env or pass -jwt-secret")
	}

	// generate a token
	token := jwt.New(jwt.SigningMethodHS256)

	// set claims
	claims := token.Claims.(jwt.MapClaims)
	claims["name"] = "John Doe"
	claims["sub"] = "1"
	claims["admin"] = true
	claims["aud"] = "example.com"
	claims["iss"] = "example.com"
	// leave this to 3 days, for easy manual testing
	if app.Action == "valid" {
		expires := time.Now().UTC().Add(time.Hour * 72)
		claims["exp"] = expires.Unix()
	} else {
		expires := time.Now().UTC().Add(time.Hour * 100 * -1)
		claims["exp"] = expires.Unix()
	}

	// create the token as a slice of bytes
	if app.Action == "valid" {
		fmt.Println("VALID Token:")
	} else {
		fmt.Println("EXPIRED Token:")
	}
	signedAccessToken, err := token.SignedString([]byte(app.JWTSecret))
	if err != nil {
		log.Fatal(err)
	}
	// print to console
	fmt.Println(string(signedAccessToken))
}
