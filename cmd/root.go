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
var forceFile bool
var forceInput bool
var forceRaw bool
var forceOAuth bool

type Token struct {
	AccessToken string `json:"access_token"`
}

var RootCmd = &cobra.Command{
	Use:   "jwt [token|path]",
	Short: "jwt can be used the debug JWT tokens.",
	Long: `A simple jwt debugging tool written in Go.
By default, jwt accepts both direct input as filenames. To force either
one of the two use the --input and --file flags respectivly.

The token can also be passed in via stdin.

jwt can also be used to parse tokens from a standard OAuth 2.0 JSON response
body. For example, passing:

  { "access_token": "ey10...", "expires": ... }

is equivalent to passing "eq10..." directly. To suppress this behaviour, use the
--raw flag, to force it --oauth.

	`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		if debug && !silent {
			log.SetLevel(log.DebugLevel)
			log.WithField("flag", "--debug").Debug("Debug on")
		}

		if silent {
			log.SetLevel(log.ErrorLevel)
		}

		if onlyJSON && !onlyInfo {
			log.WithField("flag", "--json").Debug("Only printing JSON")
		}

		if onlyInfo {
			log.WithField("flag", "--info").Debug("Only printing info")
			onlyJSON = false
		}

		if forceFile && !forceInput {
			log.WithField("flag", "--file").Debug("Forcing file mode")
		}

		if forceInput {
			log.WithField("flag", "--input").Debug("Forcing input mode")
		}

		if forceOAuth && !forceRaw {
			log.WithField("flag", "--oauth").Debug("Forcing oauth mode")
		}

		if forceRaw {
			log.WithField("flag", "--raw").Debug("Forcing raw mode")
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
	RootCmd.PersistentFlags().BoolVarP(&debug, "verbose", "v", false, "show debug output")
	RootCmd.PersistentFlags().BoolVarP(&silent, "silent", "s", false, "do not print status messages")
	RootCmd.PersistentFlags().BoolVarP(&onlyInfo, "info", "i", false, "only show info about token")
	RootCmd.PersistentFlags().BoolVarP(&onlyJSON, "json", "j", false, "only show decoded token (no info)")
	RootCmd.PersistentFlags().BoolVarP(&noColor, "no-color", "c", false, "do not colorize json")
	RootCmd.PersistentFlags().BoolVarP(&forceFile, "file", "f", false, "force input to be filename (inverse of --input)")
	RootCmd.PersistentFlags().BoolVarP(&forceInput, "input", "d", false, "force direct input (inverse of --filename)")
	RootCmd.PersistentFlags().BoolVarP(&forceOAuth, "oauth", "o", false, "force input to be OAuth 2.0 JSON response (inverse of --raw)")
	RootCmd.PersistentFlags().BoolVarP(&forceRaw, "raw", "r", false, "force input to be raw token (inverse of --oauth)")
}
