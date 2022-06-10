package app

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEnv(t *testing.T) {

	os.Setenv("ENV_PATH", ".env.not_exist")
	assertPanic(t, SetEnv)

	os.Setenv("ENV_PATH", ".env.testing")
	assert.Equal(t, "testing", os.Getenv("APP_ENV"))
}

func assertPanic(t *testing.T, f func()) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("The code did not panic")
		}
	}()
	f()
}
