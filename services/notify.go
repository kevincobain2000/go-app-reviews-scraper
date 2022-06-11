package services

import (
	"fmt"
	"log"
	"strings"

	md "github.com/JohannesKaufmann/html-to-markdown"
	"github.com/kevincobain2000/go-app-reviews-scraper/app"
	gmt "github.com/kevincobain2000/go-msteams/src"
)

type Notify struct {
}

func NewNotify() *Notify {
	return &Notify{}
}

// NotifyNewReviews sends a notification to a list of email addresses
// sends notification on microsoft teams
// also stdout the message to console in markdown of the message (html)
// when the env for MS Teams hook or EMAIL addresses are not sent the notifications are not sent
func (n *Notify) NotifyNewReviews(reviews []ReviewModel) error {
	for _, review := range reviews {
		// send message to MS Teams
		title := "You have a new review!"
		subtitle := "Store (" + review.Store + ")"
		subject := "App (" + review.AppName + ")"
		color := ""
		proxy := ""
		message := ""
		message += "<h2>" + review.Title + "</h2>" + "<br>"
		message += "<h3>@" + review.Username + "</h3>" + "<br>"
		message += "<h4>" + review.RatedAt.Format("02-Jan-2006") + "</h4>" + "<br>"
		message += "<h5>" + "Rating " + strings.Repeat("★", review.Rating) + strings.Repeat("☆", 5-review.Rating) + "</h5>" + "<br>"
		message += "<p>" + review.Body + "<p>" + "<br>"
		// For ascii output
		converter := md.NewConverter("", true, nil)
		markdown, err := converter.ConvertString(message)
		if err != nil {
			return err
		}
		log.Println("[info] Printing to console")
		for _, line := range strings.Split(markdown, "\n") {
			if strings.TrimSpace(line) != "" {
				fmt.Println(line)
			}
		}

		c := app.NewConfig()

		// For MS Teams
		if c.AppConfig.MSTeamsHookURL != "" {
			log.Println("[info] Sending to MS Teams")
			if err := gmt.Send(title, subtitle, subject, color, message, c.AppConfig.MSTeamsHookURL, proxy); err != nil {
				return err
			}
		}
	}
	return nil
}

// NotifyReviewCount sends a notification to a list of email addresses
// sends notification on microsoft teams
// also stdout the message to console in markdown of the message (html)
// when the env for MS Teams hook or EMAIL addresses are not sent the notifications are not sent
func (n *Notify) NotifyReviewCount(currentReviewCount, lastReviewCount ReviewCountsModel) error {
	uu := NewUtils()
	title := "You have a new rating!"
	subtitle := "Store (" + currentReviewCount.Store + ")"
	subject := "App (" + currentReviewCount.AppName + ")"
	color := ""
	proxy := ""
	message := ""

	message += "<b>Now </b>" + currentReviewCount.CreatedAt.Format("02-Jan-2006") + "<br>"
	message += fmt.Sprintf("Total reviews: %d", currentReviewCount.Total) + "<br>"
	message += fmt.Sprintf("Average rating: %.2f", uu.AverageRating(currentReviewCount)) + "<br>"
	message += fmt.Sprintf("★★★★★: %d", currentReviewCount.Rating5Percentage) + "%" + "<br>"
	message += fmt.Sprintf("★★★★☆: %d", currentReviewCount.Rating4Percentage) + "%" + "<br>"
	message += fmt.Sprintf("★★★☆☆: %d", currentReviewCount.Rating3Percentage) + "%" + "<br>"
	message += fmt.Sprintf("★★☆☆☆: %d", currentReviewCount.Rating2Percentage) + "%" + "<br>"
	message += fmt.Sprintf("★☆☆☆☆: %d", currentReviewCount.Rating1Percentage) + "%" + "<br>"

	if lastReviewCount.Total > 0 {
		message += "<b>Before </b>" + lastReviewCount.CreatedAt.Format("02-Jan-2006") + "<br>"
		message += fmt.Sprintf("Average rating: %.2f", uu.AverageRating(lastReviewCount)) + "<br>"
		message += fmt.Sprintf("Total reviews: %d", lastReviewCount.Total) + "<br>"
		message += fmt.Sprintf("★★★★★: %d", lastReviewCount.Rating5Percentage) + "%" + "<br>"
		message += fmt.Sprintf("★★★★☆: %d", lastReviewCount.Rating4Percentage) + "%" + "<br>"
		message += fmt.Sprintf("★★★☆☆: %d", lastReviewCount.Rating3Percentage) + "%" + "<br>"
		message += fmt.Sprintf("★★☆☆☆: %d", lastReviewCount.Rating2Percentage) + "%" + "<br>"
		message += fmt.Sprintf("★☆☆☆☆: %d", lastReviewCount.Rating1Percentage) + "%" + "<br>"
	}

	// For ascii output
	converter := md.NewConverter("", true, nil)
	markdown, err := converter.ConvertString(message)
	if err != nil {
		return err
	}
	log.Println("[info] Printing to console")
	for _, line := range strings.Split(markdown, "\n") {
		if strings.TrimSpace(line) != "" {
			fmt.Println(line)
		}
	}

	c := app.NewConfig()
	// For MS Teams
	if c.AppConfig.MSTeamsHookURL != "" {
		log.Println("[info] Sending to MS Teams")
		if err := gmt.Send(title, subtitle, subject, color, message, c.AppConfig.MSTeamsHookURL, proxy); err != nil {
			return err
		}
	}
	return nil
}
