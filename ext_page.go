package extpw

import (
	"math/rand"
	"strings"
	"time"

	dgctx "github.com/darwinOrg/go-common/context"
	dgerr "github.com/darwinOrg/go-common/enums/error"
	"github.com/darwinOrg/go-common/utils"
	dglogger "github.com/darwinOrg/go-logger"
	"github.com/playwright-community/playwright-go"
)

var (
	defaultTimeoutMillis = 10_000.0
)

type ExtPage struct {
	playwright.Page
	extBC     *ExtBrowserContext
	locked    bool
	suspended bool
}

func NewDebugExtPage(ctx *dgctx.DgContext, extPwOpt *ExtPlaywrightOption) (*ExtPage, error) {
	extBC, err := NewDebugExtBrowserContext(ctx, extPwOpt)
	if err != nil {
		return nil, err
	}

	return extBC.GetOrNewExtPage(ctx)
}

func ConnectDebugExtPage(ctx *dgctx.DgContext, extPwOpt *ExtPlaywrightOption) (*ExtPage, error) {
	extBC, err := ConnectDebugExtBrowserContext(ctx, extPwOpt)
	if err != nil {
		return nil, err
	}

	return extBC.GetOrNewExtPage(ctx)
}

func ConnectNewDebugExtPage(ctx *dgctx.DgContext, extPwOpt *ExtPlaywrightOption) (*ExtPage, error) {
	extBC, err := ConnectDebugExtBrowserContext(ctx, extPwOpt)
	if err != nil {
		return nil, err
	}

	return extBC.NewExtPage(ctx)
}

func (p *ExtPage) NewExtPage(ctx *dgctx.DgContext) (*ExtPage, error) {
	return p.extBC.NewExtPage(ctx)
}

func (p *ExtPage) ExpectExtPage(ctx *dgctx.DgContext, cb func() error) (*ExtPage, error) {
	p.CheckSuspend(ctx)

	page, err := p.extBC.ExpectPage(cb)
	if err != nil {
		dglogger.Errorf(ctx, "Page.ExpectExtPage error: %v", err)
		if page != nil {
			_ = page.Close()
		}
		return nil, err
	}

	return p.extBC.BuildExtPage(page), nil
}

func (p *ExtPage) ExtContext() *ExtBrowserContext {
	return p.extBC
}

func (p *ExtPage) Release() {
	p.extBC.pw.Lock()
	defer p.extBC.pw.Unlock()

	for _, extPage := range p.extBC.extPages {
		if extPage == p {
			extPage.locked = false
			break
		}
	}
}

func (p *ExtPage) Close() {
	_ = p.Page.Close()
}

func (p *ExtPage) CloseAll() {
	p.Close()
	p.extBC.Close()
}

func (p *ExtPage) ReNewPageByError(err error) {
	if strings.Contains(err.Error(), "target closed") {
		p.Close()
		_ = utils.Retry(3, time.Second, func() error {
			newPage, ne := p.extBC.NewPage()
			if ne != nil {
				return ne
			}
			p.Page = newPage
			return nil
		})
	}
}

func (p *ExtPage) NavigateWithLoadedState(ctx *dgctx.DgContext, url string) error {
	defer func() {
		if err := recover(); err != nil {
			dglogger.Errorf(ctx, "Page.Goto url[%s] panic: %v", url, err)
		}
	}()

	p.CheckSuspend(ctx)

	// 使用 CDP 注入脚本，确保在页面加载之前就执行
	_ = p.extBC.InjectScriptViaCDP(p.Page)

	err := p.Navigate(ctx, url, playwright.PageGotoOptions{
		WaitUntil: playwright.WaitUntilStateLoad,
	})
	if err != nil {
		dglogger.Errorf(ctx, "Page.WaitUntilStateLoad error: %v | url: %s", err, url)
		return err
	}

	// 导航后再次注入指纹脚本
	_, _ = p.Evaluate(InitScript)

	return nil
}

func (p *ExtPage) Navigate(ctx *dgctx.DgContext, url string, options ...playwright.PageGotoOptions) error {
	defer func() {
		if err := recover(); err != nil {
			dglogger.Errorf(ctx, "Page.Goto url[%s] panic: %v", url, err)
		}
	}()

	p.CheckSuspend(ctx)

	// 使用 CDP 注入脚本，确保在页面加载之前就执行
	_ = p.extBC.InjectScriptViaCDP(p.Page)

	_, err := p.Goto(url, options...)
	if err != nil {
		dglogger.Errorf(ctx, "Page.Goto url[%s] error: %v", url, err)
		p.ReNewPageByError(err)
		return err
	}

	// 导航后再次注入指纹脚本
	_, _ = p.Evaluate(InitScript)

	return nil
}

func (p *ExtPage) ReloadWithLoadedState(ctx *dgctx.DgContext) error {
	defer func() {
		if err := recover(); err != nil {
			dglogger.Errorf(ctx, "Reload[%s] panic: %v", p.URL(), err)
		}
	}()

	p.CheckSuspend(ctx)
	_, err := p.Reload(playwright.PageReloadOptions{
		WaitUntil: playwright.WaitUntilStateLoad,
	})
	if err != nil {
		dglogger.Errorf(ctx, "Page.Reload url[%s] error: %v", p.URL(), err)
		p.ReNewPageByError(err)
		return err
	}
	return nil
}

func (p *ExtPage) WaitForLoadStateLoad(ctx *dgctx.DgContext) error {
	defer func() {
		if err := recover(); err != nil {
			dglogger.Errorf(ctx, "WaitForLoadStateLoad[%s] panic: %v", p.URL(), err)
		}
	}()

	err := p.WaitForLoadState(playwright.PageWaitForLoadStateOptions{
		State: playwright.LoadStateLoad,
	})
	if err != nil {
		dglogger.Errorf(ctx, "Page.WaitForLoadStateLoad error: %v", err)
		p.ReNewPageByError(err)
		return err
	}
	return nil
}

func (p *ExtPage) WaitForDomContentLoaded(ctx *dgctx.DgContext) error {
	defer func() {
		if err := recover(); err != nil {
			dglogger.Errorf(ctx, "WaitForDomContentLoaded[%s] panic: %v", p.URL(), err)
		}
	}()

	err := p.WaitForLoadState(playwright.PageWaitForLoadStateOptions{
		State: playwright.LoadStateDomcontentloaded,
	})
	if err != nil {
		dglogger.Errorf(ctx, "Page.WaitForLoadStateLoad error: %v", err)
		p.ReNewPageByError(err)
		return err
	}
	return err
}

func (p *ExtPage) WaitForSelectorStateVisible(ctx *dgctx.DgContext, selector string) error {
	defer func() {
		if err := recover(); err != nil {
			dglogger.Errorf(ctx, "WaitForSelectorStateVisible[%s] panic: %v", selector, err)
		}
	}()

	_, err := p.WaitForSelector(selector, playwright.PageWaitForSelectorOptions{
		State:   playwright.WaitForSelectorStateVisible,
		Timeout: playwright.Float(defaultTimeoutMillis),
	})
	if err != nil {
		dglogger.Errorf(ctx, "Page.WaitForSelector error: %v", err)
		return err
	}

	return nil
}

func (p *ExtPage) RandomWaitShort(ctx *dgctx.DgContext) {
	p.RandomWaitRange(ctx, 100, 1_000)
}

func (p *ExtPage) RandomWaitMiddle(ctx *dgctx.DgContext) {
	p.RandomWaitRange(ctx, 3_000, 6_000)
}

func (p *ExtPage) RandomWaitLong(ctx *dgctx.DgContext) {
	p.RandomWaitRange(ctx, 10_000, 20_000)
}

func (p *ExtPage) RandomWaitRange(ctx *dgctx.DgContext, min, max int) {
	milli := time.Duration(rand.Intn(max-min) + min)
	dglogger.Infof(ctx, "等待 %d 毫秒", milli)
	time.Sleep(time.Millisecond * milli)
}

func (p *ExtPage) ExpectResponseText(ctx *dgctx.DgContext, urlOrPredicate string, cb func() error) (string, error) {
	defer func() {
		if err := recover(); err != nil {
			dglogger.Errorf(ctx, "ExpectResponseText[%s] panic: %v", urlOrPredicate, err)
		}
	}()

	p.CheckSuspend(ctx)
	response, err := p.ExpectResponse(urlOrPredicate, cb, playwright.PageExpectResponseOptions{
		Timeout: &defaultTimeoutMillis,
	})
	if err != nil {
		dglogger.Errorf(ctx, "Page.ExpectResponseText error: %v", err)
		return "", err
	}
	if !response.Ok() {
		return "", dgerr.SYSTEM_ERROR
	}

	text, err := response.Text()
	if err != nil {
		dglogger.Errorf(ctx, "get response text error: %v", err)
		return "", err
	}
	if text == "" {
		dglogger.Errorf(ctx, "get response text empty: %s", urlOrPredicate)
		return "", dgerr.SYSTEM_ERROR
	}

	return text, nil
}

func (p *ExtPage) HtmlContent(ctx *dgctx.DgContext) string {
	content, err := p.Content()
	if err != nil {
		dglogger.Errorf(ctx, "Page.Content error: %v", err)
	}

	return content
}

func (p *ExtPage) ExtLocator(selectors ...string) *ExtLocator {
	var locator playwright.Locator
	for _, selector := range selectors {
		if locator == nil {
			locator = p.Locator(selector)
		} else {
			locator = locator.Locator(selector)
		}
	}

	return &ExtLocator{
		extPage:   p,
		Locator:   locator,
		selectors: selectors,
	}
}

func (p *ExtPage) MustInnerText(ctx *dgctx.DgContext, selectors ...string) string {
	locator := p.ExtLocator(selectors...)

	return locator.MustInnerText(ctx)
}

func (p *ExtPage) MustTextContent(ctx *dgctx.DgContext, selectors ...string) string {
	locator := p.ExtLocator(selectors...)

	return locator.MustTextContent(ctx)
}

func (p *ExtPage) Exists(ctx *dgctx.DgContext, selector string) bool {
	return p.ExtLocator(selector).Exists(ctx)
}

func (p *ExtPage) Click(ctx *dgctx.DgContext, selector string) error {
	defer func() {
		if err := recover(); err != nil {
			dglogger.Errorf(ctx, "Click[%s] panic: %v", selector, err)
		}
	}()

	p.CheckSuspend(ctx)
	locator := p.ExtLocator(selector)
	if locator.Exists(ctx) {
		err := locator.Click()
		if err != nil {
			dglogger.Errorf(ctx, "Page.Click[%s] error: %v", selector, err)
			return err
		}
	}

	return nil
}

func (p *ExtPage) Suspend() {
	p.suspended = true
}

func (p *ExtPage) Continue() {
	p.suspended = false
}

func (p *ExtPage) CheckSuspend(ctx *dgctx.DgContext) {
	for p.suspended {
		p.RandomWaitMiddle(ctx)
	}
}
