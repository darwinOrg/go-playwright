package extpw

import (
	"testing"
	"time"
)

func TestStartChrome(t *testing.T) {
	chromeCmd, _ := StartChrome()
	time.Sleep(2 * time.Second)
	ShutdownChrome(chromeCmd, "")
}
