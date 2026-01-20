package extpw

// CheatConfig 防爬虫配置
type CheatConfig struct {
	OSType OSType
}

// PlatformConfig 平台配置
type PlatformConfig struct {
	UserAgent         string
	Platform          string
	AcceptLanguage    string
	UserAgentMetadata UserAgentMetadata
	SecChUaPlatform   string
}

// UserAgentMetadata User Agent 元数据
type UserAgentMetadata struct {
	Brands          []Brand
	FullVersion     string
	Platform        string
	PlatformVersion string
	Architecture    string
	Model           string
	Mobile          bool
}

// Brand User Agent 品牌
type Brand struct {
	Brand   string
	Version string
}

// ScreenConfig 屏幕配置
type ScreenConfig struct {
	Width             int
	Height            int
	DeviceScaleFactor float64
	Mobile            bool
	ScreenWidth       int
	ScreenHeight      int
	PositionX         int
	PositionY         int
	ViewportX         int
	ViewportY         int
	ViewportWidth     int
	ViewportHeight    int
}

// GetPlatformConfig 获取平台配置
func GetPlatformConfig(osType OSType) PlatformConfig {
	if osType == OSMacOS {
		return PlatformConfig{
			UserAgent:      "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/119.0.6045.91 Safari/537.36",
			Platform:       "MacIntel",
			AcceptLanguage: "zh-CN",
			UserAgentMetadata: UserAgentMetadata{
				Brands: []Brand{
					{Brand: "Google Chrome", Version: "119"},
					{Brand: "Chromium", Version: "119"},
					{Brand: "Not A Brand", Version: "24"},
				},
				FullVersion:     "119.0.6045.91",
				Platform:        "macOS",
				PlatformVersion: "10.15.7",
				Architecture:    "x86_64",
				Model:           "",
				Mobile:          false,
			},
			SecChUaPlatform: "\"macOS\"",
		}
	}

	// 默认返回 Windows 配置
	return PlatformConfig{
		UserAgent:      "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/119.0.6045.91 Safari/537.36",
		Platform:       "Win32",
		AcceptLanguage: "zh-CN",
		UserAgentMetadata: UserAgentMetadata{
			Brands: []Brand{
				{Brand: "Google Chrome", Version: "119"},
				{Brand: "Chromium", Version: "119"},
				{Brand: "Not A Brand", Version: "24"},
			},
			FullVersion:     "119.0.6045.91",
			Platform:        "Windows",
			PlatformVersion: "10.0",
			Architecture:    "x86_64",
			Model:           "",
			Mobile:          false,
		},
		SecChUaPlatform: "\"Windows\"",
	}
}

// GetScreenConfig 获取屏幕配置
func GetScreenConfig(osType OSType) ScreenConfig {
	if osType == OSMacOS {
		return ScreenConfig{
			Width:             1920,
			Height:            1080,
			DeviceScaleFactor: 2.0, // macOS 通常使用 2.0 的设备缩放因子
			Mobile:            false,
			ScreenWidth:       1920,
			ScreenHeight:      1080,
			PositionX:         0,
			PositionY:         25, // macOS 菜单栏高度
			ViewportX:         0,
			ViewportY:         25,
			ViewportWidth:     1920,
			ViewportHeight:    1055,
		}
	}

	// 默认返回 Windows 配置
	return ScreenConfig{
		Width:             1920,
		Height:            1040,
		DeviceScaleFactor: 1.0,
		Mobile:            false,
		ScreenWidth:       1920,
		ScreenHeight:      1080,
		PositionX:         0,
		PositionY:         40, // Windows 任务栏高度
		ViewportX:         8,
		ViewportY:         48,
		ViewportWidth:     1904,
		ViewportHeight:    1024,
	}
}
