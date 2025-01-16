package main

import (
	"github.com/joho/godotenv"
	"log"
	"reco-test-task/internal"
	"reco-test-task/internal/common"
	"time"
)

const (
	FIRST_PERIOD  = 30 * time.Second
	SECOND_PERIOD = 5 * time.Minute
)

// Initialize environment variables from .env file
func initEnv() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}
}
func main() {
	initEnv()
	conf := common.NewConfig()

	periodicExtractor := internal.NewPeriodicExtractor(conf)
	go periodicExtractor.Start(FIRST_PERIOD)
	go periodicExtractor.Start(SECOND_PERIOD)

	select {}
}
