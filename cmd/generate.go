/*
Copyright © 2020 Paul <paul@ravi>
*/
package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
	"time"

	"golang.design/x/clipboard"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/pquerna/otp/totp"
	"github.com/spf13/cobra"
	"github.com/toothbrush/clitotp-go/files"
)

var copyToClipboard bool

var generateCmd = &cobra.Command{
	Use:   "generate KEYNAME",
	Short: "Generate a TOTP",
	Long: `
Generate a TOTP and display it on stdout.  Optionally, use the --copy/-c flag to
copy it to the clipboard.

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
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		prefix := filepath.Join(home, ".totp")

		keyname := args[0]
		if !strings.HasSuffix(keyname, ".gpg") {
			keyname = keyname + ".gpg"
		}

		file := path.Join(prefix, keyname)
		out, err := exec.Command("gpg", "--batch", "--decrypt", file).Output()
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", err)
			return
		}

		noSpaceString := strings.ReplaceAll(string(out), " ", "")
		otp, err := totp.GenerateCode(
			noSpaceString,
			time.Now(),
		)

		if copyToClipboard {
			// Init returns an error if the package is not ready for use.
			err := clipboard.Init()
			if err != nil {
				fmt.Fprintf(os.Stderr, "Oh no, couldn't initialise clipboard: %v\n", err)
			} else {
				clipboard.Write(clipboard.FmtText, []byte(otp))
				fmt.Fprintln(os.Stderr, "TOTP copied to clipboard.")
			}
		}

		fmt.Printf("%s", otp)
	},
}

func init() {
	generateCmd.Flags().BoolVarP(&copyToClipboard, "copy", "c", false, "copy to clipboard")
	rootCmd.AddCommand(generateCmd)
}
