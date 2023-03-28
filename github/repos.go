package github

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strconv"

	"github.com/tarathep/ghmgr/login"
	"github.com/tarathep/ghmgr/model"
)

type Repos struct {
	Auth  login.Auth
	Owner string
}

// for poc not use on this
func (repo Repos) SparseCheckout(url string, directory string) {
	//create new dir defore
	os.Mkdir("apimtool", 0777)
	repo.Exec("git", "apimtool", "clone", "--filter=blob:none", "--no-checkout", "--depth", "1", "--sparse", url)
	repo.Exec("git", "apimtool", "sparse-checkout", "add", ".github/workflows")
	repo.Exec("git", "apimtool", "checkout")

}

// for poc not use on this
func (repo Repos) Checkout(url string) {
	repo.Exec("git", "apimtool", "clone", url)
}

// for poc not use on this
func (repo Repos) Exec(name string, dir string, arg ...string) {
	cmd := exec.Command(name, arg...)
	cmd.Dir = dir
	stdout, err := cmd.Output()

	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println(string(stdout))
}

func (repo Repos) GetReposByTeam(teamName string) []string {
	github := GitHub{Auth: repo.Auth}
	var nameRepos []string
	total := 0
	for i := 0; true; i++ {
		pagex := strconv.Itoa((i + 1))

		_, repos := func(page string) (error, model.Repos) {
			err, statusCode, bodyBytes := github.Request("GET", "https://api.github.com/orgs/"+repo.Owner+"/teams/"+teamName+"/repos?page="+page, nil)
			if statusCode != 200 {
				log.Println(statusCode, github.GetMessage(bodyBytes))
			}
			if err != nil {
				return err, nil
			}

			repos := model.Repos{}

			json.Unmarshal(bodyBytes, &repos)

			return nil, repos
		}(pagex)

		if len(repos) == 0 {
			break
		}
		//doing..
		for _, repo := range repos {
			nameRepos = append(nameRepos, repo.Name)
			total += 1
			//fmt.Print("https://github.com/"+repo.Owner.Login+"/"+repo.Name, ",")
		}
	}
	return nameRepos
}

// func (team Repos) GetRepoList(teamName string) []string {
// 	var nameRepos []string
// 	total := 0

// 	for i := 0; true; i++ {
// 		pagex := strconv.Itoa((i + 1))

// 		_, repos := team.GetRepos(pagex)

// 		if len(repos) == 0 {
// 			break
// 		}
// 		//doing..
// 		for _, repo := range repos {
// 			nameRepos = append(nameRepos, repo.Name)
// 			total += 1
// 			fmt.Println("GetRepoList : ", total, repo.Name)
// 		}
// 	}
// 	return nameRepos
// }
