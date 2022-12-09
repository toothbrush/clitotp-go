/*
Copyright © 2021 Paul <paul@ravi>
*/
package cmd

import (
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/spf13/cobra"

	"github.com/toothbrush/clitotp-go/cli"
	"github.com/toothbrush/clitotp-go/files"
)

var Force = false

var removeCmd = &cobra.Command{
	Use:   "rm KEYNAME",
	Short: "Remove a secret from the store",
	Args:  cobra.ExactValidArgs(1),
	Long: `
Removes a secret from the store.

Sometimes, you just don't need a particular secret anymore.  Remove it!
`,
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
		keyname_nogpg := strings.TrimSuffix(path.Base(filename), ".gpg")

		if _, err := os.Stat(filename); err == nil {
			// something exists at `filename` path, yay.
			var delete = false
			if Force {
				delete = true
			} else {
				// If --force isn't set, ask the user what to do.
				delete = cli.YNConfirm(fmt.Sprintf("Are you sure you would like to delete '%s'?", keyname_nogpg))
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
