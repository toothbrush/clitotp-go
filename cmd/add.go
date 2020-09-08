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
	"path"
	"strings"

	"github.com/spf13/cobra"
)

var addCmd = &cobra.Command{
	Use:   "add KEYNAME",
	Short: "Add a secret to the store",
	Args:  cobra.ExactArgs(1),
	Long: `
Spawn an interactive session to capture a secret, e.g. from a new
website you've joined, and encrypt it in your TOTP store.
`,
	Run: func(cmd *cobra.Command, args []string) {
		prefix := "/home/paul/.totp/"

		keyname := args[0]
		if !strings.HasSuffix(keyname, ".gpg") {
			keyname = keyname + ".gpg"
		}

		filename := path.Join(prefix, keyname)

		fmt.Fprintf(os.Stderr, "Will insert into: %s\n", filename)

		var newSecret string

		fmt.Fprint(os.Stderr, "Give me the secret (C-c cancels): ")
		// Ask the user for the new secret:
		fmt.Scanln(&newSecret)

		recipient := "0xF2846B1A0D32C442"
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
}
