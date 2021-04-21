package cmd

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/apex/log"
	"github.com/romeovs/jwt/util"
	"github.com/spf13/cobra"
)

var DecodeCmd = &cobra.Command{
	Use:   "decode",
	Short: "Decode the JWT token and show its claims",
	Run: func(cmd *cobra.Command, args []string) {

		token := getToken(cmd, args)
		token = tryJSON(token)

		parts := strings.Split(token, ".")

		if len(parts) != 3 {
			log.Fatal("Token has invalid number of segments")
		}

		header := unmarshal(parts[0])
		claims := unmarshal(parts[1])

		indented, err := json.MarshalIndent(claims, "", "  ")
		if err != nil {
			log.WithError(err).Fatal("Could not indent JSON")
		}

		now := time.Now().Round(time.Second)

		validity := "token is valid"

		var issued *time.Time
		if claims["iat"] != nil {
			iat := time.Unix(int64(claims["iat"].(float64)), 0)
			if iat.After(now) {
				validity = fmt.Sprintf("token is not valid for %s", iat.Sub(now))
			}
			issued = &iat
		}

		var expires *time.Time
		if claims["exp"] != nil {
			exp := time.Unix(int64(claims["exp"].(float64)), 0)
			if now.After(exp) {
				validity = fmt.Sprintf("token is expired for %s", now.Sub(exp))
			}
			expires = &exp
		}

		if !onlyJSON {
			fmt.Println()
			field("Type", header["typ"])
			field("Algorithm", header["alg"])
			field("Issuer", claims["iss"])
			field("Subject", claims["sub"])
			fmt.Println()
			field("Issued", issued)
			field("Expires", expires)
			field("Valid", validity)
			fmt.Println()
		}

		if !onlyInfo {
			json := string(indented)
			if !noColor {
				json = util.Colorize(json)
			}
			fmt.Println(json)
		}
	},
}

func init() {
	RootCmd.AddCommand(DecodeCmd)
	DecodeCmd.PersistentFlags().BoolVarP(&onlyInfo, "info", "i", false, "only show info about token")
	DecodeCmd.PersistentFlags().BoolVarP(&onlyJSON, "json", "j", false, "only show decoded token (no info)")
	DecodeCmd.PersistentFlags().BoolVarP(&noColor, "no-color", "c", false, "do not colorize json")
	DecodeCmd.PersistentFlags().BoolVarP(&forceFile, "file", "f", false, "force input to be filename (inverse of --input)")
	DecodeCmd.PersistentFlags().BoolVarP(&forceInput, "input", "d", false, "force direct input (inverse of --file)")
	DecodeCmd.PersistentFlags().BoolVarP(&forceOAuth, "oauth", "o", false, "force input to be OAuth 2.0 JSON response (inverse of --raw)")
	DecodeCmd.PersistentFlags().BoolVarP(&forceRaw, "raw", "r", false, "force input to be raw token (inverse of --oauth)")
}
