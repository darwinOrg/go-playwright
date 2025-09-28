package extpw

import (
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"syscall"
	"time"

	"github.com/go-ole/go-ole"
	"github.com/go-ole/go-ole/oleutil"
)

func StartChrome() (*exec.Cmd, error) {
	chrome, err := FindChromePath()
	if err != nil {
		return nil, err
	}

	homeDir, _ := os.UserHomeDir()
	cmd := exec.Command(chrome,
		"--remote-debugging-port=9222",
		"--no-first-run",
		"--no-default-browser-check",
		"--disable-gpu",
		"--disable-extensions",
		"--disable-plugins",
		"--disable-sync",
		"--user-data-dir="+homeDir+"/ChromeProfile",
		// 静默日志输出的关键参数
		"--log-level=3",                     // 只显示致命错误
		"--silent-startup",                  // 静默启动
		"--disable-dev-shm-usage",           // 减少崩溃
		"--disable-logging",                 // 禁用日志记录
		"--disable-ipc-flooding-protection", // 减少日志
	)

	// 将 stdout 和 stderr 重定向到空设备
	var nullWriter *os.File
	if runtime.GOOS == "windows" {
		nullWriter, _ = os.OpenFile("NUL", os.O_WRONLY, 0)
	} else {
		nullWriter, _ = os.OpenFile("/dev/null", os.O_WRONLY, 0)
	}
	cmd.Stdout = nullWriter
	cmd.Stderr = nullWriter

	if err := cmd.Start(); err != nil {
		return nil, err
	}

	log.Printf("Started Chrome with PID: %d", cmd.Process.Pid)
	return cmd, nil
}

func CloseChrome(chromeCmd *exec.Cmd) {
	if chromeCmd != nil && chromeCmd.Process != nil {
		if runtime.GOOS == "windows" {
			_ = chromeCmd.Process.Kill()
		} else {
			_ = chromeCmd.Process.Signal(syscall.SIGTERM)
		}

		done := make(chan error, 1)
		go func() { done <- chromeCmd.Wait() }()

		select {
		case <-done:
			log.Println("Chrome 已退出。")
		case <-time.After(5 * time.Second):
			log.Println("强制终止 Chrome...")
			_ = chromeCmd.Process.Kill()
		}
	}
}

// FindChromePath searches for the Google Chrome executable on different operating systems.
func FindChromePath() (string, error) {
	switch runtime.GOOS {
	case "windows":
		return findChromeOnWindows()
	case "darwin":
		return findChromeOnMacOS()
	case "linux":
		fallthrough
	default:
		return findChromeOnUnixLike()
	}
}

// findChromeOnWindows looks for Chrome in typical installation locations and the registry.
func findChromeOnWindows() (string, error) {
	// Check common installation paths first.
	for _, path := range []string{
		`C:\Program Files\Google\Chrome\Application\chrome.exe`,
		`C:\Program Files (x86)\Google\Chrome\Application\chrome.exe`,
	} {
		if _, err := os.Stat(path); err == nil {
			return path, nil
		}
	}

	// If not found in common paths, attempt to query the Windows registry.
	return findChromeInRegistry()
}

func findChromeInRegistry() (string, error) {
	// 初始化COM对象
	err := ole.CoInitialize(0)
	if err != nil {
		return "", err
	}
	defer ole.CoUninitialize()

	// 创建注册表对象
	unknown, err := oleutil.CreateObject("WScript.Shell")
	if err != nil {
		return "", err
	}
	defer unknown.Release()

	wShell, err := unknown.QueryInterface(ole.IID_IDispatch)
	if err != nil {
		return "", err
	}
	defer wShell.Release()

	reg, err := oleutil.CallMethod(wShell, "RegRead", "HKLM\\SOFTWARE\\Microsoft\\Windows\\CurrentVersion\\App Paths\\chrome.exe")
	if err != nil {
		return "", err
	}

	// 获取Chrome.exe所在路径
	chromePath := reg.ToString()

	// 校验路径有效性
	if _, err := os.Stat(chromePath); err != nil {
		return "", fmt.Errorf("chrome path from registry is invalid: %v", err)
	}

	// 返回完整路径
	return filepath.Abs(chromePath)
}

// findChromeOnMacOS looks for Chrome in the default Applications directory.
func findChromeOnMacOS() (string, error) {
	path := "/Applications/Google Chrome.app/Contents/MacOS/Google Chrome"
	if _, err := os.Stat(path); err == nil {
		return path, nil
	}
	return "", errors.New("chrome not found in default location")
}

// findChromeOnUnixLike searches for the Chrome executable in standard system directories.
func findChromeOnUnixLike() (string, error) {
	var paths []string
	if dirs := strings.Split(os.Getenv("PATH"), ":"); dirs != nil {
		paths = append(paths, dirs...)
	}

	// Add common installation directories for Unix-like systems.
	paths = append(paths,
		"/usr/bin",
		"/usr/local/bin",
		"/opt/google/chrome/bin",
	)

	var chromeExe = "google-chrome" // or "chrome" depending on your system

	for _, dir := range paths {
		if _, err := os.Stat(filepath.Join(dir, chromeExe)); err == nil {
			return filepath.Join(dir, chromeExe), nil
		}
	}

	return "", errors.New("chrome not found in standard system directories")
}
