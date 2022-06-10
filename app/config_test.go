package app

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var conf *Config

func init() {
	Setup()
	conf = NewConfig()
}

func TestConfig(t *testing.T) {
	assert.NotNil(t, conf.AppConfig)
	assert.NotNil(t, conf.DBConfig)
}
