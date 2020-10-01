/*
Copyright © 2020 Paul <paul@ravi>
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/toothbrush/clitotp-go/files"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all the TOTPs we know about",
	Long: `
Note that symlinks are not dealt with.  So don't have symlinks, i
guess.
`,
	Args: cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		totps, err := files.ListTOTPs()
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", err)
			os.Exit(1)
		}
		for _, totp := range totps {
			fmt.Println(totp)
		}
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}
