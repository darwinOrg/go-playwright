package extpw_test

import (
	"fmt"
	"testing"
	"time"

	dgctx "github.com/darwinOrg/go-common/context"
	extpw "github.com/darwinOrg/go-playwright"
)

const debugPort = 9222

func TestStartChrome(t *testing.T) {
	chromeCmd, _, _ := extpw.StartChrome(debugPort)
	defer extpw.ShutdownChrome(chromeCmd, fmt.Sprintf("http://localhost:%d", debugPort))
	time.Sleep(5 * time.Second)
}

func TestAntiCrawler(t *testing.T) {
	chromeCmd, _, _ := extpw.StartChrome(debugPort)
	defer extpw.ShutdownChrome(chromeCmd, fmt.Sprintf("http://localhost:%d", debugPort))

	ctx := dgctx.SimpleDgContext()
	extPage, err := extpw.ConnectNewDebugExtPage(ctx, &extpw.ExtPlaywrightOption{SkipInstallBrowsers: true})
	if err != nil {
		panic(err)
	}

	err = extPage.NavigateWithLoadedState(ctx, "https://www.zhipin.com/gongsi/job/5d627415a46b4a750nJ9.html?ka=company-jobs")
	if err != nil {
		panic(err)
	}
	time.Sleep(10 * time.Minute)
}
