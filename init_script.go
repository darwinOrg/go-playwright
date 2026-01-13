package extpw

var InitScript = `
// 彻底清理 cdc_ 变量 + 隐藏 webdriver
(() => {
  // 删除 webdriver
  delete navigator.__proto__.webdriver;

  // 删除所有 cdc_ 开头的属性
  const cdcKeys = Object.getOwnPropertyNames(window).filter(k => k.startsWith('cdc_'));
  cdcKeys.forEach(k => delete window[k]);

  Object.defineProperty(window, 'chrome', {
    writable: true,
    value: {
      ...window.chrome,
      runtime: undefined,
    },
  });

  // 隐藏 permissions policy 异常
  const originalPermissions = window.navigator.permissions;
  window.navigator.permissions = {
    ...originalPermissions,
    query: new Proxy(originalPermissions.query, {
      apply(target, thisArg, args) {
        return Promise.resolve({ state: 'granted' });
      }
    })
  };
})();
`
