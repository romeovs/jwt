package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/apex/log"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/mgutz/ansi"
	"github.com/spf13/cobra"
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

func fromStdin(args []string) (string, bool) {
	stat, _ := os.Stdin.Stat()
	if (stat.Mode() & os.ModeCharDevice) == 0 {
		if len(args) != 0 {
			log.Fatal("Got token both as argument and stdin")
		}

		log.Debug("Got token on stdin")

		read, err := ioutil.ReadAll(os.Stdin)
		if err != nil {
			log.WithError(err).Fatal("Could not read from std in")
		}
		return string(read), true
	}

	return "", false
}

func fromFile(cmd *cobra.Command, args []string) (string, bool) {
	if len(args) != 1 {
		cmd.UsageFunc()(cmd)
		return "", false
	}

	filename := args[0]

	if _, err := os.Stat(filename); err == nil {
		log.WithField("file", filename).Debug("Got filename")
		content, err := ioutil.ReadFile(filename)
		if err != nil {
			log.WithField("file", filename).WithError(err).Fatal("Could not read file")
		}
		return string(content), true
	} else if forceFile {
		msg := err.Error()
		parts := strings.Split(err.Error(), ": ")
		if len(parts) > 1 {
			msg = parts[1]
		}
		log.WithField("error", msg).WithField("file", filename).Fatal("Could not read file")
	}

	return "", false
}

func getToken(cmd *cobra.Command, args []string) string {

	token, ok := fromStdin(args)
	if ok {
		return token
	}

	if len(args) != 1 {
		cmd.UsageFunc()(cmd)
		os.Exit(1)
	}

	if forceInput {
		return args[0]
	}

	token, ok = fromFile(cmd, args)
	if ok {
		return token
	}

	log.Debug("Got raw token")
	return args[0]
}

func tryJSON(token string) string {
	if forceRaw {
		return token
	}
	var s Token
	err := json.Unmarshal([]byte(token), &s)
	if err == nil {
		log.Debug("Got OAuth 2 JSON")
		return s.AccessToken
	} else if forceOAuth {
		log.WithError(err).Fatal("Could not parse OAuth 2 response")
	}
	return token
}
