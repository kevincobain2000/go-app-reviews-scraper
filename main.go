package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/kevincobain2000/go-app-reviews-scraper/app"
	"github.com/kevincobain2000/go-app-reviews-scraper/services"
)

var (
	appName    = flag.String("app-name", "", "Description: Give a unique app name. Example: candy-crush")
	reviewsURL = flag.String("reviews-url", "", "Description: Link to all reviews page. Example: https://apps.apple.com/us/app/candy-crush-saga/id553834731?see-all=reviews")
	migrate    = flag.Bool("migrate", false, "Description: Run DB migration")
)

// main execution starts here for the command line interface
// sets env
// runs migrations
// scrapes reviews
// saves reviews to DB
// prints out the reviews
// notifies on MS teams
// exits with 0, 1 on fatal errors
func main() {
	// set env and flags
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	app.SetEnv()
	flag.Parse()

	// print cli args
	flag.VisitAll(func(f *flag.Flag) {
		log.Printf("[info] %s: %s\n", f.Name, f.Value)
	})

	// run db migration if flag is set. This is for first time setup or on version up
	// Migrate doesn't delete your old data from DB
	if *migrate {
		services.AutoMigrate()
		log.Println("[info] DB migration ran")
		return
	}

	// Required args check
	if *appName == "" || *reviewsURL == "" {
		log.Fatal("[fatal] Missing required flags. See -h for help.")
	}

	// prepare services and repositories
	uu := services.NewUtils()
	ss := services.NewSurfAppStore()
	nn := services.NewNotify()
	repo := services.NewReviewsRepository()
	// end prepare services and repositories

	// parse reviews url and get store
	store, err := uu.GetStoreFromURL(*reviewsURL)
	if err != nil {
		log.Fatal(err)
	}

	// create reviews object where reviews are going to be stored
	reviews := services.Reviews{
		AppName: *appName,
		Store:   store,
	}

	fmt.Println("[info] Started browser to scrape")
	if store == services.StoreIOS {
		reviews, err = ss.SurfAppleStore(*reviewsURL)
	}
	if store == services.StoreAndroid {
		reviews, err = ss.SurfAppleStore(*reviewsURL)
	}
	if err != nil {
		log.Fatal(err)
	}
	if reviews.Total == 0 {
		log.Fatal("[fatal] No reviews found, or something went wrong during fetching")
	}

	// Start procedure to update database

	// 1) from the scraped reviews, check if the reviews are in DB
	//    if not then insert the newly scraped reviews
	//    return is the new reviews that are inserted
	newReviews, err := repo.FindOrNewReviews(reviews)
	if err != nil {
		log.Fatal(err)
	}

	// 2) Get the last review count summary from DB
	lastReviewCount, err := repo.FindLastReviewCount(reviews)
	if err != nil {
		log.Fatal(err)
	}

	// 3) Only insert the newly scaped reviews summary if there is no difference
	//    when there is no diff then the last review summary is returned
	currentReviewCount, err := repo.FindOrNewReviewCount(reviews)
	if err != nil {
		log.Fatal(err)
	}

	// 4) Check if a new review count summary is created or just using previous one
	//   Use it for the notification purpose
	if currentReviewCount.ID != 0 && currentReviewCount.ID != lastReviewCount.ID {
		// we have new rating summary since last scraped
		err := nn.NotifyReviewCount(currentReviewCount, lastReviewCount)
		if err != nil {
			log.Fatal(err)
		}
	}

	// 5) Notify on new reviews if any
	err = nn.NotifyNewReviews(newReviews)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("[info] Finished!")
}
