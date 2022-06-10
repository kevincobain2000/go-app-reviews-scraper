package services

import (
	"fmt"
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

func (ut *Utils) GetStoreFromURL(urlStr string) (string, error) {
	u, err := url.Parse(urlStr)
	if err != nil {
		return "", err
	}

	if strings.Contains(u.Host, AppStoreHost) {
		return StoreIOS, nil
	}
	if strings.Contains(u.Host, PlayStoreHost) {
		return StoreAndroid, nil
	}
	return "", fmt.Errorf("[error] Unable to fetch store from the given Review URL %s", urlStr)
}
