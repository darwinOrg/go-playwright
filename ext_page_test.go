package extpw_test

import (
	"os"
	"testing"
	"time"

	dgctx "github.com/darwinOrg/go-common/context"
	extpw "github.com/darwinOrg/go-playwright"
)

func TestNewPage(t *testing.T) {
	ctx := dgctx.SimpleDgContext()
	extPage := newPage(ctx)
	defer func() {
		_ = extPage.CloseAll(ctx)
	}()

	time.Sleep(time.Minute)
}

func newPage(ctx *dgctx.DgContext) *extpw.ExtPage {
	extBC, err := newBrowserContext(ctx)
	if err != nil {
		panic(err)
	}
	extPage, err := extBC.GetOrNewExtPage(ctx)
	if err != nil {
		panic(err)
	}
	return extPage
}

func newBrowserContext(ctx *dgctx.DgContext) (*extpw.ExtBrowserContext, error) {
	cloakBrowserPath, err := extpw.FindCloakBrowserBinary()
	if err != nil {
		return nil, err
	}

	launchArgs := []string{"--start-maximized"}
	launchArgs = append(launchArgs, extpw.GetStealthArgs()...)

	return extpw.NewExtBrowserContext(ctx, &extpw.ExtPlaywrightOption{
		SkipInstallBrowsers: true,
		UserDataDir:         os.Getenv("USER_DATA_DIR"),
		BrowserPath:         cloakBrowserPath,
		LaunchArgs:          launchArgs,
	})
}
