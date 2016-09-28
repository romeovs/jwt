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

var debug bool
var silent bool
var onlyInfo bool
var onlyJSON bool
var noColor bool

type Token struct {
	AccessToken string `json:"access_token"`
}

var RootCmd = &cobra.Command{
	Use:   "jwt [token|path]",
	Short: "jwt can be used the debug JWT tokens.",
	Long:  "A simple jwt debugging tool written in Go.",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		if debug {
			log.SetLevel(log.DebugLevel)
			log.Debug("Debug on")
		}

		if silent {
			log.SetLevel(log.ErrorLevel)
			log.Debug("Debug on")
		}

		if onlyJSON && !onlyInfo {
			log.WithField("flag", "--json").Debug("Only printing JSON")
		}

		if onlyInfo {
			log.WithField("flag", "--info").Debug("Only printing info")
			onlyJSON = false
		}
	},
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

		issued := time.Unix(int64(claims["iat"].(float64)), 0)
		expires := time.Unix(int64(claims["exp"].(float64)), 0)
		now := time.Now().Round(time.Second)

		validity := "token is valid"
		if issued.After(now) {
			validity = fmt.Sprintf("token is not valid for %s", issued.Sub(now))
		}
		if now.After(expires) {
			validity = fmt.Sprintf("token is expired for %s", now.Sub(expires))
		}

		if !onlyJSON {
			fmt.Println()
			field("Type", header["typ"])
			field("Algorithm", header["alg"])
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
	RootCmd.PersistentFlags().BoolVarP(&debug, "verbose", "v", false, "show debug output")
	RootCmd.PersistentFlags().BoolVarP(&silent, "silent", "s", false, "do not print status messages")
	RootCmd.PersistentFlags().BoolVarP(&onlyInfo, "info", "i", false, "only show info about token")
	RootCmd.PersistentFlags().BoolVarP(&onlyJSON, "json", "j", false, "only show decoded token (no info)")
	RootCmd.PersistentFlags().BoolVarP(&noColor, "no-color", "c", false, "do not colorize json")
}
