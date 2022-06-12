package services

import (
	reviewsService "github.com/n0madic/google-play-scraper/pkg/reviews"
)

// SurfGoogleStore is a SurfGoogleStore
type SurfGoogleStore struct {
	// Number is the number of reviews to scrape
	Number int
}

// NewSurfGoogleStore creates a new SurfGoogleStore
func NewSurfGoogleStore(number int) *SurfGoogleStore {
	return &SurfGoogleStore{
		Number: number,
	}
}

// Surf Google Store
// Uses a library to scrape reviews from Google Play
func (s *SurfGoogleStore) Surf(urlStr string) (Reviews, error) {
	reviews := Reviews{}
	uu := NewUtils()
	id, language, err := uu.GetAppInfoGoogle(urlStr)
	if err != nil {
		return reviews, err
	}

	r := reviewsService.New(id, reviewsService.Options{
		Number:   s.Number,
		Language: language,
	})
	err = r.Run()
	if err != nil {
		return reviews, err
	}

	ratings1Count := 0
	ratings2Count := 0
	ratings3Count := 0
	ratings4Count := 0
	ratings5Count := 0
	for _, review := range r.Results {
		// fmt.Println(review.Reviewer, review.Timestamp, review.Score, review.Text)
		reviews.Titles = append(reviews.Titles, "Google Play store") // there is no title in Google, so settings it as default
		reviews.Usernames = append(reviews.Usernames, review.Reviewer)
		reviews.Datetimes = append(reviews.Datetimes, review.Timestamp)
		reviews.Bodies = append(reviews.Bodies, review.Text)
		reviews.Ratings = append(reviews.Ratings, review.Score)

		reviews.Total++
		switch review.Score {
		case 1:
			ratings1Count++
		case 2:
			ratings2Count++
		case 3:
			ratings3Count++
		case 4:
			ratings4Count++
		case 5:
			ratings5Count++
		}
	}

	reviews.Rating1Percentage = uu.CalculateRoundedPercentage(ratings1Count, reviews.Total)
	reviews.Rating2Percentage = uu.CalculateRoundedPercentage(ratings2Count, reviews.Total)
	reviews.Rating3Percentage = uu.CalculateRoundedPercentage(ratings3Count, reviews.Total)
	reviews.Rating4Percentage = uu.CalculateRoundedPercentage(ratings4Count, reviews.Total)
	reviews.Rating5Percentage = uu.CalculateRoundedPercentage(ratings5Count, reviews.Total)

	if err := VerifyReviews(&reviews); err != nil {
		return reviews, err
	}

	return reviews, nil
}
