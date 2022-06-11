package services

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func init() {
	Setup()
}

func TestVerifyReviews(t *testing.T) {
	reviews := Reviews{
		AppName: "test",
		Store:   "test",
		Usernames: []string{
			"test",
			"test2",
		},
		Titles: []string{
			"test",
			"test2",
		},
		Body: []string{
			"test",
			"test2",
		},
		Ratings: []int{
			5,
			3,
		},
		Datetimes: []time.Time{
			time.Now(),
			time.Now(),
		},
	}
	reviews.Total = 2
	reviews.Rating1Percentage = 2
	err := VerifyReviews(&reviews)
	assert.Nil(t, err)
}

func TestVerifyReviewsError(t *testing.T) {
	reviews := Reviews{
		AppName: "test",
		Store:   "test",
		Usernames: []string{
			"test",
			"test2",
		},
		Titles: []string{
			"test",
		},
		Body: []string{
			"test",
			"test2",
		},
		Ratings: []int{
			5,
		},
		Datetimes: []time.Time{
			time.Now(),
		},
	}
	err := VerifyReviews(&reviews)
	assert.NotNil(t, err)
}
