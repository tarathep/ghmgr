package login

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type Auth struct {
}

// Login with Token
func (Auth) Token() string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("----------- UPDATE GITHUB REPOSITORY ROLE TEAM ---------\n   Enter GitHub Token : ")
	token, _ := reader.ReadString('\n')
	token = strings.Trim(token, "\n")
	return token
}
