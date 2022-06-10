package services

import (
	"github.com/kevincobain2000/go-app-reviews-scraper/app"
)

// AutoMigrate will auto migrate the database
// This will not delete the data when ran again
func AutoMigrate() {
	db := app.NewDB()
	err := db.AutoMigrate(&ReviewModel{})
	if err != nil {
		panic(err)
	}
	err = db.AutoMigrate(&ReviewCountsModel{})
	if err != nil {
		panic(err)
	}
}
