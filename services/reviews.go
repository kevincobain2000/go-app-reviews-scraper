package services

import (
	"time"
)

// Reviews is a struct that represents the reviews information fetched from the scraper
// This is the data that is stored in the database
// This is the data that is returned from the scraper
// AppName and Store are set by cli args
// The arrays example: Ratings, Usernames, Titles, Body, Datetimes have items in the same order for reviews
// Non array fields are used to set the overall ratings such as Total, Rating1Percentage, Rating2Percentage..
type Reviews struct {
	AppName           string
	Store             string
	Ratings           []int
	Usernames         []string
	Titles            []string
	Datetimes         []time.Time
	Body              []string
	Total             int
	Rating1Percentage int
	Rating2Percentage int
	Rating3Percentage int
	Rating4Percentage int
	Rating5Percentage int
}

type ReviewModel struct {
	ID int `json:"id" gorm:"column:id;primary_key;AUTO_INCREMENT"`
	// AppName is not fetched by scraper but from cli args by the user
	AppName string `json:"app_name" gorm:"column:app_name;type:varchar(64); NOT NULL"`
	// Store is not fetched or judged from the URL or by scraper but from cli args by the user
	Store string `json:"store" gorm:"column:store;type:varchar(16); NOT NULL"`

	// Following items are fetched by scraper
	Username string     `json:"username" gorm:"column:username;type:string; NOT NULL"`
	Title    string     `json:"title" gorm:"column:title;type:string; NOT NULL"`
	Body     string     `json:"body" gorm:"column:body;type:string; NOT NULL"`
	Rating   int        `json:"rating" gorm:"column:rating;type:tinyint(1); NOT NULL"`
	RatedAt  *time.Time `json:"rated_at" gorm:"type:timestamp; NOT NULL"`

	// Basic timestamps
	CreatedAt *time.Time `json:"created_at,omitempty" gorm:"type:timestamp null"`
	UpdatedAt *time.Time `json:"updated_at,omitempty" gorm:"type:timestamp null"`
	DeletedAt *time.Time `json:"deleted_at,omitempty" gorm:"type:timestamp null"`
}

func (ReviewModel) TableName() string {
	return "reviews"
}
func (ReviewCountsModel) TableName() string {
	return "review_counts"
}

type ReviewCountsModel struct {
	ID int `json:"id" gorm:"column:id;primary_key;AUTO_INCREMENT"`
	// AppName is not fetched by scraper but from cli args by the user
	AppName string `json:"app_name" gorm:"column:app_name;type:varchar(64); NOT NULL"`
	// Store is not fetched or judged from the URL or by scraper but from cli args by the user
	Store string `json:"store" gorm:"column:store;type:varchar(16); NOT NULL"`

	// Following items are fetched by scraper
	Total             int `json:"total" gorm:"column:total;type:integer; NOT NULL"`
	Rating1Percentage int `json:"rating_1_percentage" gorm:"column:rating_1_percentage;type:tinyint(1); NOT NULL"`
	Rating2Percentage int `json:"rating_2_percentage" gorm:"column:rating_2_percentage;type:tinyint(1); NOT NULL"`
	Rating3Percentage int `json:"rating_3_percentage" gorm:"column:rating_3_percentage;type:tinyint(1); NOT NULL"`
	Rating4Percentage int `json:"rating_4_percentage" gorm:"column:rating_4_percentage;type:tinyint(1); NOT NULL"`
	Rating5Percentage int `json:"rating_5_percentage" gorm:"column:rating_5_percentage;type:tinyint(1); NOT NULL"`

	// Basic timestamps
	CreatedAt *time.Time `json:"created_at,omitempty" gorm:"type:timestamp null"`
	UpdatedAt *time.Time `json:"updated_at,omitempty" gorm:"type:timestamp null"`
	DeletedAt *time.Time `json:"deleted_at,omitempty" gorm:"type:timestamp null"`
}
