package csv

import (
	"encoding/csv"
	"log"
	"os"
	"time"

	"github.com/tarathep/ghmgr/model"
)

func WriteORGMemberReport(header string, name string, dataset []model.OrgMemberReport) error {
	time := time.Now().Format("20060102150405")
	file, err := os.Create(name + "-" + time + ".csv")
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	var data = [][]string{
		{header + time},
		{"No", "ID", "Username", "Name", "Email", "Team"},
	}

	//prepare dataset
	for _, d := range dataset {
		data = append(data, []string{d.No, d.ID, d.Username, d.Name, d.Email, d.Team})
	}

	for _, value := range data {
		err := writer.Write(value)
		if err != nil {
			log.Panic(err)
		}
	}
	return nil
}

func WriteTeamMemberReport(team string, header string, name string, dataset []model.TeamMemberReport) error {
	time := time.Now().Format("20060102150405")
	file, err := os.Create(name + "-" + time + ".csv")
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	var data = [][]string{
		{header + time},
		{"App/Team Name", team},
		{"No", "ID", "Username", "Name", "Email", "Role"},
	}

	//prepare dataset
	for _, d := range dataset {
		data = append(data, []string{d.No, d.ID, d.Username, d.Name, d.Email, d.Role})
	}

	for _, value := range data {
		err := writer.Write(value)
		if err != nil {
			log.Panic(err)
		}
	}
	return nil
}
