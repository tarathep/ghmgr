package login

import (
	"log"
	"os"

	"github.com/fatih/color"
)

type Auth struct {
	Token   string
	Owner   string
	Version string
}

func (auth Auth) LoginWithToken() {
	// reader := bufio.NewReader(os.Stdin)
	// fmt.Print("----------- GitHub Manager : GitHub API Management ---------\n   Enter GitHub Token : ")
	// token, _ := reader.ReadString('\n')
	// token = strings.Trim(token, "\n")

	// fmt.Println(token)

	// fmt.Println("xx>" + auth.GetToken())
}

// Login with Token
func (auth Auth) SetToken(token string) {

	if err := os.Setenv("GHMGR_TOKEN", token); err != nil {
		log.Panic(err)
	}

	if auth.GetToken() == token {
		color.New(color.FgHiGreen).Print("Login with Token Success\n")
		println(auth.GetToken())
	} else {
		color.New(color.FgHiRed).Print("Login with Token Error\n")
		os.Exit(1)
	}

}

// Get Login with Token
func (auth Auth) GetToken() string {
	return os.Getenv("GHMGR_TOKEN")
}

func (auth Auth) SetOwner(owner string) {
	if err := os.Setenv("GHMGR_OWNER", owner); err != nil {
		log.Panic(err)
	}

	if auth.GetOwner() == owner {
		color.New(color.FgHiGreen).Print("Set Owner Success\n")
	} else {
		color.New(color.FgHiRed).Print("Set Owner Error\n")
		os.Exit(1)
	}
}

func (auth Auth) GetOwner() string {
	return os.Getenv("GHMGR_OWNER")
}
