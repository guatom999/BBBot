package botRepositories

import (
	"encoding/csv"
	"errors"
	"log"
	"os"

	"github.com/ahmdrz/goinsta/v2"
)

type (
	IBotRepository interface {
		GetFollowers(target string, isOldFollowers bool) error
		GetLastFollowers() []string
		GetNowFollowers() []string
	}

	botRepository struct {
		instaBot *goinsta.Instagram
	}
)

func NewBotRepository(instaBot *goinsta.Instagram) IBotRepository {
	return &botRepository{instaBot: nil}
}

func (r *botRepository) GetFollowers(target string, isWriteOldFile bool) error {

	user, err := r.instaBot.Profiles.ByName(target)
	if err != nil {
		log.Printf("Failed to get followers : %s", err.Error())
		return err
	}
	// defer r.instaBot.Logout()

	followers := user.Followers()

	var file *os.File

	if isWriteOldFile {
		file, err = os.Create("followers_last.csv")
	} else {
		file, err = os.Create("followers_now.csv")
	}

	if err != nil {
		log.Fatalf("Error: Failed to create file: %v", err)
		return errors.New("failed to create file")
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	for followers.Next() {
		for _, follower := range followers.Users {
			writer.Write([]string{follower.Username})
		}
	}

	log.Printf("Write followers_now.csv successfully")

	return nil
}

func (r *botRepository) GetLastFollowers() []string {

	lastFile, err := os.Open("followers_last.csv")
	if err != nil {
		log.Fatalf("Error: Failed to read file: %v", err)
		return make([]string, 0)
	}

	lastReader := csv.NewReader(lastFile)
	lastRecords, err := lastReader.ReadAll()
	if err != nil {
		log.Fatalf("Error: Failed to read file: %v", err)
		return make([]string, 0)
	}

	result := make([]string, 0)

	for _, record := range lastRecords {
		result = append(result, record[0])
	}

	return result

}

func (r *botRepository) GetNowFollowers() []string {

	lastFile, err := os.Open("followers_now.csv")
	if err != nil {
		log.Fatalf("Error: Failed to read file: %v", err)
		return make([]string, 0)
	}

	reader := csv.NewReader(lastFile)
	nowRecords, err := reader.ReadAll()
	if err != nil {
		log.Fatalf("Error: Failed to read file: %v", err)
		return make([]string, 0)
	}

	result := make([]string, 0)

	for _, record := range nowRecords {
		result = append(result, record[0])
	}

	return result

}
