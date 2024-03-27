package extpw

import (
	dgctx "github.com/darwinOrg/go-common/context"
	dgerr "github.com/darwinOrg/go-common/enums/error"
	dglogger "github.com/darwinOrg/go-logger"
	"github.com/playwright-community/playwright-go"
	"math/rand"
	"time"
)

type ExtPage struct {
	playwright.Page
	extBC  *ExtBrowserContext
	locked bool
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
	page, err := p.extBC.ExpectPage(cb)
	if err != nil {
		dglogger.Errorf(ctx, "Page.ExpectExtPage error: %v", err)
		return nil, err
	}

	return p.extBC.buildExtPage(page), nil
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

func (p *ExtPage) NavigateWithLoadedState(ctx *dgctx.DgContext, url string) error {
	return p.Navigate(ctx, url, playwright.PageGotoOptions{
		WaitUntil: playwright.WaitUntilStateLoad,
	})
}

func (p *ExtPage) WaitForLoadStateLoad(ctx *dgctx.DgContext) error {
	err := p.WaitForLoadState(playwright.PageWaitForLoadStateOptions{
		State: playwright.LoadStateLoad,
	})
	if err != nil {
		dglogger.Errorf(ctx, "Page.WaitForLoadStateLoad error: %v", err)
		return err
	}
	return err
}

func (p *ExtPage) WaitForDomContentLoaded(ctx *dgctx.DgContext) error {
	err := p.WaitForLoadState(playwright.PageWaitForLoadStateOptions{
		State: playwright.LoadStateDomcontentloaded,
	})
	if err != nil {
		dglogger.Errorf(ctx, "Page.WaitForLoadStateLoad error: %v", err)
		return err
	}
	return err
}

func (p *ExtPage) Navigate(ctx *dgctx.DgContext, url string, options ...playwright.PageGotoOptions) error {
	_, err := p.Goto(url, options...)
	if err != nil {
		dglogger.Errorf(ctx, "Page.Goto url[%s] error: %v", url, err)
		return err
	}
	return nil
}

func (p *ExtPage) ReloadWithLoadedState(ctx *dgctx.DgContext) error {
	_, err := p.Reload(playwright.PageReloadOptions{
		WaitUntil: playwright.WaitUntilStateLoad,
	})
	if err != nil {
		dglogger.Errorf(ctx, "Page.Reload url[%s] error: %v", p.URL(), err)
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
	response, err := p.ExpectResponse(urlOrPredicate, cb)
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
