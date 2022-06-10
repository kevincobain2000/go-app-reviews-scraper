package services

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/araddon/dateparse"
	"github.com/headzoo/surf"
	"github.com/headzoo/surf/agent"
	"github.com/headzoo/surf/browser"
	"github.com/headzoo/surf/jar"
)

// getBrowser returns a new browser instance
// with a custom User-Agent
// and a cookie jar
func (s *SurfAppStore) getBrowser() *browser.Browser {
	bow := surf.NewBrowser()
	bow.SetUserAgent(agent.Chrome())
	bow.SetAttributes(browser.AttributeMap{
		browser.SendReferer:         surf.DefaultSendReferer,
		browser.MetaRefreshHandling: surf.DefaultMetaRefreshHandling,
		browser.FollowRedirects:     surf.DefaultFollowRedirects,
	})
	bow.SetCookieJar(jar.NewMemoryCookies())
	return bow
}

const (
	// ratingReviewUserClass is the css class where the Total number of reviews on the left side of ★ bars are shown
	// Example: 278 is fetched
	// 276件の評価 ★★★★★　----
	//            ★★★★　　----
	//            ★★★　　　-----
	//            ★★　　　　-------
	//            ★　　　　　----------
	ratingsTotalCssClass = ".we-customer-ratings__count"

	// ratingsBarCssClass is the css class where the ★ bars are shown
	// Example: 17 is fetched as a percentage for 5 star rating
	//          and so on
	// 276件の評価 ★★★★★　----　style="width:17%"
	//            ★★★★　　----
	//            ★★★　　　-----
	//            ★★　　　　-------
	//            ★　　　　　----------
	ratingsBarCssClass = ".we-star-bar-graph__bar__foreground-bar"

	// ratingStarCssClass is the css class where the ★ stars are shown for each review
	// In order to check if the rating is a 5 stars ★★★★☆,
	// it will have a another class in same element 'we-star-rating-stars-4' as 'we-star-rating-stars we-star-rating-stars-4'
	// declared as ratingStarCssClass below
	//
	// Example: 3 star rating is to be fetched and converted to int
	// ★★★☆☆
	// John Doe 2021/12/02
	// Great App
	// Hi, this is a great application..
	ratingStarCssClass = ".we-star-rating-stars"

	// These are the subclasses of the main class ratingStarCssClass
	// it is the sub class of the main css class
	// Example: a rating is 4 stars when element has class 'we-star-rating-stars we-star-rating-stars-4'
	rating1StarCssClassName = "we-star-rating-stars-1"
	rating2StarCssClassName = "we-star-rating-stars-2"
	rating3StarCssClassName = "we-star-rating-stars-3"
	rating4StarCssClassName = "we-star-rating-stars-4"
	rating5StarCssClassName = "we-star-rating-stars-5"

	// ratingReviewDateClass is the css class where the review date is shown
	// Example: "2021/12/02" is fetched and converted into date
	// ★★★☆☆
	// John Doe 2021/12/02
	// Great App
	// Hi, this is a great application..
	ratingReviewDateClass = ".we-customer-review__date"

	// ratingReviewUserClass is the css class where the username of the review is shown
	// Example: "John Doe" is fetched and trimmed
	// ★★★☆☆
	// John Doe 2021/12/02
	// Great App
	// Hi, this is a great application..
	ratingReviewUserClass = ".we-customer-review__user"

	// ratingReviewTitleClass is the css class where the title of the review is shown
	// Example: "Great App" is fetched and trimmed
	// ★★★☆☆
	// John Doe 2021/12/02
	// Great App
	// Hi, this is a great application..
	ratingReviewTitleClass = ".we-customer-review__title"

	// ratingReviewBodyClass is the css class where the body of the review is shown
	// Example:  "Hi, this is a great application.."" is fetched. This is not trimmed
	// ★★★☆☆
	// John Doe 2021/12/02
	// Great App
	// Hi, this is a great application..
	ratingReviewBodyClass = ".we-customer-review__body"
)

type SurfAppStore struct {
}

func NewSurfAppStore() *SurfAppStore {
	return &SurfAppStore{}
}

// Surf reviews for a given app
// and returns a Reviews struct
func (s *SurfAppStore) Surf(url string) (Reviews, error) {
	reviews := Reviews{}
	bow := s.getBrowser()

	var err error

	err = bow.Open(url)
	if err != nil {
		return reviews, fmt.Errorf("[error] unable to open the URL")
	}

	// Setting the mutable state of the struct pointer
	// Following procedure can be done in any order
	err = s.setRatingTotal(&reviews, bow)
	if err != nil {
		return reviews, err
	}
	err = s.setRatingsPercentage(&reviews, bow)
	if err != nil {
		return reviews, err
	}
	err = s.setRatings(&reviews, bow)
	if err != nil {
		return reviews, err
	}
	err = s.setReviewDate(&reviews, bow)
	if err != nil {
		return reviews, err
	}
	err = s.setReviewUser(&reviews, bow)
	if err != nil {
		return reviews, err
	}
	err = s.setReviewTitle(&reviews, bow)
	if err != nil {
		return reviews, err
	}
	err = s.setReviewBody(&reviews, bow)
	if err != nil {
		return reviews, err
	}

	// Finally verify the fetched reviews count is equal to the total reviews count
	// This is to verify that all classes were found and fetched in order
	err = s.verifyReviews(&reviews, bow)
	if err != nil {
		return reviews, err
	}

	return reviews, nil
}

// verifyReviews checks if the reviews surfed for all items equally
func (s *SurfAppStore) verifyReviews(reviews *Reviews, bow *browser.Browser) error {
	// check if all data was fetched successfully
	if len(reviews.Usernames) != len(reviews.Body) ||
		len(reviews.Usernames) != len(reviews.Titles) ||
		len(reviews.Usernames) != len(reviews.Datetimes) ||
		len(reviews.Usernames) != len(reviews.Ratings) {
		return fmt.Errorf("[error] fetched counts do not match up. All counts must be equal")
	}
	return nil
}

// setRatingTotal sets the total number of reviews
// for the app
// and returns an error if it fails
// or if the total is not found
// or if the total is not a number or float
// This is the total number of reviews shown on the page E.g 2.5 M review, or 200 reviews or 200 万 reviews
func (s *SurfAppStore) setRatingTotal(reviews *Reviews, bow *browser.Browser) error {
	hasError := false
	bow.Dom().Find(ratingsTotalCssClass).Each(func(idx int, s *goquery.Selection) {
		re := regexp.MustCompile(`[-]?\d[\d,]*[\.]?[\d{2}]*`)
		count := re.FindAllString(s.Text(), -1)
		if len(count) == 0 {
			hasError = true
		}

		total, err := strconv.ParseFloat(count[0], 64)
		if err != nil {
			hasError = true
		}

		//Check if in string, dirty code to interpret M as million
		if strings.Contains(s.Text(), "M ") {
			total = total * 1000000
		}
		if strings.Contains(s.Text(), "万") {
			total = total * 10000
		}
		// convert float to int
		reviews.Total = int(total)
	})
	if hasError {
		return fmt.Errorf("[error] unable to fetch total")
	}
	return nil
}

// setReviewDate sets the date of the review
// and returns an error if it fails
// or if the date is not found
// or if the date is not a valid date
// This is the date from the datetime attribute of the css class where date of review for a review is shown
func (s *SurfAppStore) setReviewDate(reviews *Reviews, bow *browser.Browser) error {
	hasError := false
	bow.Dom().Find(ratingReviewDateClass).Each(func(_ int, s *goquery.Selection) {
		datetimeStr, has := s.Attr("datetime")
		if !has {
			hasError = true
		}
		if has {
			datetime, err := dateparse.ParseAny(datetimeStr)
			if err != nil {
				hasError = true
			} else {
				reviews.Datetimes = append(reviews.Datetimes, datetime)
			}
		}
	})
	if hasError {
		return fmt.Errorf("[error] unable to parse datetime")
	}
	return nil
}

// setReviewUser sets the username of the review
// and returns an error if it fails
// or if the username is not found
// This is the username from the css class where username of review for a review is shown
func (s *SurfAppStore) setReviewUser(reviews *Reviews, bow *browser.Browser) error {
	bow.Dom().Find(ratingReviewUserClass).Each(func(_ int, s *goquery.Selection) {
		reviews.Usernames = append(reviews.Usernames, strings.TrimSpace(s.Text()))
	})
	return nil
}

// setReviewTitle sets the title of the review
// and returns an error if it fails
// or if the title is not found
// This is the title from the css class where title of review for a review is shown
func (s *SurfAppStore) setReviewTitle(reviews *Reviews, bow *browser.Browser) error {
	bow.Dom().Find(ratingReviewTitleClass).Each(func(_ int, s *goquery.Selection) {
		reviews.Titles = append(reviews.Titles, strings.TrimSpace(s.Text()))
	})
	return nil
}

// setReviewBody sets the body of the review
// and returns an error if it fails
// or if the body is not found
// This is the body from the css class where body of review for a review is shown
func (s *SurfAppStore) setReviewBody(reviews *Reviews, bow *browser.Browser) error {
	bow.Dom().Find(ratingReviewBodyClass).Each(func(_ int, s *goquery.Selection) {
		reviews.Body = append(reviews.Body, strings.TrimSpace(s.Text()))
	})
	return nil
}

// setRatings sets the rating of the review
// and returns an error if it fails
// or if the rating is not found
// This is the rating from the css class where rating of review for a review is shown
// There are many elements for this css class, so we loop through them in order
func (s *SurfAppStore) setRatings(reviews *Reviews, bow *browser.Browser) error {
	bow.Dom().Find(ratingStarCssClass).Each(func(_ int, s *goquery.Selection) {
		rating := 0
		if s.HasClass(rating1StarCssClassName) {
			rating = 5
		}
		if s.HasClass(rating2StarCssClassName) {
			rating = 4
		}
		if s.HasClass(rating3StarCssClassName) {
			rating = 3
		}
		if s.HasClass(rating4StarCssClassName) {
			rating = 2
		}
		if s.HasClass(rating5StarCssClassName) {
			rating = 1
		}
		reviews.Ratings = append(reviews.Ratings, rating)
	})
	return nil
}

// setRatingsPercentage sets the rating percentage of the review
// and returns an error if it fails
// or if the rating percentage is not found
// This is the rating percentage from the style attribute of the element which shows the bar of the rating
func (s *SurfAppStore) setRatingsPercentage(reviews *Reviews, bow *browser.Browser) error {
	hasError := false
	bow.Dom().Find(ratingsBarCssClass).Each(func(idx int, s *goquery.Selection) {
		style, has := s.Attr("style")
		if !has {
			hasError = true
		}
		// extract width: 17% from the style attribute
		re := regexp.MustCompile(`width: \d+\%`)
		width := re.FindAllString(style, -1)
		reP := regexp.MustCompile(`\d+`)
		if len(width) == 0 {
			hasError = true
		}
		// extract just the number 17 from the matched string "width: 17%"
		percentage := reP.FindAllString(width[0], -1)
		if len(percentage) == 0 {
			hasError = true
		}
		//convert to int
		percentageInt, err := strconv.Atoi(percentage[0])
		if err != nil {
			hasError = true
		}
		switch idx {
		// is first bar element 5 stars
		case 0:
			reviews.Rating5Percentage = percentageInt
		// is second bar element
		case 1:
			reviews.Rating4Percentage = percentageInt
		// and so on...
		case 2:
			reviews.Rating3Percentage = percentageInt
		case 3:
			reviews.Rating2Percentage = percentageInt
		case 4:
			reviews.Rating1Percentage = percentageInt
		}
	})
	if hasError {
		return fmt.Errorf("[error] unable to fetch percentage")
	}
	return nil
}
