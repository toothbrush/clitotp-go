/*
Copyright © 2021 Paul <paul@ravi>
*/
package cmd

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"

	"github.com/toothbrush/clitotp-go/cli"
)

var Force = false

var removeCmd = &cobra.Command{
	Use:   "rm KEYNAME",
	Short: "Remove a secret from the store",
	Args:  cobra.ExactArgs(1),
	Long: `
Sometimes, you just don't need a particular secret anymore.  Remove it!
`,
	ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return nil, cobra.ShellCompDirectiveNoFileComp
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
		keyname_nogpg := strings.TrimSuffix(keyname, ".gpg")

		filename := path.Join(prefix, keyname)

		if _, err := os.Stat(filename); err == nil {
			// something exists at `filename` path, yay.
			var delete = false
			if Force {
				delete = true
			} else {
				// If --force isn't set, ask the user what to do.
				delete = cli.YNConfirm(fmt.Sprintf("Are you sure you would like to delete %s?", keyname_nogpg))
			}

			if delete {
				err = os.Remove(filename)
				if err == nil {
					fmt.Fprintf(os.Stderr, "removed '%s'\n", filename)
				} else {
					fmt.Fprintf(os.Stderr, "Error!  Unable to remove '%s'!\n", filename)
					os.Exit(1)
				}
			} else {
				fmt.Fprintln(os.Stderr, "Cancelled by user, aborting.")
			}
		} else {
			// couldn't find mentioned secret file
			fmt.Fprintf(os.Stderr, "Error, could not find secret named '%s'.  Aborting.\n", keyname_nogpg)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(removeCmd)
	removeCmd.Flags().BoolVarP(&Force, "force", "f", false, "delete without confirmation")
}
