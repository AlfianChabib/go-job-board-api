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
		}

		skills = append(skills, domain.Skill{
			Name:  record[1],
			Label: record[2],
		})
	}

	result := db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "name"}},
		DoNothing: true,
	}).CreateInBatches(&skills, 100)

	if result.Error != nil {
		log.Fatal("Error creating skills:", result.Error)
	}

	fmt.Println("Successfully created", result.RowsAffected, "skills")
}

func main() {
	SeedSkils()
}
