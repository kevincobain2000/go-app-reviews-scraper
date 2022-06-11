package services

import (
	"fmt"
	"math"
	"net/url"
	"strings"
)

type Utils struct {
}

func NewUtils() *Utils {
	return &Utils{}
}

const (
	// StoreIOS is the store name for iOS
	StoreIOS = "ios"
	// AppStoreHost is the hostname for the App Store
	AppStoreHost = "apps.apple.com"

	// StoreAndroid is the store name for Android
	StoreAndroid = "android"
	// PlayStoreHost is the hostname for the Play Store
	PlayStoreHost = "play.google.com"
)

// GetStoreFromURL returns the store name from the given URL
// The store name is one of the following:
// - ios
// - android
// - ""
// - error
// The error is returned when the URL is not valid
// The error is returned when the URL is not a valid App Store URL
// The error is returned when the URL is not a valid Play Store URL
func (ut *Utils) GetStoreFromURL(urlStr string) (string, error) {
	u, err := url.Parse(urlStr)
	if err != nil {
		return "", err
	}

	if strings.HasPrefix(u.Host, AppStoreHost) {
		return StoreIOS, nil
	}
	if strings.HasPrefix(u.Host, PlayStoreHost) {
		return StoreAndroid, nil
	}
	return "", fmt.Errorf("[error] Unable to fetch store from the given Review URL %s", urlStr)
}

// GetAppInfoGoogle returns the app info from the given Google Play URL
// The app info is one of the following:
// - id, location, nil
// - "", "", error
//
// urlStr = play.google.com/store/apps/details?id=com.king.candycrushsaga&hl=en&gl=US
// id returned as com.king.candycrushsaga
// hl returned as hl
//
// The error is returned when the URL is not valid
// The error is returned when the URL is not a valid Google Play URL
func (ut *Utils) GetAppInfoGoogle(urlStr string) (string, string, error) {
	u, err := url.Parse(urlStr)
	if err != nil {
		return "", "", err
	}
	if !strings.HasPrefix(u.Host, PlayStoreHost) {
		return "", "", fmt.Errorf("[error] Unable to fetch store from the given Review URL %s", urlStr)
	}

	queries, err := url.ParseQuery(u.RawQuery)
	if err != nil {
		return "", "", err
	}
	if queries["id"] == nil {
		return "", "", fmt.Errorf("[error] Unable to fetch app ID from the given Review URL %s", urlStr)
	}
	if queries["hl"] == nil {
		return "", "", fmt.Errorf("[error] Unable to fetch hl language from Review URL %s", urlStr)
	}

	id := queries["id"][0]
	language := queries["hl"][0]
	return id, language, nil
}

// CalculateRoundedPercentage returns the rounded percentage
// E.g 127/999 x 100 = 12.71%, rounded to int is 13
// E.g 123/999 x 100 = 12.31%, rounded to int is 12
func (ut *Utils) CalculateRoundedPercentage(count, total int) int {
	percentage := int(math.Round(float64(count) / float64(total) * 100))
	if percentage <= 0 {
		return 0
	}
	return percentage
}

// AverageRating returns the average rating
func (ut *Utils) AverageRating(reviewCounts ReviewCountsModel) float64 {
	if reviewCounts.Rating1Percentage+
		reviewCounts.Rating2Percentage+
		reviewCounts.Rating3Percentage+
		reviewCounts.Rating4Percentage+
		reviewCounts.Rating5Percentage <= 0 {
		return 0.0
	}
	avg := (1*float64(reviewCounts.Rating1Percentage) +
		2*float64(reviewCounts.Rating2Percentage) +
		3*float64(reviewCounts.Rating3Percentage) +
		4*float64(reviewCounts.Rating4Percentage) +
		5*float64(reviewCounts.Rating5Percentage)) /
		(float64(reviewCounts.Rating1Percentage) +
			float64(reviewCounts.Rating2Percentage) +
			float64(reviewCounts.Rating3Percentage) +
			float64(reviewCounts.Rating4Percentage) +
			float64(reviewCounts.Rating5Percentage))
	return avg
}
