package services

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func init() {
	Setup()
}

const (
	candyCrushURL = "https://apps.apple.com/us/app/candy-crush-saga/id553834731?see-all=reviews"
)

func TestGetBrowser(t *testing.T) {
	s := NewSurfAppStore()
	browser := s.getBrowser()
	assert.NotNil(t, browser)
}

func TestSurf(t *testing.T) {
	s := NewSurfAppStore()
	// will do actual test
	reviews, err := s.SurfAppleStore(candyCrushURL)
	assert.Nil(t, err)
	assert.NotNil(t, reviews)
	assert.LessOrEqual(t, 10000, reviews.Total)         // at least more than this many reviews are there
	assert.LessOrEqual(t, 1, reviews.Rating1Percentage) // at least more than this many reviews are there
	assert.LessOrEqual(t, 1, reviews.Rating2Percentage)
	assert.LessOrEqual(t, 1, reviews.Rating3Percentage)
	assert.LessOrEqual(t, 1, reviews.Rating4Percentage)
	assert.LessOrEqual(t, 10, reviews.Rating5Percentage) // at least more than this many reviews are there, test will fail if users bombard the ratings
}
