package cmd

import (
	"github.com/apex/log"
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
	Use: "jwt",
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
		DecodeCmd.Run(cmd, args)
	},
}

func init() {
	RootCmd.PersistentFlags().BoolVarP(&debug, "verbose", "v", false, "show debug output")
	RootCmd.PersistentFlags().BoolVarP(&silent, "silent", "s", false, "do not print status messages")
}
