package login

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type Auth struct {
	Token string
}

func (auth Auth) LoginWithToken() string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("----------- GITHUBY : GitHub API Management ---------\n   Enter GitHub Token : ")
	token, _ := reader.ReadString('\n')
	token = strings.Trim(token, "\n")
	auth.Token = token

	return token

}

// Login with Token
func (auth Auth) GetToken() string {
	return "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"
}
