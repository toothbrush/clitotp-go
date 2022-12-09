package files

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"

	homedir "github.com/mitchellh/go-homedir"
)

func FullKeyPath(keyname string) (string, error) {
	home, err := homedir.Dir()
	if err != nil {
		return "", err
	}
	prefix := filepath.Join(home, ".totp")

	if !strings.HasSuffix(keyname, ".gpg") {
		keyname = keyname + ".gpg"
	}

	return path.Join(prefix, keyname), nil
}

func RetrieveSecret(keyname string) (string, error) {
	file, err := FullKeyPath(keyname)
	if err != nil {
		return "", err
	}

	// does it exist?
	if _, err := os.Stat(file); err == nil {
		// exists
	} else if errors.Is(err, os.ErrNotExist) {
		// does *not* exist
		return "", fmt.Errorf("File doesn't exist: %s", file)
	} else {
		// no idea.  stuff is weird.
		return "", err
	}

	out, err := exec.Command("gpg", "--batch", "--decrypt", file).Output()
	if err != nil {
		return "", err
	}

	return strings.ReplaceAll(string(out), " ", ""), nil
}
