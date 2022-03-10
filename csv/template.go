package csv

import (
	"encoding/csv"
	"log"
	"os"
	"time"

	"github.com/tarathep/ghmgr/model"
)

type Template struct{}

func (Template) ReadFile(filename string) (project string, csvlist []model.CSV) {

	records, err := readData(filename)

	if err != nil {
		log.Fatal(err)
	}

	for i, record := range records {
		switch i {
		case 0:
			project = record[1]
			// fmt.Println("TeamName/Project :>", record[1])
		case 1:
		default:
			csv := model.CSV{
				ID:             record[0],
				MemberName:     record[1],
				Email:          record[2],
				Role:           record[3],
				GitHubTeamRole: record[4],
				GitHubUser:     record[5]}

			//fmt.Println((i - 1), csv.Email)
			csvlist = append(csvlist, csv)
		}
	}
	return project, csvlist
}

func readData(fileName string) ([][]string, error) {

	f, err := os.Open(fileName)

	if err != nil {
		return [][]string{}, err
	}

	defer f.Close()

	r := csv.NewReader(f)

	records, err := r.ReadAll()

	if err != nil {
		return [][]string{}, err
	}

	return records, nil
}

func (Template) WriteCSV(teamName string, dataset []model.CSV) bool {
	file, err := os.Create(teamName + "-" + time.Now().Format("20060102150405") + ".csv")
	if err != nil {
		log.Panic(err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	var data = [][]string{
		{"TeamName/Project", teamName},
		{"ID (Employee)", "MemberName", "Email", "Role", "GitHub Team Role", "UserName(GitHub)"},
	}

	//prepare dataset
	for _, d := range dataset {
		data = append(data, []string{d.ID, d.MemberName, d.Email, d.Role, d.GitHubTeamRole, d.GitHubUser})
	}

	for _, value := range data {
		err := writer.Write(value)
		if err != nil {
			log.Panic(err)
		}
	}

	return true
}
