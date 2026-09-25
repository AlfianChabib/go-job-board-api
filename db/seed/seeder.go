package main

import (
	"AlfianChabib/go-job-board-api/internal/config"
	"AlfianChabib/go-job-board-api/internal/database"
	"AlfianChabib/go-job-board-api/internal/model/domain"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"gorm.io/gorm/clause"
)

func SeedSkils() {
	env := config.LoadEnv()
	db := database.OpenConnection(env)

	file, err := os.Open("db/seed/skills-dataset.csv")
	if err != nil {
		log.Fatal("Error opening CSV:", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)

	var skills []domain.Skill
	seen := make(map[string]int)

	_, err = reader.Read()
	if err != nil && err != io.EOF {
		log.Fatal("Error reading CSV header:", err)
	}

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Println(err)
			continue
		}

		if len(record) < 3 {
			continue
		}

		name := strings.TrimSpace(record[1])
		if name == "" {
			continue
		}

		label := strings.TrimSpace(record[2])
		abbr := ""
		if len(record) > 3 {
			abbr = strings.TrimSpace(record[3])
		}

		if idx, exists := seen[name]; exists {
			if abbr != "" && skills[idx].Abbreviation == "" {
				skills[idx].Abbreviation = abbr
			}
			if label != "" && skills[idx].Label == "" {
				skills[idx].Label = label
			}
			continue
		}

		seen[name] = len(skills)
		skills = append(skills, domain.Skill{
			Name:         name,
			Label:        label,
			Abbreviation: abbr,
		})
	}

	result := db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "name"}},
		DoUpdates: clause.AssignmentColumns([]string{"label", "abbreviation"}),
	}).CreateInBatches(&skills, 100)

	if result.Error != nil {
		log.Fatal("Error creating skills:", result.Error)
	}

	fmt.Println("Successfully created/updated", result.RowsAffected, "skills")
}

func main() {
	SeedSkils()
}
