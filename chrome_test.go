package extpw_test

import (
	"fmt"
	"testing"
	"time"

	extpw "github.com/darwinOrg/go-playwright"
)

func TestStartChrome(t *testing.T) {
	debugPort := 9222
	chromeCmd, _, _ := extpw.StartChrome(debugPort)
	time.Sleep(2 * time.Second)
	extpw.ShutdownChrome(chromeCmd, fmt.Sprintf("http://localhost:%d", debugPort))
}
