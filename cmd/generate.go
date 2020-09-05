/*
Copyright © 2020 Paul <paul@ravi>
*/
package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path"
	"strings"
	"time"

	"github.com/pquerna/otp/totp"
	"github.com/spf13/cobra"
)

var generateCmd = &cobra.Command{
	Use:   "generate KEY",
	Short: "Generate a TOTP",
	Args:  cobra.ExactArgs(1),
	Long: `
Pass a filename of something in your $HOME/.totp directory.
`,
	Run: func(cmd *cobra.Command, args []string) {
		prefix := "/home/paul/.totp"
		keyname := args[0]
		if !strings.HasSuffix(keyname, ".gpg") {
			keyname = keyname + ".gpg"
		}

		file := path.Join(prefix, keyname)
		out, err := exec.Command("gpg", "--decrypt", file).Output()
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s", err)
			return
		}

		otp, err := totp.GenerateCode(
			string(out),
			time.Now(),
		)

		fmt.Printf("%s", otp)
	},
}

func init() {
	rootCmd.AddCommand(generateCmd)
}
