/*
Copyright © 2020 Paul <paul@ravi>
*/
package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all the TOTPs we know about",
	Long: `
Note that symlinks are not dealt with.  So don't have symlinks, i
guess.
`,
	Run: func(cmd *cobra.Command, args []string) {
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		prefix := strings.Join([]string{filepath.Join(home, ".totp"), string(filepath.Separator)}, "")

		err = filepath.Walk(prefix,
			func(path string, info os.FileInfo, err error) error {
				if err != nil {
					return err
				}
				relative, err := filepath.Rel(prefix, path)
				if err != nil {
					return err
				}

				// Real-deal files will end with ".gpg"
				if strings.HasSuffix(relative, ".gpg") {
					fmt.Println(strings.TrimSuffix(relative, ".gpg"))
				}
				return nil
			})
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}
