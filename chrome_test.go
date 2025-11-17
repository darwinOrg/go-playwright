package extpw

import (
	"testing"
	"time"
)

func TestStartChrome(t *testing.T) {
	debugPort := 9222
	_, _, _ = StartChrome(debugPort)
	time.Sleep(2 * time.Second)
	//ShutdownChrome(chromeCmd, fmt.Sprintf("http://localhost:%d", debugPort))
}
