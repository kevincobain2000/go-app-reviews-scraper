package services

import (
	"errors"
	"time"

	"github.com/kevincobain2000/go-app-reviews-scraper/app"
	"gorm.io/gorm"
)

// ReviewsRepository is the respository for the sales data
type ReviewsRepository struct {
	db *gorm.DB
}

// NewReviewsRepository the constructor for NewReviewsRepository
// db is injected as DUI to the constructor
func NewReviewsRepository() *ReviewsRepository {
	return &ReviewsRepository{
		db: app.NewDB(),
	}
}

// FindOrNewReviews finds the review count or creates a new one
// looks for existing review by store, app name, username with it's rating date
// username and rating_date are scraped from the review page from the review card
func (r *ReviewsRepository) FindOrNewReviews(reviews Reviews) ([]ReviewModel, error) {
	query := `app_name = ?
		AND store = ?
		AND username = ?
		AND rated_at = ?
		AND deleted_at IS NULL`

	newReviews := []ReviewModel{}
	for i := 0; i < len(reviews.Ratings); i++ {
		var review = ReviewModel{}
		result := r.db.Where(
			query,
			reviews.AppName,
			reviews.Store,
			reviews.Usernames[i],
			reviews.Datetimes[i],
		).First(&review)

		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			review, err := r.insertReview(
				reviews.AppName,
				reviews.Store,
				reviews.Usernames[i],
				reviews.Titles[i],
				reviews.Bodies[i],
				reviews.Ratings[i],
				reviews.Datetimes[i],
			)
			if err != nil {
				return newReviews, err
			}
			newReviews = append(newReviews, review)
		}
	}

	return newReviews, nil
}

// FindLastReviewCount finds the last review count
// ID desc is used to get the last record
// ORDER BY `review_counts`.`id` DESC LIMIT 1
func (r *ReviewsRepository) FindLastReviewCount(reviews Reviews) (ReviewCountsModel, error) {

	var reviewCount = ReviewCountsModel{}
	query := `app_name = ?
		AND store = ?
		AND deleted_at IS NULL`
	result := r.db.Where(
		query,
		reviews.AppName,
		reviews.Store,
	).Last(&reviewCount)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return reviewCount, nil
	}
	return reviewCount, nil
}

// FindOrNewReviewCount finds the review count or creates a new one
// looks for existing review count by store, app name, with all the ratings and total fetched by scraping
// if not found a match, it creates a new one
// previous review count summary is kept as it is and new one is used at the time of fetching
// Also see @FindLastReviewCount
func (r *ReviewsRepository) FindOrNewReviewCount(reviews Reviews) (ReviewCountsModel, error) {

	var reviewCount = ReviewCountsModel{}
	query := `app_name = ?
		AND store = ?
		AND total = ?
		AND rating_1_percentage = ?
		AND rating_2_percentage = ?
		AND rating_3_percentage = ?
		AND rating_4_percentage = ?
		AND rating_5_percentage = ?
		AND deleted_at IS NULL`
	result := r.db.Where(
		query,
		reviews.AppName,
		reviews.Store,
		reviews.Total,
		reviews.Rating1Percentage,
		reviews.Rating2Percentage,
		reviews.Rating3Percentage,
		reviews.Rating4Percentage,
		reviews.Rating5Percentage,
	).Last(&reviewCount)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return r.InsertReviewCount(reviews)
	}
	return reviewCount, nil
}

// insertReview inserts a new review
// that's it
// DELETED_AT is NULL by default
func (r *ReviewsRepository) insertReview(appName, store, username, title, body string, rating int, rated_at time.Time) (ReviewModel, error) {
	now := time.Now()
	review := ReviewModel{
		AppName:   appName,
		Store:     store,
		Username:  username,
		Title:     title,
		Body:      body,
		Rating:    rating,
		RatedAt:   &rated_at,
		CreatedAt: &now,
		UpdatedAt: &now,
	}
	result := r.db.Create(&review)
	return review, result.Error
}

// InsertReviewCount inserts a new review count
// that's it
// DELETED_AT is NULL by default
func (r *ReviewsRepository) InsertReviewCount(reviews Reviews) (ReviewCountsModel, error) {
	now := time.Now()
	reviewCount := ReviewCountsModel{
		AppName:           reviews.AppName,
		Store:             reviews.Store,
		Total:             reviews.Total,
		Rating1Percentage: reviews.Rating1Percentage,
		Rating2Percentage: reviews.Rating2Percentage,
		Rating3Percentage: reviews.Rating3Percentage,
		Rating4Percentage: reviews.Rating4Percentage,
		Rating5Percentage: reviews.Rating5Percentage,
		CreatedAt:         &now,
		UpdatedAt:         &now,
	}
	result := r.db.Create(&reviewCount)
	return reviewCount, result.Error
}
