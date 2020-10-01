package files

import (
	"os"
	"path/filepath"
	"strings"

	homedir "github.com/mitchellh/go-homedir"
)

func ListTOTPs() ([]string, error) {
	home, err := homedir.Dir()
	if err != nil {
		return nil, err
	}
	prefix := strings.Join([]string{filepath.Join(home, ".totp"), string(filepath.Separator)}, "")

	var list = []string{}

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
				list = append(list, strings.TrimSuffix(relative, ".gpg"))
			}
			return nil
		})
	if err != nil {
		return nil, err
	}

	return list, nil
}
