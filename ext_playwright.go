package extpw

import (
	"fmt"
	dgctx "github.com/darwinOrg/go-common/context"
	dglogger "github.com/darwinOrg/go-logger"
	"github.com/playwright-community/playwright-go"
	"net"
)

var MyBrowserType = struct {
	Chrome  int
	Firefox int
	Webkit  int
}{0, 1, 2}

type ExtPlaywrightOption struct {
	SkipInstallBrowsers bool
	Headless            bool
	BrowserType         int
	Channel             string
	BrowserPath         string
	UserDataDir         string
	RemoteDebuggingHost string
	RemoteDebuggingPort int
}

func (opt *ExtPlaywrightOption) getBrowserType(pw *playwright.Playwright) playwright.BrowserType {
	switch opt.BrowserType {
	case MyBrowserType.Firefox:
		return pw.Firefox
	case MyBrowserType.Webkit:
		return pw.WebKit
	default:
		return pw.Chromium
	}
}

func (opt *ExtPlaywrightOption) mustGetRemoteDebuggingPort() int {
	if opt.RemoteDebuggingPort > 0 {
		return opt.RemoteDebuggingPort
	} else {
		return 9222
	}
}

func (opt *ExtPlaywrightOption) mustGetBrowserPath() string {
	if opt.BrowserPath != "" {
		return opt.BrowserPath
	} else {
		return "/Applications/Google Chrome.app/Contents/MacOS/Google Chrome"
	}
}

func newPlaywright(ctx *dgctx.DgContext, extPwOpt *ExtPlaywrightOption) (*playwright.Playwright, error) {
	runOption := &playwright.RunOptions{
		SkipInstallBrowsers: extPwOpt.SkipInstallBrowsers,
	}
	pw, err := playwright.Run(runOption)
	if err != nil {
		dglogger.Errorf(ctx, "could not start playwright: %v", err)
		return nil, err
	}
	return pw, nil
}

func connectOverCDP(pw *playwright.Playwright, extPwOpt *ExtPlaywrightOption) (playwright.Browser, error) {
	if extPwOpt.RemoteDebuggingHost == "" {
		extPwOpt.RemoteDebuggingHost = "localhost"
	}

	return pw.Chromium.ConnectOverCDP(fmt.Sprintf("http://%s:%d", extPwOpt.RemoteDebuggingHost, extPwOpt.mustGetRemoteDebuggingPort()))
}

func getFreePort(ctx *dgctx.DgContext) (int, error) {
	addr, err := net.ResolveTCPAddr("tcp", "localhost:0")
	if err != nil {
		dglogger.Errorf(ctx, "net.ResolveTCPAddr: %v", err)
		return 0, err
	}

	l, err := net.ListenTCP("tcp", addr)
	if err != nil {
		dglogger.Errorf(ctx, "net.ListenTCP: %v", err)
		return 0, err
	}
	defer func() {
		_ = l.Close()
	}()

	return l.Addr().(*net.TCPAddr).Port, nil
}
