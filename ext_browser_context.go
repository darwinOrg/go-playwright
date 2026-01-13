package extpw

import (
	"fmt"
	"os"

	dgcoll "github.com/darwinOrg/go-common/collection"
	dgctx "github.com/darwinOrg/go-common/context"
	dglogger "github.com/darwinOrg/go-logger"
	"github.com/playwright-community/playwright-go"
)

type ExtBrowserContext struct {
	playwright.BrowserContext
	pw          *playwright.Playwright
	browserType playwright.BrowserType
	browser     playwright.Browser
	extPages    []*ExtPage
}

func NewDebugExtBrowserContext(ctx *dgctx.DgContext, extPwOpt *ExtPlaywrightOption) (*ExtBrowserContext, error) {
	extPwOpt.SkipInstallBrowsers = true
	pw, err := newPlaywright(ctx, extPwOpt)
	if err != nil {
		return nil, err
	}

	if extPwOpt.UserDataDir == "" {
		userHomeDir, err := os.UserHomeDir()
		if err != nil {
			dglogger.Errorf(ctx, "could not get work directory: %v", err)
			return nil, err
		}
		extPwOpt.UserDataDir = userHomeDir + "/ChromeProfile"
	}

	var remoteDebuggingPort int
	if extPwOpt.RemoteDebuggingPort > 0 {
		remoteDebuggingPort = extPwOpt.RemoteDebuggingPort
	} else {
		remoteDebuggingPort, err = getFreePort(ctx)
		if err != nil {
			return nil, err
		}
	}

	browserContext, err := pw.Chromium.LaunchPersistentContext(extPwOpt.UserDataDir,
		playwright.BrowserTypeLaunchPersistentContextOptions{
			Args:           []string{fmt.Sprintf("--remote-debugging-port=%d", remoteDebuggingPort)},
			ExecutablePath: playwright.String(extPwOpt.mustGetBrowserPath()),
			Headless:       playwright.Bool(extPwOpt.Headless),
		})
	if err != nil {
		dglogger.Errorf(ctx, "could not create browser context: %v", err)
		return nil, err
	}

	return buildExtBrowserContext(ctx, pw, pw.Chromium, nil, browserContext)
}

func ConnectDebugExtBrowserContext(ctx *dgctx.DgContext, extPwOpt *ExtPlaywrightOption) (*ExtBrowserContext, error) {
	pw, err := newPlaywright(ctx, extPwOpt)
	if err != nil {
		return nil, err
	}

	browser, err := connectOverCDP(pw, extPwOpt)
	if err != nil {
		dglogger.Errorf(ctx, "could not connect over CDP: %v", err)
		// 确保在连接失败时清理资源
		_ = pw.Stop()
		return nil, err
	}

	// 添加检查确保浏览器连接有效
	if browser == nil {
		dglogger.Error(ctx, "browser connection is nil")
		_ = pw.Stop()
		return nil, fmt.Errorf("failed to establish browser connection")
	}

	// 检查浏览器是否已经关闭
	if browser.IsConnected() == false {
		dglogger.Error(ctx, "browser is not connected")
		_ = browser.Close()
		_ = pw.Stop()
		return nil, fmt.Errorf("browser is not connected")
	}

	extBrowserContext, err := buildExtBrowserContext(ctx, pw, pw.Chromium, browser, nil)
	if err != nil {
		dglogger.Errorf(ctx, "could not build ext browser context: %v", err)
		_ = browser.Close()
		_ = pw.Stop()
		return nil, err
	}

	return extBrowserContext, nil
}

func NewExtBrowserContext(ctx *dgctx.DgContext, extPwOpt *ExtPlaywrightOption) (*ExtBrowserContext, error) {
	pw, err := newPlaywright(ctx, extPwOpt)
	if err != nil {
		return nil, err
	}

	var executablePath *string
	if extPwOpt.BrowserPath != "" {
		executablePath = playwright.String(extPwOpt.BrowserPath)
	}
	var channel *string
	if extPwOpt.Channel != "" {
		channel = playwright.String(extPwOpt.Channel)
	}

	browserType := extPwOpt.getBrowserType(pw)
	var browserContext playwright.BrowserContext
	var browser playwright.Browser

	if extPwOpt.UserDataDir != "" {
		browserContext, err = browserType.LaunchPersistentContext(extPwOpt.UserDataDir,
			playwright.BrowserTypeLaunchPersistentContextOptions{
				ExecutablePath: executablePath,
				Channel:        channel,
				Headless:       playwright.Bool(extPwOpt.Headless),
			})
		if err != nil {
			dglogger.Errorf(ctx, "could not create browser context: %v", err)
			return nil, err
		}
	} else {
		browser, err = browserType.Launch(playwright.BrowserTypeLaunchOptions{
			ExecutablePath: executablePath,
			Channel:        channel,
			Headless:       playwright.Bool(extPwOpt.Headless),
		})
		if err != nil {
			dglogger.Errorf(ctx, "could not launch browser: %v", err)
			return nil, err
		}
	}

	return buildExtBrowserContext(ctx, pw, browserType, browser, browserContext)
}

func buildExtBrowserContext(ctx *dgctx.DgContext, pw *playwright.Playwright, browserType playwright.BrowserType,
	browser playwright.Browser, browserContext playwright.BrowserContext) (*ExtBrowserContext, error) {

	if pw == nil {
		return nil, fmt.Errorf("playwright instance is nil")
	}

	if browserType == nil {
		browserType = pw.Chromium
	}

	if browserContext == nil {
		if browser == nil {
			var err error
			browser, err = browserType.Launch()
			if err != nil {
				dglogger.Errorf(ctx, "could not launch browser: %v", err)
				return nil, err
			}
		}

		// 检查浏览器连接状态
		if !browser.IsConnected() {
			return nil, fmt.Errorf("browser is not connected")
		}

		browserContexts := browser.Contexts()
		if len(browserContexts) > 0 {
			browserContext = browserContexts[0]
		} else {
			var err error
			browserContext, err = browser.NewContext()
			if err != nil {
				dglogger.Errorf(ctx, "could not create browser context: %v", err)
				return nil, err
			}
		}
	}

	// 检查 browserContext 是否有效
	if browserContext == nil {
		return nil, fmt.Errorf("browser context is nil")
	}

	extBC := &ExtBrowserContext{
		pw:             pw,
		browserType:    browserType,
		browser:        browser,
		BrowserContext: browserContext,
	}

	extBC.extPages = dgcoll.MapToList(browserContext.Pages(), func(page playwright.Page) *ExtPage {
		return &ExtPage{
			Page:   page,
			extBC:  extBC,
			locked: false,
		}
	})

	return extBC, nil
}

func (bc *ExtBrowserContext) Close() {
	_ = bc.BrowserContext.Close()
	if bc.browser != nil {
		_ = bc.browser.Close()
	}
	_ = bc.pw.Stop()
}

func (bc *ExtBrowserContext) GetOrNewExtPage(ctx *dgctx.DgContext) (*ExtPage, error) {
	var openedPages []*ExtPage
	if len(bc.extPages) > 0 {
		for _, extPage := range bc.extPages {
			if extPage.IsClosed() {
				extPage.Close()
				extPage.locked = false
			} else {
				openedPages = append(openedPages, extPage)
			}
		}
	}

	bc.extPages = openedPages
	if len(bc.extPages) > 0 {
		bc.pw.Lock()
		defer bc.pw.Unlock()

		for _, extPage := range bc.extPages {
			if !extPage.locked {
				extPage.locked = true
				return extPage, nil
			}
		}
	}

	return bc.NewExtPage(ctx)
}

func (bc *ExtBrowserContext) NewExtPage(ctx *dgctx.DgContext) (*ExtPage, error) {
	page, err := bc.NewPage()
	if err != nil {
		dglogger.Errorf(ctx, "could not create page: %v", err)
		return nil, err
	}

	return bc.BuildExtPage(page), nil
}

func (bc *ExtBrowserContext) BuildExtPage(page playwright.Page) *ExtPage {
	extPage := &ExtPage{
		Page:   page,
		extBC:  bc,
		locked: true,
	}
	bc.extPages = append(bc.extPages, extPage)

	return extPage
}
