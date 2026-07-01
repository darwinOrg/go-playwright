package main

import (
	"log"

	"github.com/mxschmitt/playwright-go"
)

func main() {
	driver, err := playwright.NewDriver(&playwright.RunOptions{})
	if err != nil {
		log.Fatalf("could not start driver: %v", err)
	}
	if err = driver.DownloadDriver(); err != nil {
		log.Fatalf("could not download driver: %v", err)
	}
}
