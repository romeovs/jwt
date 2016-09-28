package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/apex/log"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/mgutz/ansi"
)

func unmarshal(segment string) jwt.MapClaims {
	decoded, err := jwt.DecodeSegment(segment)
	if err != nil {
		log.WithError(err).Fatal("Could not decode token")
	}

	var claims jwt.MapClaims
	err = json.Unmarshal(decoded, &claims)
	if err != nil {
		log.WithError(err).Fatal("Could not unmarshal JSON in token")
	}

	return claims
}

func field(name string, value interface{}) {
	key := fmt.Sprintf("%10s", name)
	if !noColor {
		key = ansi.Color(key, "blue")
	}
	fmt.Printf("%s  %s\n", key, value)
}
