package extpw

import (
	"fmt"
	"net"

	dgctx "github.com/darwinOrg/go-common/context"
	dglogger "github.com/darwinOrg/go-logger"
	"github.com/playwright-community/playwright-go"
)

var MyBrowserType = struct {
	Chrome  string
	Firefox string
	Webkit  string
}{"chrome", "firefox", "webkit"}

type ExtPlaywrightOption struct {
	SkipInstallBrowsers bool   `json:"skipInstallBrowsers" mapstructure:"skipInstallBrowsers"`
	Headless            bool   `json:"headless" mapstructure:"headless"`
	BrowserType         string `json:"browserType" mapstructure:"browserType"`
	Channel             string `json:"channel" mapstructure:"channel"`
	DriverDirectory     string `json:"driverDirectory" mapstructure:"driverDirectory"`
	BrowserPath         string `json:"browserPath" mapstructure:"browserPath"`
	UserDataDir         string `json:"userDataDir" mapstructure:"userDataDir"`
	RemoteDebuggingHost string `json:"remoteDebuggingHost" mapstructure:"remoteDebuggingHost"`
	RemoteDebuggingPort int    `json:"remoteDebuggingPort" mapstructure:"remoteDebuggingPort"`
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
		chromePath, err := FindChromePath()
		if err != nil {
			panic(err)
		}

		return chromePath
	}
}

func newPlaywright(ctx *dgctx.DgContext, extPwOpt *ExtPlaywrightOption) (*playwright.Playwright, error) {
	runOption := &playwright.RunOptions{
		SkipInstallBrowsers: extPwOpt.SkipInstallBrowsers,
		DriverDirectory:     extPwOpt.DriverDirectory,
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
