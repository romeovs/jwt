package cmd

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/apex/log"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/spf13/cobra"
)

var RootCmd = &cobra.Command{
	Use:   "jwt [token]",
	Short: "jwt can be used the debug JWT tokens.",
	Long:  "A simple jwt debugging tool written in Go.",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			cmd.UsageFunc()
			return
		}

		token := args[0]
		parts := strings.Split(token, ".")

		if len(parts) != 3 {
			log.Fatal("Token has invalid number of segments")
		}

		segment, err := jwt.DecodeSegment(parts[1])
		if err != nil {
			log.WithError(err).Fatal("Could not decode token")
		}

		var claims jwt.MapClaims
		err = json.Unmarshal(segment, &claims)
		if err != nil {
			log.WithError(err).Fatal("Could not unmarshal JSON in token")
		}

		indented, err := json.MarshalIndent(claims, "", "  ")
		if err != nil {
			log.WithError(err).Fatal("Could not indent JSON")
		}

		fmt.Println(string(indented))
	},
}
