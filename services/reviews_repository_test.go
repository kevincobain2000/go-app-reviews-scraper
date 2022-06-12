package services

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func init() {
	Setup()
}

func TestInsertReview(t *testing.T) {
	r := NewReviewsRepository()
	assert.NotNil(t, r)
	now := time.Now()

	rating := 5
	review, err := r.insertReview("app_name", "store", "username", "title", "body", rating, now)
	assert.Nil(t, err)
	assert.Equal(t, "app_name", review.AppName)
	assert.Equal(t, "store", review.Store)
	assert.Equal(t, "username", review.Username)
	assert.Equal(t, "title", review.Title)
	assert.Equal(t, "body", review.Body)
	assert.Equal(t, rating, review.Rating)
	assert.Equal(t, now, *review.RatedAt)
}

func TestFindOrNewReviewCount(t *testing.T) {
	r := NewReviewsRepository()

	ReviewCountsModel, err := r.FindOrNewReviewCount(Reviews{})
	assert.Nil(t, err)
	assert.Equal(t, "", ReviewCountsModel.AppName)
	assert.Equal(t, "", ReviewCountsModel.Store)
	assert.Equal(t, 0, ReviewCountsModel.Total)
	//.. on so on

	now := time.Now()
	reviews := Reviews{
		AppName:           "app",
		Store:             "store",
		Usernames:         []string{"username"},
		Titles:            []string{"title"},
		Bodies:            []string{"body"},
		Ratings:           []int{1},
		Datetimes:         []time.Time{now},
		Rating1Percentage: 10,
		Rating2Percentage: 0,
	}
	ReviewCountsModel, err = r.FindOrNewReviewCount(reviews)
	assert.Nil(t, err)
	assert.Equal(t, "app", ReviewCountsModel.AppName)
	assert.Equal(t, "store", ReviewCountsModel.Store)
}

func TestFindLastReviewCount(t *testing.T) {
	r := NewReviewsRepository()
	assert.NotNil(t, r)
	now := time.Now()

	reviews := Reviews{
		AppName:           "app",
		Store:             "store",
		Usernames:         []string{"username"},
		Titles:            []string{"title"},
		Bodies:            []string{"body"},
		Ratings:           []int{1},
		Datetimes:         []time.Time{now},
		Rating1Percentage: 20,
		Rating2Percentage: 0,
	}

	reviewCount, err := r.FindLastReviewCount(reviews)
	assert.Nil(t, err)
	assert.Equal(t, "app", reviewCount.AppName)
	assert.Equal(t, "store", reviewCount.Store)
	assert.Equal(t, 10, reviewCount.Rating1Percentage)
	assert.Equal(t, 0, reviewCount.Rating2Percentage)
	assert.Equal(t, 0, reviewCount.Rating3Percentage)
	assert.Equal(t, 0, reviewCount.Rating4Percentage)
	assert.Equal(t, 0, reviewCount.Rating5Percentage)

	_, err = r.InsertReviewCount(reviews)
	assert.Nil(t, err)

	reviewCount, err = r.FindLastReviewCount(reviews)
	assert.Nil(t, err)
	assert.Equal(t, "app", reviewCount.AppName)
	assert.Equal(t, "store", reviewCount.Store)
	assert.Equal(t, 20, reviewCount.Rating1Percentage)
	assert.Equal(t, 0, reviewCount.Rating2Percentage)
	assert.Equal(t, 0, reviewCount.Rating3Percentage)
	assert.Equal(t, 0, reviewCount.Rating4Percentage)
	assert.Equal(t, 0, reviewCount.Rating5Percentage)
}
func TestFindOrNewReviews(t *testing.T) {
	r := NewReviewsRepository()
	assert.NotNil(t, r)
	now := time.Now()
	newReviews, err := r.FindOrNewReviews(Reviews{})
	assert.Nil(t, err)
	assert.Equal(t, 0, len(newReviews))
	reviews := Reviews{
		AppName:           "app",
		Store:             "store",
		Usernames:         []string{"username"},
		Titles:            []string{"title"},
		Bodies:            []string{"body"},
		Ratings:           []int{1},
		Datetimes:         []time.Time{now},
		Rating1Percentage: 20,
		Rating2Percentage: 0,
	}

	newReviews, err = r.FindOrNewReviews(reviews)
	assert.Nil(t, err)
	assert.Equal(t, 1, len(newReviews))

}
