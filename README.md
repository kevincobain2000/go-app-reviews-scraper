[![codecov](https://codecov.io/gh/kevincobain2000/go-app-reviews-scraper/branch/master/graph/badge.svg)](https://codecov.io/gh/kevincobain2000/go-app-reviews-scraper)


<p align="center">
  <a href="https://github.com/kevincobain2000/go-app-reviews-scraper">
    <img alt="go-app-reviews-scraper" src="logo.png" width="360">
  </a>
</p>

<h3 align="center">Scrape Reviews and Ratings.<br>Notify when new review arrives.</h3>

<p align="center">
  Monitor your app store reviews and get notified when new reviews are published.
  <br>
  Apple app store reviews command line API for iOS and Google Play apps.
  <br>
  Written in Go :heart:

</p>

**Supports:** Apple App store and Google Play Store

**Command line:** Arch free binary to run as scheduler on any platform.

**Headless:** Uses headless browser to run without selenium or chromium drivers.

**Proxy support:** Works behind a proxy.

**Notifications:** Supports multiple notification channels Microsoft Teams or console output.

**Dependencies:** None. Works with sqlite on disk or in memory.

### Notification sample on new review


<h3 align="center">
    Microsoft Teams
</h3>

<p align="center">
  <img src="screenshot1.png" alt="teams">
</p>

<h3 align="center">
   Terminal
</h3>

<p align="center">
  <img src="screenshot2.png" alt="teams">
</p>

### Notification sample on new rating

<h3 align="center">
    Microsoft Teams
</h3>

<p align="center">
  <img src="screenshot3.png" alt="teams">
</p>

<h3 align="center">
   Terminal
</h3>

<p align="center">
  <img src="screenshot4.png" alt="teams">
</p>



## Installation

```sh
go install github.com/kevincobain2000/go-app-reviews-scraper@latest
```

### From Command Line:

```sh
cp .env.local .env
go-app-reviews-scraper -migrate
ENV_PATH=./.env go-app-reviews-scraper -app-name="candy-crush" -reviews-url="https://apps.apple.com/us/app/candy-crush-saga/id553834731?see-all=reviews"
ENV_PATH=./.env go-app-reviews-scraper -app-name="candy-crush" -reviews-url="https://play.google.com/store/apps/details?id=com.king.candycrushsaga&hl=en&gl=US"
```

--

### Command Line Params Help:

```sh
go-app-reviews-scraper -h
  -app-name string
    	Description: Give a unique app name. Example: candy-crush
  -migrate
    	Description: Run DB migration
  -reviews-url string

    	Description: Apple's link to reviews page. Example: https://apps.apple.com/us/app/candy-crush-saga/id553834731?see-all=reviews
    	Description: Google's link reviews page. Example: https://play.google.com/store/apps/details?id=com.king.candycrushsaga&hl=en&gl=US
```

### CHANGE LOG

- v1.0 - Initial release includes iOS App store reviews scraper and notification to MS Teams.
- v1.0 - Support for Google Play Store.

### ROADMAP

- v1.2 - Work in progress. Send notification via Email.
- v1.3 - Work in progress. Send notification via Slack.
