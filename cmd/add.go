/*
Copyright © 2020 Paul <paul@ravi>
*/
package cmd

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/toothbrush/clitotp-go/cli"
	"github.com/toothbrush/clitotp-go/files"
)

var addCmd = &cobra.Command{
	Use:   "add KEYNAME",
	Short: "Add a secret to the store",
	Args:  cobra.ExactArgs(1),
	Long: `
Adds a secret to your TOTP store.

This command will spawn an interactive session to capture a secret, e.g. from
a new website you've joined, and encrypt it in your TOTP store.
`,
	ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return nil, cobra.ShellCompDirectiveNoFileComp
	},
	Run: func(cmd *cobra.Command, args []string) {
		filename, err := files.FullKeyPath(args[0])
		if err != nil {
			log.Panic(err)
		}

		recipient := viper.GetString("recipient")

		if recipient == "" {
			log.Panic("Please set a GPG recipient to encrypt to!")
		}
		fmt.Fprintf(os.Stderr, "Encrypting to %s.\n", recipient)
		fmt.Fprintf(os.Stderr, "Will store secret in: %s\n", filename)

		if _, err := os.Stat(filename); err == nil {
			// something exists at `filename` path.
			overwrite := cli.YNConfirm("File exists.  Overwrite?")
			if !overwrite {
				fmt.Fprintln(os.Stderr, "Aborting.")
				os.Exit(0)
			}
		}

		var newSecret string

		fmt.Fprint(os.Stderr, "Give me the secret (C-c cancels): ")
		// Ask the user for the new secret:
		scanner := bufio.NewScanner(os.Stdin)
		if scanner.Scan() {
			newSecret = scanner.Text()
		} else {
			log.Panic("oh no, couldn't read input.")
		}

		gpgCmd := exec.Command("gpg", "--batch", "--encrypt", "--recipient", recipient)
		stdin, err := gpgCmd.StdinPipe()
		if err != nil {
			log.Panic(err)
		}
		stdout, err := gpgCmd.StdoutPipe()
		if err != nil {
			log.Panic(err)
		}
		err = gpgCmd.Start()
		if err != nil {
			log.Panic(err)
		}

		// open the out file for writing
		outfileHandle, err := os.Create(filename)
		if err != nil {
			log.Panic(err)
		}

		io.Copy(stdin, bytes.NewBufferString(newSecret))
		stdin.Close()

		writer := bufio.NewWriter(outfileHandle)

		io.Copy(writer, stdout)
		writer.Flush()
		outfileHandle.Close()
		err = gpgCmd.Wait()
		if err != nil {
			log.Panic(err)
		}

		fmt.Fprintln(os.Stderr, "Encrypted and saved.")
	},
}

func init() {
	rootCmd.AddCommand(addCmd)
	addCmd.PersistentFlags().StringP("recipient", "r", "", "GPG key id to encrypt TOTP value")
	viper.BindPFlag("recipient", addCmd.PersistentFlags().Lookup("recipient"))
}
