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
	pw               *playwright.Playwright
	browserType      playwright.BrowserType
	browser          playwright.Browser
	extPages         []*ExtPage
	cheatInitialized bool // 防爬虫信息是否已初始化
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
			Args: []string{
				fmt.Sprintf("--remote-debugging-port=%d", remoteDebuggingPort),
				"--disable-blink-features=AutomationControlled",
				"--disable-features=IsolateOrigins,site-per-process",
				"--no-sandbox",
				"--disable-setuid-sandbox",
				"--disable-dev-shm-usage",
				"--disable-accelerated-2d-canvas",
				"--disable-gpu",
				"--window-size=1920,1080",
			},
			ExecutablePath:    playwright.String(extPwOpt.mustGetBrowserPath()),
			Headless:          playwright.Bool(extPwOpt.Headless),
			IgnoreDefaultArgs: []string{"--enable-automation", "--enable-blink-features=IdleDetection"},
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
				Args: []string{
					"--disable-blink-features=AutomationControlled",
					"--disable-features=IsolateOrigins,site-per-process",
					"--no-sandbox",
					"--disable-setuid-sandbox",
					"--disable-dev-shm-usage",
					"--disable-accelerated-2d-canvas",
					"--disable-gpu",
					"--window-size=1920,1080",
				},
				IgnoreDefaultArgs: []string{"--enable-automation", "--enable-blink-features=IdleDetection"},
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
			Args: []string{
				"--disable-blink-features=AutomationControlled",
				"--disable-features=IsolateOrigins,site-per-process",
				"--no-sandbox",
				"--disable-setuid-sandbox",
				"--disable-dev-shm-usage",
				"--disable-accelerated-2d-canvas",
				"--disable-gpu",
				"--window-size=1920,1080",
			},
			IgnoreDefaultArgs: []string{"--enable-automation", "--enable-blink-features=IdleDetection"},
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
			browserContext, err = browser.NewContext(playwright.BrowserNewContextOptions{
				UserAgent:  playwright.String("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36"),
				Locale:     playwright.String("zh-CN"),
				TimezoneId: playwright.String("Asia/Shanghai"),
			})
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

	// 获取当前操作系统类型并生成对应的脚本
	osType := GetCurrentOS()
	initScript := GetInitScript(osType)

	// 为上下文注入指纹脚本（这样所有新页面都会自动应用）
	_ = browserContext.AddInitScript(playwright.Script{Content: &initScript})

	extBC := &ExtBrowserContext{
		pw:             pw,
		browserType:    browserType,
		browser:        browser,
		BrowserContext: browserContext,
	}

	// 为已存在的页面初始化防爬虫信息（只初始化一次）
	pages := browserContext.Pages()
	if len(pages) > 0 {
		// 只在第一个页面上初始化 CDP 配置
		_ = extBC.InitCheatInfoOnPage(pages[0])
	}

	extBC.extPages = dgcoll.MapToList(pages, func(page playwright.Page) *ExtPage {
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

	// 脚本已经通过 BrowserContext.AddInitScript 注入，不需要在这里重复初始化
	return bc.BuildExtPage(page), nil
}

func (bc *ExtBrowserContext) BuildExtPage(page playwright.Page) *ExtPage {
	// 脚本已经通过 BrowserContext.AddInitScript 注入，不需要在这里重复初始化
	extPage := &ExtPage{
		Page:   page,
		extBC:  bc,
		locked: true,
	}
	bc.extPages = append(bc.extPages, extPage)

	return extPage
}

// InjectScriptViaCDP 使用 CDP 注入脚本到页面
func (bc *ExtBrowserContext) InjectScriptViaCDP(page playwright.Page) error {
	return bc.InjectScriptViaCDPWithOS(page)
}

// InjectScriptViaCDPWithOS 使用 CDP 注入脚本到页面（根据操作系统类型生成脚本）
func (bc *ExtBrowserContext) InjectScriptViaCDPWithOS(page playwright.Page) error {
	// 尝试使用 CDP 注入脚本
	cdpSession, err := page.Context().NewCDPSession(page)
	if err != nil {
		// CDP 不可用，静默失败
		return err
	}

	// 获取当前操作系统类型并生成对应的脚本
	osType := GetCurrentOS()
	initScript := GetInitScript(osType)

	_, err = cdpSession.Send("Page.addScriptToEvaluateOnNewDocument", map[string]interface{}{
		"source": initScript,
	})
	if err != nil {
		// CDP 注入失败，静默失败
		return err
	}

	return nil
}

// InitCheatInfoOnPage 在页面上初始化所有防爬虫信息（只在页面创建时调用一次）
func (bc *ExtBrowserContext) InitCheatInfoOnPage(page playwright.Page) error {
	// 只在第一次初始化时设置平台和屏幕信息
	if !bc.cheatInitialized {
		// 初始化所有防爬虫信息
		if err := bc.initCheatInfo(page); err != nil {
			// CDP 不可用，静默失败
			return err
		}
		bc.cheatInitialized = true
	}

	// 注入脚本（使用动态生成的脚本）
	return bc.InjectScriptViaCDPWithOS(page)
}

// initCheatInfo 初始化所有防爬虫信息
func (bc *ExtBrowserContext) initCheatInfo(page playwright.Page) error {
	if err := bc.initCheatPlatform(page); err != nil {
		return err
	}
	if err := bc.initCheatScreenInfo(page); err != nil {
		return err
	}
	if err := bc.initCheatWebGLCanvasInfo(page); err != nil {
		return err
	}
	return nil
}

// initCheatPlatform 设置 User Agent 和平台信息
func (bc *ExtBrowserContext) initCheatPlatform(page playwright.Page) error {
	cdpSession, err := page.Context().NewCDPSession(page)
	if err != nil {
		return err
	}

	// 获取当前操作系统类型
	osType := GetCurrentOS()
	config := GetPlatformConfig(osType)

	// 转换 brands 格式
	var brands []map[string]string
	for _, brand := range config.UserAgentMetadata.Brands {
		brands = append(brands, map[string]string{
			"brand":   brand.Brand,
			"version": brand.Version,
		})
	}

	// 设置 User Agent Override
	_, err = cdpSession.Send("Network.setUserAgentOverride", map[string]interface{}{
		"userAgent":      config.UserAgent,
		"platform":       config.Platform,
		"acceptLanguage": config.AcceptLanguage,
		"userAgentMetadata": map[string]interface{}{
			"brands":          brands,
			"fullVersion":     config.UserAgentMetadata.FullVersion,
			"platform":        config.UserAgentMetadata.Platform,
			"platformVersion": config.UserAgentMetadata.PlatformVersion,
			"architecture":    config.UserAgentMetadata.Architecture,
			"model":           config.UserAgentMetadata.Model,
			"mobile":          config.UserAgentMetadata.Mobile,
		},
	})
	if err != nil {
		return err
	}

	// 构建完整的 Client Hints
	secChUa := ""
	for i, brand := range config.UserAgentMetadata.Brands {
		if i > 0 {
			secChUa += ", "
		}
		secChUa += fmt.Sprintf("\"%s\";v=\"%s\"", brand.Brand, brand.Version)
	}

	// 设置额外的 HTTP Headers
	_, err = cdpSession.Send("Network.setExtraHTTPHeaders", map[string]interface{}{
		"headers": map[string]string{
			"sec-ch-ua":                 secChUa,
			"sec-ch-ua-mobile":          "?0",
			"sec-ch-ua-platform":        config.SecChUaPlatform,
			"sec-ch-ua-arch":            fmt.Sprintf("\"%s\"", config.UserAgentMetadata.Architecture),
			"sec-ch-ua-full-version":    fmt.Sprintf("\"%s\"", config.UserAgentMetadata.FullVersion),
			"sec-ch-ua-bitness":         "\"64\"",
			"sec-ch-ua-model":           "\"\"",
			"sec-fetch-dest":            "document",
			"sec-fetch-mode":            "navigate",
			"sec-fetch-site":            "none",
			"sec-fetch-user":            "?1",
			"upgrade-insecure-requests": "1",
		},
	})
	return err
}

// initCheatScreenInfo 设置屏幕信息
func (bc *ExtBrowserContext) initCheatScreenInfo(page playwright.Page) error {
	cdpSession, err := page.Context().NewCDPSession(page)
	if err != nil {
		return err
	}

	// 获取当前操作系统类型
	osType := GetCurrentOS()
	config := GetScreenConfig(osType)

	// 设置设备指标
	_, err = cdpSession.Send("Emulation.setDeviceMetricsOverride", map[string]interface{}{
		"width":             config.Width,
		"height":            config.Height,
		"deviceScaleFactor": config.DeviceScaleFactor,
		"mobile":            config.Mobile,
		"screenWidth":       config.ScreenWidth,
		"screenHeight":      config.ScreenHeight,
		"positionX":         config.PositionX,
		"positionY":         config.PositionY,
		"viewport": map[string]interface{}{
			"x":      config.ViewportX,
			"y":      config.ViewportY,
			"width":  config.ViewportWidth,
			"height": config.ViewportHeight,
			"scale":  config.DeviceScaleFactor,
		},
	})
	return err
}

// initCheatWebGLCanvasInfo 设置 WebGL Canvas 信息
func (bc *ExtBrowserContext) initCheatWebGLCanvasInfo(page playwright.Page) error {
	cdpSession, err := page.Context().NewCDPSession(page)
	if err != nil {
		return err
	}

	// 通过 Runtime.evaluate 注入 WebGL 伪装脚本
	_, err = cdpSession.Send("Runtime.evaluate", map[string]interface{}{
		"expression":   InitScript,
		"awaitPromise": false,
	})
	return err
}
