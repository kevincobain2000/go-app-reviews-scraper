package services

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetStoreFromURL(t *testing.T) {
	uu := NewUtils()
	tests := []struct {
		urlStr    string
		storeWant string
		errWant   error
	}{
		{
			urlStr:    "https://apps.apple.com/us/app/candy-crush-saga/id553834731?see-all=reviews",
			storeWant: "ios",
			errWant:   nil,
		},

		{
			urlStr:    "https://play.google.com/store/apps/details?id=com.king.candycrushsaga&hl=en&gl=US",
			storeWant: "android",
			errWant:   nil,
		},
		// no protocol, bad urls following, not starts with the required domain
		{
			urlStr:    "https://sumdomain.apps.apple.com/us/app/candy-crush-saga/id553834731?see-all=reviews",
			storeWant: "",
			errWant:   fmt.Errorf("error"),
		},

		{
			urlStr:    "https://subdomain.play.google.com/store/apps/details?id=com.king.candycrushsaga&hl=en&gl=US",
			storeWant: "",
			errWant:   fmt.Errorf("error"),
		},
		{
			urlStr:    "play.google.com/store/apps/details?id=com.king.candycrushsaga&hl=en&gl=US",
			storeWant: "",
			errWant:   fmt.Errorf("error"),
		},
		{
			urlStr:    "apps.apple.com/us/app/candy-crush-saga/id553834731?see-all=reviews",
			storeWant: "",
			errWant:   fmt.Errorf("error"),
		},
		{
			urlStr:    "example.com/apps/candy-crush-saga",
			storeWant: "",
			errWant:   fmt.Errorf("error"),
		},
		{
			urlStr:    "bad/url/apps/candy-crush-saga",
			storeWant: "",
			errWant:   fmt.Errorf("error"),
		},
	}
	for _, test := range tests {
		t.Run(test.urlStr, func(t *testing.T) {
			store, err := uu.GetStoreFromURL(test.urlStr)
			assert.Equal(t, test.storeWant, store)
			if test.errWant != nil {
				assert.NotNil(t, err)
			}
		})
	}
}

func TestGetAppInfoGoogle(t *testing.T) {
	uu := NewUtils()
	tests := []struct {
		urlStr       string
		idWant       string
		languageWant string
		errWant      error
	}{
		{
			urlStr:       "https://play.google.com/store/apps/details?id=com.king.candycrushsaga&hl=en&gl=US",
			idWant:       "com.king.candycrushsaga",
			languageWant: "en",
			errWant:      nil,
		},
		{
			urlStr:       "https://play.google.com/store/apps/details?hl=en&id=com.king.candycrushsaga&gl=US",
			idWant:       "com.king.candycrushsaga",
			languageWant: "en",
			errWant:      nil,
		},
		// no protocol, bad urls following
		{
			urlStr:       "https://play.google.com/store/apps/details?hl=en&id=com.king.candycrushsaga&gl=US",
			idWant:       "com.king.candycrushsaga",
			languageWant: "en",
			errWant:      nil,
		},
		{
			urlStr:       "example.com/apps/candy-crush-saga",
			idWant:       "",
			languageWant: "",
			errWant:      fmt.Errorf("error"),
		},
		{
			urlStr:  "bad/url/apps/candy-crush-saga?id=com.king.candycrushsaga&hl=en&gl=US",
			idWant:  "",
			errWant: fmt.Errorf("error"),
		},
		{
			urlStr:  "invalid.com/url/apps/candy-crush-saga?id=com.king.candycrushsaga&hl=en&gl=US",
			idWant:  "",
			errWant: fmt.Errorf("error"),
		},
	}
	for _, test := range tests {
		t.Run(test.urlStr, func(t *testing.T) {
			id, language, err := uu.GetAppInfoGoogle(test.urlStr)
			assert.Equal(t, test.idWant, id)
			assert.Equal(t, test.languageWant, language)
			if test.errWant != nil {
				assert.NotNil(t, err)
			}
		})
	}
}

func TestCalculateRoundedPercentage(t *testing.T) {
	uu := NewUtils()
	tests := []struct {
		count int
		total int
		want  int
	}{
		{
			count: 127,
			total: 999,
			want:  13,
		},
		{
			count: 123,
			total: 999,
			want:  12,
		},
		{
			count: 123,
			total: 0,
			want:  0,
		},
		{
			count: 0,
			total: 0,
			want:  0,
		},
		{
			count: 0,
			total: 999,
			want:  0,
		},
	}
	for _, test := range tests {
		t.Run("percentage test", func(t *testing.T) {
			percentage := uu.CalculateRoundedPercentage(test.count, test.total)
			assert.Equal(t, test.want, percentage)
		})
	}
}
func TestAverageRating(t *testing.T) {
	uu := NewUtils()

	tests := []struct {
		reviewCounts  ReviewCountsModel
		wantAvgRating float64
	}{
		{
			reviewCounts: ReviewCountsModel{
				Rating1Percentage: 123,
				Rating2Percentage: 123,
				Rating3Percentage: 123,
				Rating4Percentage: 123,
				Rating5Percentage: 123,
			},
			wantAvgRating: 3,
		},
		{
			reviewCounts: ReviewCountsModel{
				Rating1Percentage: 10,
				Rating2Percentage: 20,
				Rating3Percentage: 30,
				Rating4Percentage: 40,
				Rating5Percentage: 50,
			},
			wantAvgRating: 3.6666666666666665,
		},
		// irregular
		{
			reviewCounts: ReviewCountsModel{
				Rating1Percentage: 0,
				Rating2Percentage: -1,
				Rating3Percentage: 30,
				Rating4Percentage: 40,
				Rating5Percentage: 50,
			},
			wantAvgRating: 4.184873949579832, // still not go to error
		},
		{
			reviewCounts: ReviewCountsModel{
				Rating1Percentage: 0,
				Rating2Percentage: 0,
				Rating3Percentage: 0,
				Rating4Percentage: 0,
				Rating5Percentage: 0,
			},
			wantAvgRating: 0, // still not go to error
		},
		{
			reviewCounts: ReviewCountsModel{
				Rating1Percentage: 0,
				Rating2Percentage: -1,
				Rating3Percentage: 0,
				Rating4Percentage: 0,
				Rating5Percentage: 0,
			},
			wantAvgRating: 0, // still not go to error
		},
	}
	for _, test := range tests {
		t.Run("percentage test", func(t *testing.T) {
			avg := uu.AverageRating(test.reviewCounts)
			assert.Equal(t, test.wantAvgRating, avg)
		})
	}
}
