package login

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

type Auth struct {
	Token string
}

func (auth Auth) LoginWithToken() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("----------- GitHub Manager : GitHub API Management ---------\n   Enter GitHub Token : ")
	token, _ := reader.ReadString('\n')
	token = strings.Trim(token, "\n")

	// c := encrypt("hello", "world")
	// fmt.Println(c)
	// fmt.Println(decrypt("hello", c))
	// fmt.Println(decrypt("hello", "c2932347953ad4a4-25f496d260de9c150fc9e4c6-20bc1f8439796cc914eb783b9996a8d9c32d45e2df"))

	os.Setenv("GHP_TOKEN", token)
}

// Login with Token
func (auth Auth) SetToken(token string) {
	fmt.Print(token)
	if err := os.Setenv("GHP_TOKEN", token); err != nil {
		log.Panic(err)
	}
}

// Get Login with Token
func (auth Auth) GetToken() string {
	return os.Getenv("GHP_TOKEN")
	//return "ghp_6JL92oupGZThnb38UyJtlf8LO2Uw6j2OvnZy"
}
