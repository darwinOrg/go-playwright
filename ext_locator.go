package extpw

import (
	"fmt"
	"strings"

	dgcoll "github.com/darwinOrg/go-common/collection"
	dgctx "github.com/darwinOrg/go-common/context"
	dglogger "github.com/darwinOrg/go-logger"
	"github.com/mxschmitt/playwright-go"
)

type ExtLocator struct {
	extPage *ExtPage
	playwright.Locator
	selector string
}

func NewExtLocator(extPage *ExtPage, locator playwright.Locator, selector string) *ExtLocator {
	return &ExtLocator{
		extPage:  extPage,
		Locator:  locator,
		selector: selector,
	}
}

func (l *ExtLocator) ExtPage() *ExtPage {
	return l.extPage
}

func (l *ExtLocator) ExtLocator(selector string) *ExtLocator {
	return &ExtLocator{
		extPage:  l.extPage,
		Locator:  l.Locator.Locator(selector),
		selector: fmt.Sprintf("%s %s", l.selector, selector),
	}
}

func (l *ExtLocator) Exists(ctx *dgctx.DgContext) bool {
	defer func() {
		if err := recover(); err != nil {
			dglogger.Errorf(ctx, "ExtLocator.Exists[%s] panic: %v", l.selector, err)
		}
	}()

	count, err := l.Locator.Count()
	if err != nil {
		dglogger.Errorf(ctx, "locator[%s] count error: %v", l.selector, err)
		return false
	}
	return count > 0
}

func (l *ExtLocator) MustInnerText(ctx *dgctx.DgContext) string {
	if !l.Exists(ctx) {
		return ""
	}

	text, err := l.Locator.InnerText()
	if err != nil {
		dglogger.Errorf(ctx, "locator[%s] inner text error: %v", l.selector, err)
	}
	return strings.TrimSpace(text)
}

func (l *ExtLocator) MustInnerTexts(ctx *dgctx.DgContext) []string {
	allLocators, err := l.ExtAll(ctx)
	if err != nil {
		return []string{}
	}

	return dgcoll.MapToList(allLocators, func(locator *ExtLocator) string {
		return locator.MustInnerText(ctx)
	})
}

func (l *ExtLocator) MustTextContent(ctx *dgctx.DgContext) string {
	if !l.Exists(ctx) {
		return ""
	}

	text, err := l.Locator.TextContent()
	if err != nil {
		dglogger.Errorf(ctx, "locator[%s] text content error: %v", l.selector, err)
	}
	return strings.TrimSpace(text)
}

func (l *ExtLocator) MustTextContents(ctx *dgctx.DgContext) []string {
	allLocators, err := l.ExtAll(ctx)
	if err != nil {
		return []string{}
	}

	return dgcoll.MapToList(allLocators, func(locator *ExtLocator) string {
		return locator.MustTextContent(ctx)
	})
}

func (l *ExtLocator) MustAttribute(ctx *dgctx.DgContext, attr string) string {
	if !l.Exists(ctx) {
		return ""
	}

	text, err := l.Locator.GetAttribute(attr)
	if err != nil {
		dglogger.Errorf(ctx, "locator[%s] get attribute error: %v", l.selector, err)
	}
	return strings.TrimSpace(text)
}

func (l *ExtLocator) MustAttributes(ctx *dgctx.DgContext, attr string) []string {
	allLocators, err := l.ExtAll(ctx)
	if err != nil {
		return []string{}
	}

	return dgcoll.MapToList(allLocators, func(locator *ExtLocator) string {
		return locator.MustAttribute(ctx, attr)
	})
}

func (l *ExtLocator) ExtAll(ctx *dgctx.DgContext) ([]*ExtLocator, error) {
	if !l.Exists(ctx) {
		return []*ExtLocator{}, nil
	}

	count, err := l.Locator.Count()
	if err != nil {
		dglogger.Errorf(ctx, "locator[%s] all error: %v", l.selector, err)
		return nil, err
	}
	if count == 0 {
		return []*ExtLocator{}, nil
	}

	extLocators := make([]*ExtLocator, count)
	for i := 0; i < count; i++ {
		selector := fmt.Sprintf("%s:nth-child(%d)", l.selector, i+1)
		extLocators[i] = l.extPage.ExtLocator(selector)
	}

	return extLocators, nil
}

func (l *ExtLocator) MustClick(ctx *dgctx.DgContext) {
	l.CheckSuspend(ctx)
	err := l.Click()
	if err != nil {
		dglogger.Errorf(ctx, "locator[%s] click error: %v", l.selector, err)
	}
}

func (l *ExtLocator) HasClass(ctx *dgctx.DgContext, class string) bool {
	cls := l.MustAttribute(ctx, "class")
	if cls == "" {
		return false
	}

	if strings.Contains(cls, " ") {
		classes := strings.Split(cls, " ")
		return dgcoll.Contains(classes, class)
	}

	return cls == class
}

func (l *ExtLocator) GetSelector() string {
	return l.selector
}

func (l *ExtLocator) Suspend() {
	l.extPage.suspended = true
}

func (l *ExtLocator) Continue() {
	l.extPage.suspended = false
}

func (l *ExtLocator) CheckSuspend(ctx *dgctx.DgContext) {
	for l.extPage.suspended {
		l.extPage.RandomWaitMiddle(ctx)
	}
}
