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
		{
			// no protocol, bad urls following
			urlStr:    "play.google.com/store/apps/details?id=com.king.candycrushsaga&hl=en&gl=US",
			storeWant: "",
			errWant:   fmt.Errorf(""),
		},
		{
			urlStr:    "apps.apple.com/us/app/candy-crush-saga/id553834731?see-all=reviews",
			storeWant: "",
			errWant:   fmt.Errorf(""),
		},
		{
			urlStr:    "example.com/apps/candy-crush-saga",
			storeWant: "",
			errWant:   fmt.Errorf(""),
		},
		{
			urlStr:    "bad/url/apps/candy-crush-saga",
			storeWant: "",
			errWant:   fmt.Errorf(""),
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
