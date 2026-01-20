package extpw

// GetInitScript 获取初始化脚本
func GetInitScript(osType OSType) string {
	var commonFonts string
	var commonDevices string

	if osType == OSMacOS {
		commonFonts = `
  const commonFonts = [
    'Arial', 'Helvetica', 'Times New Roman', 'Courier New', 'Georgia',
    'Verdana', 'Tahoma', 'Trebuchet MS', 'Lucida Grande', 'Palatino',
    'Menlo', 'Monaco', 'Consolas', 'Andale Mono', 'Charcoal', 'Geneva'
  ];
`
		commonDevices = `
  const commonDevices = [
    { deviceId: 'mac-audio-input-1', label: 'Built-in Microphone', kind: 'audioinput' },
    { deviceId: 'mac-audio-input-2', label: 'External Microphone', kind: 'audioinput' },
    { deviceId: 'mac-audio-output-1', label: 'Built-in Speakers', kind: 'audiooutput' }
  ];
`
	} else {
		commonFonts = `
  const commonFonts = [
    'Arial', 'Calibri', 'Cambria', 'Courier New', 'Georgia',
    'Microsoft YaHei', 'SimSun', 'SimHei', 'Times New Roman', 'Verdana',
    'Tahoma', 'Trebuchet MS', 'Lucida Console', 'Segoe UI', 'Comic Sans MS'
  ];
`
		commonDevices = `
  const commonDevices = [
    { deviceId: 'b0a2b8c3-45f9-4d42-b748-d925e4f8f302', label: 'Realtek High Definition Audio', kind: 'audioinput' },
    { deviceId: 'd3c39e1f-bf6c-4e8b-87b4-e1a3e45db602', label: 'Intel(R) Audio Display', kind: 'audioinput' },
    { deviceId: 'a1e7b1ac-5f0b-4d8d-8ff3-b08eaf08a4d4', label: 'Generic Audio Device', kind: 'audiooutput' }
  ];
`
	}

	return `
// 彻底清理 cdc_ 变量 + 隐藏 webdriver
(() => {
  // 删除 webdriver
  delete navigator.__proto__.webdriver;
  delete Object.getPrototypeOf(navigator).webdriver;

  // 删除所有 cdc_ 开头的属性
  const cdcKeys = Object.getOwnPropertyNames(window).filter(k => k.startsWith('cdc_'));
  cdcKeys.forEach(k => delete window[k]);

  // 删除所有 __playwright 开头的属性
  const playwrightKeys = Object.getOwnPropertyNames(window).filter(k => k.startsWith('__playwright'));
  playwrightKeys.forEach(k => delete window[k]);

  // 删除所有 __pw 开头的属性
  const pwKeys = Object.getOwnPropertyNames(window).filter(k => k.startsWith('__pw'));
  pwKeys.forEach(k => delete window[k]);

  // 伪装 chrome 对象
  Object.defineProperty(window, 'chrome', {
    writable: true,
    value: {
      ...window.chrome,
      runtime: undefined,
      loadTimes: function() {
        return {
          requestTime: Date.now() / 1000 - 0.5,
          startLoadTime: Date.now() / 1000 - 0.4,
          commitLoadTime: Date.now() / 1000 - 0.3,
          finishDocumentLoadTime: Date.now() / 1000 - 0.2,
          finishLoadTime: Date.now() / 1000 - 0.1,
          firstPaintTime: Date.now() / 1000 - 0.15,
          firstPaintAfterLoadTime: 0,
          navigationType: 'Other',
        };
      },
      csi: function() {
        return {
          startE: Date.now(),
          onloadT: Date.now() - 100,
          pageT: Date.now() - 200,
          tran: 15,
        };
      },
      app: {
        isInstalled: false,
        InstallState: { DISABLED: 'disabled', INSTALLED: 'installed', NOT_INSTALLED: 'not_installed' },
        RunningState: { CANNOT_RUN: 'cannot_run', READY_TO_RUN: 'ready_to_run', RUNNING: 'running' },
      },
    },
  });

  // 伪装 navigator.plugins
  const commonPlugins = [
    { name: 'Chrome PDF Plugin', description: 'Portable Document Format', filename: 'internal-pdf-viewer' },
    { name: 'Chrome PDF Viewer', description: '', filename: 'mhjfbmdgcfjbbpaeojofohoefgiehjai' },
    { name: 'Native Client', description: '', filename: 'internal-nacl-plugin' },
  ];

  Object.defineProperty(navigator, 'plugins', {
    get: function() {
      const plugins = [];
      commonPlugins.forEach(plugin => {
        const pluginObj = {
          name: plugin.name,
          description: plugin.description,
          filename: plugin.filename,
          length: 0,
        };
        plugins.push(pluginObj);
      });
      return plugins;
    },
    configurable: true
  });

  // 伪装 navigator.languages
  Object.defineProperty(navigator, 'languages', {
    get: function() {
      return ['zh-CN', 'zh', 'en-US', 'en'];
    },
    configurable: true
  });

  // 伪装 navigator.hardwareConcurrency
  Object.defineProperty(navigator, 'hardwareConcurrency', {
    get: function() {
      return 4;
    },
    configurable: true
  });

  // 伪装 navigator.deviceMemory
  Object.defineProperty(navigator, 'deviceMemory', {
    get: function() {
      return 8;
    },
    configurable: true
  });

  // 伪装 navigator.maxTouchPoints
  Object.defineProperty(navigator, 'maxTouchPoints', {
    get: function() {
      return 0;
    },
    configurable: true
  });

  // 伪装 navigator.connection
  if (!navigator.connection) {
    Object.defineProperty(navigator, 'connection', {
      value: {
        effectiveType: '4g',
        rtt: 50,
        downlink: 10,
        saveData: false,
      },
      configurable: true
    });
  }

  // 伪装 navigator.permissions
  const originalPermissions = window.navigator.permissions;
  window.navigator.permissions = {
    ...originalPermissions,
    query: new Proxy(originalPermissions.query, {
      apply(target, thisArg, args) {
        return Promise.resolve({ state: 'granted' });
      }
    })
  };

  // 伪装 navigator.getBattery
  if (navigator.getBattery) {
    const originalGetBattery = navigator.getBattery;
    navigator.getBattery = function() {
      return Promise.resolve({
        charging: true,
        chargingTime: 0,
        dischargingTime: Infinity,
        level: 1.0,
      });
    };
  }

  // 伪装 navigator.mediaCapabilities
  if (navigator.mediaCapabilities) {
    const originalMediaCapabilities = navigator.mediaCapabilities;
    navigator.mediaCapabilities = {
      ...originalMediaCapabilities,
      decodingInfo: function(config) {
        return Promise.resolve({
          supported: true,
          smooth: true,
          powerEfficient: true,
        });
      },
      encodingInfo: function(config) {
        return Promise.resolve({
          supported: true,
          smooth: true,
          powerEfficient: true,
        });
      },
    };
  }

  // WebGL Canvas 伪装
` + commonFonts + `
` + commonDevices + `
  const commonSampleRates = [44100, 48000, 96000];

  function delayExecution(fn, delay = 500) {
    setTimeout(fn, delay);
  }

  function generateAudioDevice() {
    const device = commonDevices[Math.floor(Math.random() * commonDevices.length)];
    return new Proxy(device, {
      get(target, prop) {
        if (prop === 'toJSON') {
          return () => target;
        }
        return target[prop];
      }
    });
  }

  function getRandomFont() {
    return commonFonts[Math.floor(Math.random() * commonFonts.length)];
  }

  // 工具函数：生成高斯噪声
  function getGaussianNoise(mean, stddev) {
    let u = 1 - Math.random();
    let v = 1 - Math.random();
    let noise = Math.sqrt(-2.0 * Math.log(u)) * Math.cos(2.0 * Math.PI * v);
    return mean + noise * stddev;
  }

  let cachedNoise = null;
  function getCachedNoise(width, height) {
    if (!cachedNoise) {
      cachedNoise = Array.from({ length: width * height }, () => getGaussianNoise(0, 5));
    }
    return cachedNoise;
  }

  function getRandomFakeIP() {
    return 192.168.${Math.floor(Math.random() * 256)}.${Math.floor(Math.random() * 256)};
  }

  // WebGL Proxy 处理
  const WebGLProxyHandler = {
    get(target, prop, receiver) {
      if (prop === 'getSupportedExtensions') {
        return function() {
          return [
            'ANGLE_instanced_arrays', 'EXT_blend_minmax', 'EXT_clip_control',
            'EXT_color_buffer_half_float', 'EXT_depth_clamp', 'EXT_disjoint_timer_query',
            'EXT_float_blend', 'EXT_frag_depth', 'EXT_polygon_offset_clamp',
            'EXT_shader_texture_lod', 'EXT_texture_compression_bptc', 'EXT_texture_compression_rgtc',
            'EXT_texture_filter_anisotropic', 'EXT_texture_mirror_clamp_to_edge', 'EXT_sRGB',
            'KHR_parallel_shader_compile', 'OES_element_index_uint', 'OES_fbo_render_mipmap',
            'OES_standard_derivatives', 'OES_texture_float', 'OES_texture_float_linear',
            'OES_texture_half_float', 'OES_texture_half_float_linear', 'OES_vertex_array_object',
            'WEBGL_blend_func_extended', 'WEBGL_color_buffer_float', 'WEBGL_compressed_texture_s3tc',
            'WEBGL_compressed_texture_s3tc_srgb', 'WEBGL_debug_renderer_info', 'WEBGL_debug_shaders',
            'WEBGL_depth_texture', 'WEBGL_draw_buffers', 'WEBGL_lose_context', 'WEBGL_multi_draw',
            'WEBGL_polygon_mode', 'WEBGL_compressed_texture_etc', 'WEBGL_compressed_texture_etc1',
            'WEBGL_compressed_texture_pvrtc', 'WEBGL_compressed_texture_atc', 'WEBGL_compressed_texture_astc'
          ];
        };
      }

      if (prop === 'getParameter') {
        return function(parameter) {
          const UNMASKED_RENDERER_WEBGL = 0x9246;
          const UNMASKED_VENDOR_WEBGL = 0x9245;
          const UNMASKED_RENDERER_WEBGL2 = 0x9332;
          const UNMASKED_VENDOR_WEBGL2 = 0x9331;

          const fakeParameters = {
            [target.MAX_TEXTURE_SIZE]: 16384,
            [target.MAX_VERTEX_ATTRIBS]: 16,
            [target.MAX_TEXTURE_IMAGE_UNITS]: 16,
            [target.MAX_COMBINED_TEXTURE_IMAGE_UNITS]: 32,
            [target.MAX_VERTEX_UNIFORM_VECTORS]: 4096,
            [target.MAX_FRAGMENT_UNIFORM_VECTORS]: 1024,
            [target.MAX_VARYING_VECTORS]: 30,
            [target.MAX_RENDERBUFFER_SIZE]: 16384,
            [target.MAX_VIEWPORT_DIMS]: [32767, 32767],
            [target.SAMPLES]: 4,
            [target.SAMPLE_BUFFERS]: 1,
            [UNMASKED_RENDERER_WEBGL]: 'ANGLE (Intel(R) UHD Graphics 630 Direct3D11 vs_5_0 ps_5_0)',
            [UNMASKED_VENDOR_WEBGL]: 'Intel Inc.',
            [UNMASKED_RENDERER_WEBGL2]: 'ANGLE (Intel(R) UHD Graphics 630 Direct3D11 vs_5_0 ps_5_0)',
            [UNMASKED_VENDOR_WEBGL2]: 'Intel Inc.',
            [target.ALIASED_LINE_WIDTH_RANGE]: [1, 1],
            [target.ALIASED_POINT_SIZE_RANGE]: [1, 1024],
            [target.MAX_CUBE_MAP_TEXTURE_SIZE]: 16384,
            [target.MAX_RENDERBUFFER_SIZE]: 16384,
            [target.MAX_3D_TEXTURE_SIZE]: 2048,
            [target.MAX_ARRAY_TEXTURE_LAYERS]: 2048,
          };
          return fakeParameters[parameter] || Reflect.get(target, prop, receiver).call(target, parameter);
        };
      }

      if (prop === 'getExtension') {
        return function(name) {
          if (name === 'WEBGL_debug_renderer_info') {
            return {
              UNMASKED_RENDERER_WEBGL: 0x9246,
              UNMASKED_VENDOR_WEBGL: 0x9245,
            };
          }
          return Reflect.get(target, prop, receiver).call(target, name);
        };
      }

      return Reflect.get(target, prop, receiver);
    },
    has(target, prop) {
      return prop in target; // 确保属性检查一致性
    }
  };

  // 修改 Canvas
  const originalGetContext = HTMLCanvasElement.prototype.getContext;
  const originalDescriptor = Object.getOwnPropertyDescriptor(HTMLCanvasElement.prototype, 'getContext');
  Object.defineProperty(HTMLCanvasElement.prototype, 'getContext', {
    ...originalDescriptor,
    value: function (type, ...args) {
      const context = originalGetContext.apply(this, [type, ...args]);

      if (type === '2d' && context) {
        const originalGetImageData = context.getImageData;
        context.getImageData = function (x, y, width, height) {
          const imageData = originalGetImageData.call(this, x, y, width, height);
          const noise = getCachedNoise(width, height);

          for (let i = 0; i < imageData.data.length; i += 4) {
            imageData.data[i] += noise[i];       // 红色通道
            imageData.data[i + 1] += noise[i + 1]; // 绿色通道
            imageData.data[i + 2] += noise[i + 2]; // 蓝色通道
          }
          return imageData;
        };
      }

      return (type === 'webgl' || type === 'webgl2')
        ? new Proxy(context, WebGLProxyHandler)
        : context;
    }
  });

  const originalToDataURL = HTMLCanvasElement.prototype.toDataURL;
  HTMLCanvasElement.prototype.toDataURL = function(type, quality) {
    const context = this.getContext('2d');
    if (context) {
      const imageData = context.getImageData(0, 0, this.width, this.height);
      const noise = getCachedNoise(this.width, this.height);

      for (let i = 0; i < imageData.data.length; i += 4) {
        imageData.data[i] += noise[i];
        imageData.data[i + 1] += noise[i + 1];
        imageData.data[i + 2] += noise[i + 2];
      }
      context.putImageData(imageData, 0, 0);
    }
    return originalToDataURL.call(this, type, quality);
  };

  // 代理 document.fonts.check
  const originalFontCheck = document.fonts.check;
  document.fonts.check = new Proxy(originalFontCheck, {
    apply(target, thisArg, args) {
      if (args.length > 0) {
        // 提取字体家族名称
        const fontDescriptor = args[0];
        const fontFamilyMatches = fontDescriptor.match(/"(.*?)"/) || fontDescriptor.match(/'([^']*)'/);
        const fontFamily = fontFamilyMatches ? fontFamilyMatches[1] : fontDescriptor.split(/\\s+/).pop();

        // 伪装返回 true 如果字体在伪装列表中
        if (commonFonts.includes(fontFamily)) {
          return true;
        }
      }
      // 默认行为
      return target.apply(thisArg, args);
    }
  });

  // 代理 CSSStyleDeclaration.prototype.setProperty
  const originalSetProperty = CSSStyleDeclaration.prototype.setProperty;
  CSSStyleDeclaration.prototype.setProperty = new Proxy(originalSetProperty, {
    apply(target, thisArg, args) {
      if (args[0] === 'font-family' && args[1]) {
        const requestedFonts = args[1].split(',').map(f => f.trim());
        const modifiedFonts = requestedFonts.map(f => commonFonts.includes(f) ? f : getRandomFont());
        return target.apply(thisArg, [args[0], modifiedFonts.join(', '), args[2]]);
      }
      return target.apply(thisArg, args);
    }
  });

  // 代理 window.getComputedStyle
  const originalGetComputedStyle = window.getComputedStyle;
  window.getComputedStyle = new Proxy(originalGetComputedStyle, {
    apply(target, thisArg, args) {
      const style = target.apply(thisArg, args);
      const originalFontFamily = style.fontFamily;

      Object.defineProperty(style, 'fontFamily', {
        get() {
          const fonts = originalFontFamily.split(',').map(f => f.trim());
          return fonts.map(f => commonFonts.includes(f) ? f : getRandomFont()).join(', ');
        },
        configurable: true
      });

      return style;
    }
  });

  // 代理 FontFace.load
  FontFace.prototype.load = new Proxy(FontFace.prototype.load, {
    apply(target, thisArg) {
      return Promise.resolve(thisArg);
    }
  });

  // 代理 FontFace.prototype.family
  Object.defineProperty(FontFace.prototype, 'family', {
    get() {
      return commonFonts.includes(this._family) ? this._family : getRandomFont();
    },
    set(value) {
      this._family = value;
    },
    configurable: true
  });

  // 记录真实 FontFace.family
  FontFace.prototype._family = '';

  // 代理 navigator.mediaDevices.enumerateDevices
  const originalEnumerateDevices = navigator.mediaDevices.enumerateDevices;
  navigator.mediaDevices.enumerateDevices = function() {
    return new Promise((resolve) => {
      delayExecution(() => {
        const devices = Array.from({ length: 3 }, generateAudioDevice);
        resolve(devices);
      });
    });
  };

  // 代理 window.AudioContext
  const originalAudioContext = window.AudioContext;
  window.AudioContext = new Proxy(originalAudioContext, {
    construct(target, args) {
      const context = new target(...args);
      const sampleRate = commonSampleRates[Math.floor(Math.random() * commonSampleRates.length)];
      Object.defineProperty(context, 'sampleRate', {
        get: () => sampleRate,
        configurable: true
      });
      return context;
    }
  });

  // 代理 navigator.mediaDevices.getUserMedia
  const originalGetUserMedia = navigator.mediaDevices.getUserMedia;
  navigator.mediaDevices.getUserMedia = function(constraints) {
    return new Promise((resolve) => {
      delayExecution(() => {
        const fakeStream = new MediaStream();
        resolve(fakeStream);
      });
    });
  };

  // 强制 WebRTC 配置使用 IP 地址
  const originalRTCPeerConnection = window.RTCPeerConnection;
  const peerConnectionHandlingMap = new WeakMap();

  window.RTCPeerConnection = function (...args) {
    const pc = new originalRTCPeerConnection(...args);
    peerConnectionHandlingMap.set(pc, false);

    pc.addEventListener("icecandidate", (event) => {
      if (!event.candidate) return;

      const isHandlingEvent = peerConnectionHandlingMap.get(pc);
      if (isHandlingEvent) return;

      peerConnectionHandlingMap.set(pc, true);

      setTimeout(() => {
        const candidate = event.candidate.candidate;

        const ipRegex = /(?:\\d{1,3}\\.){3}\\d{1,3}|([a-zA-Z0-9-]+(?:\\.[a-zA-Z0-9-]+)*\\.local)|(?:[a-fA-F0-9:]+:+)+[a-fA-F0-9]+/g;
        const fakeCandidate = candidate.replace(ipRegex, () => getRandomFakeIP());

        const modifiedEvent = new Event("icecandidate");
        Object.defineProperty(modifiedEvent, "candidate", {
          value: {
            ...event.candidate,
            candidate: fakeCandidate,
          },
          writable: false,
          enumerable: true,
          configurable: false,
        });
        pc.dispatchEvent(modifiedEvent);

        peerConnectionHandlingMap.set(pc, false);
      }, 0);
    });

    return pc;
  };

  // 代理 Function.prototype.toString 以隐藏代理痕迹
  Function.prototype.toString = new Proxy(Function.prototype.toString, {
    apply(target, thisArg, args) {
      if (thisArg === HTMLCanvasElement.prototype.getContext) {
        return "function getContext() { [native code] }";
      }
      if (thisArg === HTMLCanvasElement.prototype.toDataURL) {
        return "function toDataURL() { [native code] }";
      }
      if (thisArg === window.RTCPeerConnection) {
        return "function RTCPeerConnection() { [native code] }";
      }
      return target.apply(thisArg, args);
    },
  });

  // 伪装屏幕方向
  if (!screen.orientation) {
    Object.defineProperty(screen, 'orientation', {
      value: {
        type: 'landscape-primary',
        angle: 0,
        lock: function() {},
        unlock: function() {},
      },
      configurable: true
    });
  }

  // 伪装 window.outerWidth 和 window.outerHeight
  Object.defineProperty(window, 'outerWidth', {
    get: function() {
      return 1920;
    },
    configurable: true
  });

  Object.defineProperty(window, 'outerHeight', {
    get: function() {
      return 1080;
    },
    configurable: true
  });

  // 伪装 window.screenX 和 window.screenY
  Object.defineProperty(window, 'screenX', {
    get: function() {
      return 0;
    },
    configurable: true
  });

  Object.defineProperty(window, 'screenY', {
    get: function() {
      return 0;
    },
    configurable: true
  });

  // 伪装 window.devicePixelRatio
  Object.defineProperty(window, 'devicePixelRatio', {
    get: function() {
      return 1.0;
    },
    configurable: true
  });

  // 伪装 screen.colorDepth 和 screen.pixelDepth
  Object.defineProperty(screen, 'colorDepth', {
    get: function() {
      return 24;
    },
    configurable: true
  });

  Object.defineProperty(screen, 'pixelDepth', {
    get: function() {
      return 24;
    },
    configurable: true
  });

  // 伪装 screen.orientation
  if (!screen.orientation) {
    Object.defineProperty(screen, 'orientation', {
      value: {
        type: 'landscape-primary',
        angle: 0,
        lock: function() {},
        unlock: function() {},
      },
      configurable: true
    });
  }

  // 伪装 Date.prototype.getTimezoneOffset
  const originalGetTimezoneOffset = Date.prototype.getTimezoneOffset;
  Date.prototype.getTimezoneOffset = function() {
    return -480; // UTC+8
  };

  // 伪装 Intl.DateTimeFormat
  const originalDateTimeFormat = Intl.DateTimeFormat;
  Intl.DateTimeFormat = function() {
    const formatter = new originalDateTimeFormat(...arguments);
    const originalResolvedOptions = formatter.resolvedOptions;
    formatter.resolvedOptions = function() {
      const options = originalResolvedOptions.call(this);
      options.timeZone = 'Asia/Shanghai';
      return options;
    };
    return formatter;
  };
})();
`
}

// InitScript 默认初始化脚本（向后兼容）
var InitScript = GetInitScript(GetCurrentOS())
