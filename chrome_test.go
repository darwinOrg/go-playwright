package extpw_test

import (
	"fmt"
	"testing"
	"time"

	dgctx "github.com/darwinOrg/go-common/context"
	dglogger "github.com/darwinOrg/go-logger"
	extpw "github.com/darwinOrg/go-playwright"
)

func TestStartChrome(t *testing.T) {
	debugPort := 9222
	chromeCmd, _, _ := extpw.StartChrome(debugPort)
	time.Sleep(2 * time.Second)
	extpw.ShutdownChrome(chromeCmd, fmt.Sprintf("http://localhost:%d", debugPort))
}

func TestBossScraper(t *testing.T) {
	ctx := dgctx.SimpleDgContext()
	browserPath, err := extpw.FindCloakBrowserBinary()
	if err != nil {
		panic(err)
	}
	dglogger.Debugf(ctx, "browserPath: %s", browserPath)

	extBC, err := extpw.NewExtBrowserContext(ctx, &extpw.ExtPlaywrightOption{
		SkipInstallBrowsers: true,
		BrowserType:         extpw.MyBrowserType.Chrome,
		BrowserPath:         browserPath,
		UserDataDir:         "/Users/mac/.browser-boss",
		LaunchArgs:          extpw.GetStealthArgs(),
	})
	if err != nil {
		panic(err)
	}

	extPage, err := extBC.NewExtPage(ctx)
	if err != nil {
		panic(err)
	}
	defer extPage.Close()

	url := "https://www.zhipin.com/gongsi/job/480261c022ea03d81nV53tQ~.html?ka=company-jobs"
	err = extPage.NavigateWithLoadedState(ctx, url)
	if err != nil {
		panic(err)
	}
	extPage.RandomWaitLong(ctx)
}
