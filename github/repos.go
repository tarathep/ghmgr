package github

import (
	"encoding/json"
	"errors"
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

func (repo Repos) GetReposByTeam(teamName string) []model.Repo {
	github := GitHub{Auth: repo.Auth}

	var repositories []model.Repo

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
			repositories = append(repositories, repo)
			total += 1
		}
	}
	return repositories
}

// Remove a repository from a team
func (repo Repos) RemovingRepositoryTeam(teamName, repoName string) (error, bool) {
	github := GitHub{Auth: repo.Auth}

	err, statusCode, bodyBytes := github.Request("DELETE", "https://api.github.com/orgs/"+repo.Owner+"/teams/"+teamName+"/repos/"+repo.Owner+"/"+repoName, nil)
	if statusCode != 204 {
		//fmt.Println(statusCode, github.GetMessage(bodyBytes))
		return errors.New(github.GetMessage(bodyBytes)), false
	}
	if err != nil {
		fmt.Println(err.Error())
		return err, false
	}
	return nil, true

}
