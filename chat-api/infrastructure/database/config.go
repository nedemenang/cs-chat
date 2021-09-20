package database

import (
	"os"
	"time"
)

type config struct {
	host     string
	database string
	port     string
	driver   string
	user     string
	password string
	uri      string

	ctxTimeout time.Duration
}

func newConfigMongoDB() *config {
	return &config{
		uri:        os.Getenv("MONGODB_URI"),
		database:   os.Getenv("MONGODB_DATABASE"),
		ctxTimeout: 60 * time.Second,
	}
}
