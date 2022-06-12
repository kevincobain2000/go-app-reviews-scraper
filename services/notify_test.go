package services

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func init() {
	Setup()
}

func TestNotify(t *testing.T) {
	nn := NewNotify()
	repo := NewReviewsRepository()
	reviews := Reviews{
		AppName: "test",
		Store:   "test",
		Usernames: []string{
			"test",
		},
		Titles: []string{
			"test",
		},
		Bodies: []string{
			"test",
		},
		Ratings: []int{
			5,
		},
		Datetimes: []time.Time{
			time.Now(),
		},
	}
	newReviews, err := repo.FindOrNewReviews(reviews)
	assert.Nil(t, err)
	err = nn.NotifyNewReviews(newReviews)
	assert.Nil(t, err)

	_, err = repo.FindLastReviewCount(reviews)
	assert.Nil(t, err)
	_, err = repo.FindOrNewReviewCount(reviews)
	assert.Nil(t, err)
}
