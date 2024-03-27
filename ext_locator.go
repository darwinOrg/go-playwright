package extpw

import (
	dgctx "github.com/darwinOrg/go-common/context"
	dglogger "github.com/darwinOrg/go-logger"
	"github.com/playwright-community/playwright-go"
	"strings"
)

type ExtLocator struct {
	playwright.Locator
	selectors []string
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
	return text
}

func (l *ExtLocator) MustTextContent(ctx *dgctx.DgContext) string {
	if !l.Exists(ctx) {
		return ""
	}

	text, err := l.Locator.TextContent()
	if err != nil {
		dglogger.Errorf(ctx, "locator[%s] text content error: %v", strings.Join(l.selectors, " "), err)
	}
	return text
}

func (l *ExtLocator) ExtLocator(selector string) *ExtLocator {
	return &ExtLocator{
		Locator:   l.Locator.Locator(selector),
		selectors: append(l.selectors, selector),
	}
}
