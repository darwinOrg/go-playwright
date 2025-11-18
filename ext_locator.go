package extpw

import (
	"strings"

	dgcoll "github.com/darwinOrg/go-common/collection"
	dgctx "github.com/darwinOrg/go-common/context"
	dglogger "github.com/darwinOrg/go-logger"
	"github.com/playwright-community/playwright-go"
)

type ExtLocator struct {
	extPage *ExtPage
	playwright.Locator
	selectors []string
}

func (l *ExtLocator) ExtLocator(selector string) *ExtLocator {
	return &ExtLocator{
		extPage:   l.extPage,
		Locator:   l.Locator.Locator(selector),
		selectors: append(l.selectors, selector),
	}
}

func (l *ExtLocator) Exists(ctx *dgctx.DgContext) bool {
	count, err := l.Locator.Count()
	if err != nil {
		dglogger.Errorf(ctx, "locator[%s] count error: %v", strings.Join(l.selectors, " "), err)
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
		dglogger.Errorf(ctx, "locator[%s] inner text error: %v", strings.Join(l.selectors, " "), err)
	}
	return strings.TrimSpace(text)
}

func (l *ExtLocator) MustTextContent(ctx *dgctx.DgContext) string {
	if !l.Exists(ctx) {
		return ""
	}

	text, err := l.Locator.TextContent()
	if err != nil {
		dglogger.Errorf(ctx, "locator[%s] text content error: %v", strings.Join(l.selectors, " "), err)
	}
	return strings.TrimSpace(text)
}

func (l *ExtLocator) MustGetAttribute(ctx *dgctx.DgContext, attr string) string {
	if !l.Exists(ctx) {
		return ""
	}

	text, err := l.Locator.GetAttribute(attr)
	if err != nil {
		dglogger.Errorf(ctx, "locator[%s] get attribute error: %v", strings.Join(l.selectors, " "), err)
	}
	return strings.TrimSpace(text)
}

func (l *ExtLocator) ExtAll(ctx *dgctx.DgContext) ([]*ExtLocator, error) {
	if !l.Exists(ctx) {
		return []*ExtLocator{}, nil
	}

	allLocators, err := l.Locator.All()
	if err != nil {
		dglogger.Errorf(ctx, "locator[%s] all error: %v", strings.Join(l.selectors, " "), err)
		return nil, err
	}

	return dgcoll.MapToList(allLocators, func(locator playwright.Locator) *ExtLocator {
		return &ExtLocator{
			extPage:   l.extPage,
			Locator:   locator,
			selectors: l.selectors,
		}
	}), nil
}

func (l *ExtLocator) MustAllInnerTexts(ctx *dgctx.DgContext) []string {
	allLocators, err := l.ExtAll(ctx)
	if err != nil {
		return []string{}
	}

	return dgcoll.MapToList(allLocators, func(locator *ExtLocator) string {
		return locator.MustInnerText(ctx)
	})
}

func (l *ExtLocator) MustAllTextContents(ctx *dgctx.DgContext) []string {
	allLocators, err := l.ExtAll(ctx)
	if err != nil {
		return []string{}
	}

	return dgcoll.MapToList(allLocators, func(locator *ExtLocator) string {
		return locator.MustTextContent(ctx)
	})
}

func (l *ExtLocator) MustAllGetAttributes(ctx *dgctx.DgContext, attr string) []string {
	allLocators, err := l.ExtAll(ctx)
	if err != nil {
		return []string{}
	}

	return dgcoll.MapToList(allLocators, func(locator *ExtLocator) string {
		return locator.MustGetAttribute(ctx, attr)
	})
}

func (l *ExtLocator) MustClick(ctx *dgctx.DgContext) {
	l.CheckSuspend(ctx)
	err := l.Click()
	if err != nil {
		dglogger.Errorf(ctx, "locator[%s] click error: %v", strings.Join(l.selectors, " "), err)
	}
}

func (l *ExtLocator) HasClass(ctx *dgctx.DgContext, class string) bool {
	cls := l.MustGetAttribute(ctx, "class")
	if cls == "" {
		return false
	}

	if strings.Contains(cls, " ") {
		classes := strings.Split(cls, " ")
		return dgcoll.Contains(classes, class)
	}

	return cls == class
}

func (l *ExtLocator) GetSelectors() []string {
	return l.selectors
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
