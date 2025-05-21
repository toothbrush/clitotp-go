/*
Copyright © 2022 Paul <paul@ravi>
*/
package cmd

import (
	"fmt"
	"net/url"
	"os"
	"path"

	"github.com/mdp/qrterminal/v3"

	"github.com/spf13/cobra"
	"github.com/toothbrush/clitotp-go/files"
)

var nameOverride string

var qrCmd = &cobra.Command{
	Use:   "qr KEYNAME",
	Short: "Generate a QR code to export TOTP secret",
	Long: `
Generate a QR code to export a TOTP secret into Google Authenticator or similar.

KEYNAME is a filename of something in your $HOME/.totp directory.
`,
	Args: cobra.ExactValidArgs(1),
	ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {

		if len(args) != 0 {
			return nil, cobra.ShellCompDirectiveNoFileComp
		}

		totps, err := files.ListTOTPs()
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", err)
			os.Exit(1)
		}

		return totps, cobra.ShellCompDirectiveNoFileComp
	},
	Run: func(cmd *cobra.Command, args []string) {
		filename, err := files.FullKeyPath(args[0])
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", err)
			os.Exit(1)
		}
		secret, err := files.RetrieveSecret(args[0])
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", err)
			os.Exit(1)
		}

		config := qrterminal.Config{
			Level:     qrterminal.L,
			Writer:    os.Stdout, // write straight to stdout...
			BlackChar: qrterminal.BLACK_BLACK,
			WhiteChar: qrterminal.WHITE_WHITE,
			QuietZone: 3,
		}

		issuer := url.QueryEscape(path.Base(filename))

		// allow overriding the "display name" in apps like Google Authenticator
		if nameOverride != "" {
			issuer = nameOverride
		}

		authURI := fmt.Sprintf("otpauth://totp/%s:clitotp-go?secret=%s&issuer=%s",
			issuer,
			secret,
			issuer)
		qrterminal.GenerateWithConfig(authURI, config)
	},
}

func init() {
	qrCmd.PersistentFlags().StringVarP(&nameOverride, "name", "n", "", "Display name of TOTP entry in Google Authenticator")
	rootCmd.AddCommand(qrCmd)
}
