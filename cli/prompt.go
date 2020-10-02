package cli

import (
	"github.com/AlecAivazis/survey/v2"
)

func YNConfirm(promptText string) bool {
	name := false
	prompt := &survey.Confirm{
		Message: promptText,
	}
	survey.AskOne(prompt, &name)
	return name
}
