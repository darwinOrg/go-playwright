package extpw_test

import (
	"fmt"
	"os"
	"testing"
	"time"

	dgctx "github.com/darwinOrg/go-common/context"
	dgsys "github.com/darwinOrg/go-common/sys"
	extpw "github.com/darwinOrg/go-playwright"
)

const debugPort = 9222

func TestStartChrome(t *testing.T) {
	chromeCmd, _, _ := extpw.StartChrome(debugPort, extpw.GetDefaultUserDataDir())
	defer extpw.ShutdownChrome(chromeCmd, fmt.Sprintf("http://localhost:%d", debugPort))
	time.Sleep(5 * time.Second)
}

func TestAntiCrawler(t *testing.T) {
	userDataDir := extpw.GetDefaultUserDataDir()
	//defer os.RemoveAll(userDataDir)

	ctx := dgctx.SimpleDgContext()
	extBC, err := extpw.ConnectDebugExtBrowserContext(ctx, &extpw.ExtPlaywrightOption{
		SkipInstallBrowsers: true,
		UserDataDir:         userDataDir,
	})
	if err != nil {
		panic(err)
	}
	defer extBC.Close()

	extPage, err := extBC.NewExtPage(ctx)
	if err != nil {
		panic(err)
	}
	defer extPage.Close()

	err = extPage.NavigateWithLoadedState(ctx, "https://www.zhipin.com/gongsi/job/5d627415a46b4a750nJ9.html?ka=company-jobs")
	if err != nil {
		panic(err)
	}
	dgsys.HangupApplication()
}

func TestAntiCrawler2(t *testing.T) {
	// 创建临时用户数据目录
	tempUserDataDir := extpw.GetTempUserDataDir()
	defer os.RemoveAll(tempUserDataDir)

	fmt.Println("开始测试防爬虫功能...")

	ctx := dgctx.SimpleDgContext()
	extBC, err := extpw.NewDebugExtBrowserContext(ctx, &extpw.ExtPlaywrightOption{
		SkipInstallBrowsers: true,
		UserDataDir:         tempUserDataDir,
	})
	if err != nil {
		panic(err)
	}
	defer extBC.Close()

	extPage, err := extBC.NewExtPage(ctx)
	if err != nil {
		panic(err)
	}
	defer extPage.Close()

	fmt.Println("导航到测试页面...")
	err = extPage.NavigateWithLoadedState(ctx, "https://www.zhipin.com/gongsi/job/5d627415a46b4a750nJ9.html?ka=company-jobs")
	if err != nil {
		panic(err)
	}

	fmt.Println("等待 30 秒观察页面行为...")
	time.Sleep(30 * time.Second)

	currentURL := extPage.URL()
	fmt.Printf("当前页面 URL: %s\n", currentURL)

	if currentURL != "https://www.zhipin.com/gongsi/job/5d627415a46b4a750nJ9.html?ka=company-jobs" {
		fmt.Printf("警告：页面已被重定向到: %s\n", currentURL)
	} else {
		fmt.Println("成功：页面没有被重定向")
	}
}
