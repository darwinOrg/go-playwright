package extpw

import (
	"fmt"
	"testing"
	"time"
)

func TestStartChrome(t *testing.T) {
	debugPort := 9221
	chromeCmd, _, _ := StartChrome(debugPort)
	time.Sleep(2 * time.Second)
	ShutdownChrome(chromeCmd, fmt.Sprintf("http://localhost:%d", debugPort))
}
