package services

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func init() {
	Setup()
}

const (
	iosReviewsCandyCrushURL    = "https://apps.apple.com/us/app/candy-crush-saga/id553834731?see-all=reviews"
	googleReviewsCandyCrushURL = "https://play.google.com/store/apps/details?id=com.king.candycrushsaga&hl=en&gl=US"
)

func TestGetBrowser(t *testing.T) {
	s := NewSurfAppStore()
	browser := s.getBrowser()
	assert.NotNil(t, browser)
}

func TestSurfAppleStore(t *testing.T) {
	s := NewSurfAppStore()
	// will do actual test
	reviews, err := s.Surf(iosReviewsCandyCrushURL)
	assert.Nil(t, err)
	assert.NotNil(t, reviews)
	assert.LessOrEqual(t, 10000, reviews.Total)         // at least more than this many reviews are there
	assert.LessOrEqual(t, 1, reviews.Rating1Percentage) // at least more than this many reviews are there
}

func TestSurfGoogleStore(t *testing.T) {
	s := NewSurfGoogleStore(15) // extract a few for faster tests
	// will do actual test
	reviews, err := s.Surf(googleReviewsCandyCrushURL)
	assert.Nil(t, err)
	assert.NotNil(t, reviews)
	assert.LessOrEqual(t, 15, reviews.Total)            // at least more than this many reviews are there
	assert.LessOrEqual(t, 1, reviews.Rating1Percentage) // at least more than this many reviews are there
}
