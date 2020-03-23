package conf

import (
	"errors"
	"github.com/jinzhu/gorm"
	"log"
	"os"
)

var (
	ApiKey string
	Debug  bool
	Db     *gorm.DB
)

// Validate validate settings fields or panic
func Validate() (valid bool, err error) {
	valid = ApiKey != ""
	if !valid {
		err = errors.New("API_KEY is missing")
	}
	return valid, err
}

// Setup startup bot configuration
func Setup() {

	db, err := gorm.Open("sqlite3", "courses_bot.db")
	if err != nil {
		log.Panic("failed to connect database")
	}

	ApiKey = os.Getenv("API_KEY")
	Debug = os.Getenv("DEBUG") == "1"
	Db = db

	valid, err := Validate()
	if !valid {
		log.Fatal(err)
	}
}
