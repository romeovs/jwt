package cmd

import (
	"encoding/json"

	"github.com/apex/log"
	jwt "github.com/dgrijalva/jwt-go"
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
