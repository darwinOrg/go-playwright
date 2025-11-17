package extpw

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"syscall"
	"time"

	"github.com/go-ole/go-ole"
	"github.com/go-ole/go-ole/oleutil"
	"github.com/gorilla/websocket"
)

const (
	defaultCdpURL = "http://localhost:9222"
)

type BrowserInfo struct {
	WebSocketDebuggerURL string `json:"webSocketDebuggerUrl"`
}

func StartChrome(port int) (*exec.Cmd, string, error) {
	homeDir, _ := os.UserHomeDir()
	userDataDir := homeDir + "/ChromeProfile"

	if port == 0 {
		port = 9222
	}
	if isPortOpen("localhost", port) {
		return nil, userDataDir, nil
	}

	chrome, err := FindChromePath()
	if err != nil {
		return nil, "", err
	}

	cmd := exec.Command(chrome,
		fmt.Sprintf("--remote-debugging-port=%d", port),
		"--remote-allow-origins=*",
		"--remote-debugging-address=0.0.0.0", // 可选：允许外部访问
		"--no-first-run",
		"--no-default-browser-check",
		"--disable-gpu",
		"--disable-extensions",
		"--disable-plugins",
		"--disable-sync",
		"--user-data-dir="+userDataDir,
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
		return nil, "", err
	}

	log.Printf("Started Chrome with PID: %d", cmd.Process.Pid)
	return cmd, userDataDir, nil
}

func ShutdownChrome(cmd *exec.Cmd, baseURL string) {
	if baseURL == "" {
		baseURL = defaultCdpURL
	}

	// WebSocket Browser.close
	if err := closeChromeViaCDP(baseURL); err == nil {
		log.Println("通过 Browser.close 成功关闭")
		return
	}

	// 最终方案: SIGTERM
	log.Println("CDP 关闭失败，回退到 SIGTERM...")
	gracefulShutdown(cmd)
}

// closeChromeViaCDP 使用 WebSocket 发送 Browser.close
func closeChromeViaCDP(baseURL string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cdpWebSocketURL, err := getBrowserWebSocketURL(baseURL)
	if err != nil {
		return err
	}

	conn, _, err := websocket.DefaultDialer.DialContext(ctx, cdpWebSocketURL, nil)
	if err != nil {
		log.Printf(fmt.Sprintf("websocket.DefaultDialer.DialContext error: %v", err))
		return err
	}
	defer func() {
		_ = conn.Close()
	}()

	// 发送 Browser.close 命令
	closeCmd := `{
		"id": 1,
		"method": "Browser.close"
	}`

	if err := conn.WriteMessage(websocket.TextMessage, []byte(closeCmd)); err != nil {
		log.Printf(fmt.Sprintf("conn.WriteMessage error: %v", err))
		return err
	}

	log.Println("已发送 Browser.close 命令")
	return nil
}

func getBrowserWebSocketURL(baseURL string) (string, error) {
	resp, err := http.Get(baseURL + "/json/version")
	if err != nil {
		log.Printf(fmt.Sprintf(baseURL+"/json/version get error: %v", err))
		return "", err
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	var info BrowserInfo
	if err := json.NewDecoder(resp.Body).Decode(&info); err != nil {
		log.Printf(fmt.Sprintf("json.NewDecoder(resp.Body).Decode error: %v", err))
		return "", err
	}

	return info.WebSocketDebuggerURL, nil
}

func gracefulShutdown(cmd *exec.Cmd) {
	if cmd.Process == nil {
		log.Println("Chrome 进程不存在")
		return
	}

	log.Println("发送 SIGTERM 以优雅关闭 Chrome...")
	if runtime.GOOS == "windows" {
		// Windows 不支持 SIGTERM，尝试其他方式或直接 Kill
		_ = cmd.Process.Kill()
	} else {
		_ = cmd.Process.Signal(syscall.SIGTERM)
	}

	// 等待退出
	done := make(chan struct{})
	go func() {
		_ = cmd.Wait()
		close(done)
	}()

	select {
	case <-done:
		log.Println("Chrome 已优雅退出")
	case <-time.After(10 * time.Second):
		log.Println("超时，强制终止")
		_ = cmd.Process.Kill()
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

// isPortOpen 检查本地指定端口是否已经开启（监听）
func isPortOpen(host string, port int) bool {
	conn, err := net.DialTimeout("tcp", fmt.Sprintf("%s:%d", host, port), 2*time.Second)
	if err != nil {
		return false
	}
	defer func() {
		_ = conn.Close()
	}()
	return true
}
