(async ({
  scriptURL,
  publicKey,
  cookieNames: cookieNameList,
  contextURL,
  identifier,
  castleTokenBatchSize,
}) => {
  const cookieNames = new Set(cookieNameList);
  const castleStorageKey = "fi.mau.twitter.castle_token";

  let browserStatusText;
  function showBrowserLoginStatus(message) {
    try {
      document.title = "Signing in to X";
      if (!browserStatusText || !browserStatusText.isConnected) {
        const body = document.body || document.documentElement.appendChild(document.createElement("body"));
        const container = document.createElement("main");
        container.id = "mautrix-twitter-login-status";
        container.setAttribute("role", "status");
        container.setAttribute("aria-live", "polite");

        const title = document.createElement("h1");
        title.textContent = "Signing in to X";
        Object.assign(title.style, {
          margin: "0",
          color: "#eff3f4",
          fontSize: "28px",
          fontWeight: "700",
          lineHeight: "1.2",
        });

        const progress = document.createElement("progress");
        progress.setAttribute("aria-label", "Signing in");
        Object.assign(progress.style, {
          width: "220px",
          height: "6px",
          margin: "28px 0 22px",
          accentColor: "#1d9bf0",
        });

        browserStatusText = document.createElement("p");
        Object.assign(browserStatusText.style, {
          margin: "0",
          color: "#8b98a5",
          fontSize: "15px",
          lineHeight: "1.5",
        });

        container.append(title, progress, browserStatusText);
        Object.assign(container.style, {
          width: "min(420px, calc(100% - 48px))",
          textAlign: "center",
        });
        Object.assign(document.documentElement.style, {
          minHeight: "100%",
          background: "#000000",
          colorScheme: "dark",
        });
        Object.assign(body.style, {
          minHeight: "100vh",
          margin: "0",
          display: "grid",
          placeItems: "center",
          background: "#000000",
          fontFamily: "-apple-system, BlinkMacSystemFont, 'Segoe UI', sans-serif",
        });
        body.replaceChildren(container);
      }
      browserStatusText.textContent = message;
    } catch (_) {}
  }
  showBrowserLoginStatus("Preparing secure login...");

  const browserLog = message => {
    try {
      console.log("[BrowserAuth] mautrix-twitter Castle: " + message);
    } catch (_) {}
  };
  function getLocalStorage() {
    try {
      return window.localStorage;
    } catch (_) {}
    return null;
  }
  function castleTokenResultKey(index) {
    return index === 1 ? "castle_token" : "castle_token_" + index;
  }
  function castleTokenStorageKeyForIndex(index) {
    return index === 1 ? castleStorageKey : castleStorageKey + "_" + index;
  }
  const sleep = ms => new Promise(resolve => setTimeout(resolve, ms));
  async function waitFor(predicate, message, timeout = 30000) {
    const deadlineAt = Date.now() + timeout;
    while (Date.now() < deadlineAt) {
      const value = predicate();
      if (value) {
        return value;
      }
      await sleep(100);
    }
    throw new Error(message);
  }
  function resetStoredResult() {
    try {
      const storage = getLocalStorage();
      if (!storage) {
        return;
      }
      storage.removeItem(castleStorageKey);
      for (let index = 2; index <= castleTokenBatchSize; index++) {
        storage.removeItem(castleTokenStorageKeyForIndex(index));
      }
      for (const name of cookieNames) {
        storage.removeItem("fi.mau.twitter.cookie." + name);
      }
    } catch (_) {}
  }
  function storeBrowserAuthResult(result) {
    try {
      globalThis.__MAUTRIX_TWITTER_CASTLE_RESULT__ = result;
      globalThis.__MAUTRIX_TWITTER_CASTLE_IN_PROGRESS__ = false;
      globalThis.__BEEP_BEEP_AUTH_RESULTS__ = result;
      if (typeof window !== "undefined") {
        window.__BEEP_BEEP_AUTH_RESULTS__ = result;
      }
    } catch (_) {}
    try {
      const storage = getLocalStorage();
      if (!storage) {
        return;
      }
      storage.setItem(castleStorageKey, result.castle_token || "");
      for (let index = 2; index <= castleTokenBatchSize; index++) {
        const key = castleTokenResultKey(index);
        if (result[key]) {
          storage.setItem(castleTokenStorageKeyForIndex(index), result[key]);
        }
      }
      for (const name of cookieNames) {
        if (result[name]) {
          storage.setItem("fi.mau.twitter.cookie." + name, result[name]);
        }
      }
    } catch (_) {}
  }
  function storedBrowserAuthResult() {
    try {
      const result = globalThis.__MAUTRIX_TWITTER_CASTLE_RESULT__;
      if (result && result.castle_token) {
        return result;
      }
    } catch (_) {}
    try {
      const storage = getLocalStorage();
      if (!storage) {
        return null;
      }
      const token = storage.getItem(castleStorageKey);
      if (!token) {
        return null;
      }
      const result = { castle_token: token };
      for (let index = 2; index <= castleTokenBatchSize; index++) {
        const token = storage.getItem(castleTokenStorageKeyForIndex(index));
        if (token) {
          result[castleTokenResultKey(index)] = token;
        }
      }
      for (const name of cookieNames) {
        const value = storage.getItem("fi.mau.twitter.cookie." + name);
        if (value) {
          result[name] = value;
        }
      }
      return result;
    } catch (_) {
      return null;
    }
  }
  const existingResult = storedBrowserAuthResult();
  if (existingResult) {
    browserLog("returning stored result");
    return existingResult;
  }
  if (globalThis.__MAUTRIX_TWITTER_CASTLE_IN_PROGRESS__) {
    browserLog("waiting for in-flight result");
    return await waitFor(() => storedBrowserAuthResult(), "X Castle token generation did not finish", 30000);
  }
  globalThis.__MAUTRIX_TWITTER_CASTLE_IN_PROGRESS__ = true;
  resetStoredResult();

  function addModules(entry, modules) {
    if (!entry || typeof entry !== "object") {
      return;
    }
    const defs = entry[1];
    if (!defs || typeof defs !== "object") {
      return;
    }
    for (const id of Object.keys(defs)) {
      modules[id] = defs[id];
    }
  }

  function loadScript(doc, url, timeout = 10000) {
    return new Promise((resolve, reject) => {
      const script = doc.createElement("script");
      const timer = setTimeout(() => {
        try {
          script.remove();
        } catch (_) {}
        reject(new Error("Timed out loading X Castle script"));
      }, timeout);
      script.src = url;
      script.async = true;
      script.onload = () => {
        clearTimeout(timer);
        resolve();
      };
      script.onerror = () => {
        clearTimeout(timer);
        reject(new Error("Failed to load X Castle script"));
      };
      (doc.head || doc.documentElement).appendChild(script);
    });
  }

  function installModuleCapture(win) {
    const chunk = win.webpackChunk_twitter_responsive_web = win.webpackChunk_twitter_responsive_web || [];
    const modules = {};
    for (const entry of chunk) {
      addModules(entry, modules);
    }
    if (!chunk.__mautrixCastleCaptured) {
      const nativePush = chunk.push.bind(chunk);
      chunk.push = (...entries) => {
        for (const entry of entries) {
          addModules(entry, modules);
        }
        return nativePush(...entries);
      };
      Object.defineProperty(chunk, "__mautrixCastleCaptured", { value: true });
    }
    return modules;
  }

  async function createTokenInContext(win, doc) {
    const modules = installModuleCapture(win);
    await loadScript(doc, scriptURL);
    await waitFor(() => Object.keys(modules).length > 0, "X Castle script did not register modules", 15000);

    const cache = {};
    function req(id) {
      const key = String(id);
      if (cache[key]) {
        return cache[key].exports;
      }
      const fn = modules[key];
      if (typeof fn !== "function") {
        throw new Error("Missing X Castle module " + key);
      }
      const module = { exports: {} };
      cache[key] = module;
      fn(module, module.exports, req);
      return module.exports;
    }
    req.d = (exports, definition) => {
      for (const key of Object.keys(definition)) {
        if (!Object.prototype.hasOwnProperty.call(exports, key)) {
          Object.defineProperty(exports, key, { enumerable: true, get: definition[key] });
        }
      }
    };
    req.o = (obj, prop) => Object.prototype.hasOwnProperty.call(obj, prop);
    req.r = exports => {
      if (typeof win.Symbol !== "undefined" && win.Symbol.toStringTag) {
        Object.defineProperty(exports, win.Symbol.toStringTag, { value: "Module" });
      }
      Object.defineProperty(exports, "__esModule", { value: true });
    };

    let configure;
    for (const id of Object.keys(modules)) {
      try {
        const exports = req(id);
        if (exports && typeof exports.configure === "function") {
          configure = exports.configure;
          break;
        }
      } catch (_) {}
    }
    if (typeof configure !== "function") {
      throw new Error("X Castle module is unavailable");
    }

    const castle = configure({ pk: publicKey });
    if (!castle || typeof castle.createRequestToken !== "function") {
      throw new Error("X Castle token generator is unavailable");
    }
    const tokens = [];
    for (let index = 0; index < castleTokenBatchSize; index++) {
      const token = await castle.createRequestToken();
      if (token) {
        tokens.push(token);
      }
      await sleep(10);
    }
    return tokens;
  }

  async function synthesizeCastleActivity(win, doc) {
    const target = doc.body || doc.documentElement;
    if (!target) {
      return;
    }
    const input = doc.createElement("input");
    input.type = "text";
    input.autocomplete = "username";
    input.style.cssText = "position:fixed;left:24px;top:24px;width:240px;height:32px;opacity:0;pointer-events:none;";
    target.appendChild(input);
    input.focus();
    const points = [[52, 48], [96, 52], [148, 60], [212, 68]];
    for (const [x, y] of points) {
      const eventInit = { bubbles: true, cancelable: true, clientX: x, clientY: y, screenX: x + 16, screenY: y + 88, pointerType: "mouse", isPrimary: true, buttons: 0 };
      if (typeof win.PointerEvent === "function") {
        target.dispatchEvent(new win.PointerEvent("pointermove", eventInit));
      } else {
        target.dispatchEvent(new win.MouseEvent("mousemove", eventInit));
      }
      await sleep(20);
    }
    const downInit = { bubbles: true, cancelable: true, clientX: 96, clientY: 52, screenX: 112, screenY: 140, pointerType: "mouse", isPrimary: true, buttons: 1 };
    if (typeof win.PointerEvent === "function") {
      input.dispatchEvent(new win.PointerEvent("pointerdown", downInit));
      input.dispatchEvent(new win.PointerEvent("pointerup", { ...downInit, buttons: 0 }));
    }
    input.dispatchEvent(new win.MouseEvent("mousedown", downInit));
    input.dispatchEvent(new win.MouseEvent("mouseup", { ...downInit, buttons: 0 }));
    input.dispatchEvent(new win.MouseEvent("click", { ...downInit, buttons: 0 }));
    const chars = identifier || "x";
    const valueSetter = Object.getOwnPropertyDescriptor(win.HTMLInputElement.prototype, "value").set;
    for (const ch of chars.slice(0, 16)) {
      input.dispatchEvent(new win.KeyboardEvent("keydown", { bubbles: true, cancelable: true, key: ch }));
      valueSetter.call(input, input.value + ch);
      input.dispatchEvent(new win.InputEvent("input", { bubbles: true, inputType: "insertText", data: ch }));
      input.dispatchEvent(new win.KeyboardEvent("keyup", { bubbles: true, cancelable: true, key: ch }));
      await sleep(15);
    }
    input.dispatchEvent(new win.Event("change", { bubbles: true }));
    await sleep(100);
  }

  async function loadFetchedContextFrame() {
    const resp = await fetch(contextURL, { credentials: "include" });
    if (!resp.ok) {
      throw new Error("Failed to fetch X login context: HTTP " + resp.status);
    }
    const html = await resp.text();
    const iframe = document.createElement("iframe");
    iframe.tabIndex = -1;
    iframe.style.cssText = "position:fixed;left:-10000px;top:-10000px;width:1024px;height:768px;border:0;opacity:0;pointer-events:none;";
    document.documentElement.appendChild(iframe);
    const win = iframe.contentWindow;
    const doc = iframe.contentDocument;
    if (!win || !doc) {
      throw new Error("Synthetic X login context is not accessible");
    }
    doc.open();
    doc.write(html);
    doc.close();
    try {
      win.history.replaceState(null, "", contextURL);
    } catch (_) {}
    await waitFor(() => doc.readyState === "interactive" || doc.readyState === "complete", "Synthetic X login context did not become ready", 15000);
    return { win, doc };
  }

  function copyCookies(result, cookieText) {
    for (const part of cookieText.split(";")) {
      const idx = part.indexOf("=");
      if (idx <= 0) {
        continue;
      }
      const name = part.slice(0, idx).trim();
      if (cookieNames.has(name)) {
        result[name] = part.slice(idx + 1).trim();
      }
    }
  }

  function quoteClientHint(value) {
    return '"' + String(value).replace(/\\/g, "\\\\").replace(/"/g, '\\"') + '"';
  }

  function captureBrowserHeaders() {
    const headers = {
      browser_user_agent: String(navigator.userAgent || ""),
    };
    const userAgentData = navigator.userAgentData;
    if (!userAgentData) {
      return headers;
    }
    const brands = Array.from(userAgentData.brands || []);
    if (brands.length > 0) {
      headers.browser_sec_ch_ua = brands.map(item =>
        quoteClientHint(item.brand) + ";v=" + quoteClientHint(item.version)
      ).join(", ");
    }
    if (userAgentData.platform) {
      headers.browser_sec_ch_ua_platform = quoteClientHint(userAgentData.platform);
    }
    headers.browser_sec_ch_ua_mobile = userAgentData.mobile ? "?1" : "?0";
    return headers;
  }

  try {
    browserLog("loading X context");
    const context = await loadFetchedContextFrame();
    browserLog("context ready");
    browserLog("creating module token");
    showBrowserLoginStatus("Completing sign-in...");
    await synthesizeCastleActivity(context.win, context.doc);
    const castleTokens = await createTokenInContext(context.win, context.doc);
    const castleToken = castleTokens[0] || "";
    browserLog("module tokens ready, count " + castleTokens.length + ", first length " + String(castleToken || "").length);
    const result = { castle_token: castleToken };
    Object.assign(result, captureBrowserHeaders());
    for (let index = 1; index < castleTokens.length; index++) {
      result[castleTokenResultKey(index + 1)] = castleTokens[index];
    }
    copyCookies(result, document.cookie);
    copyCookies(result, context.doc.cookie);
    storeBrowserAuthResult(result);
    browserLog("result stored");
    showBrowserLoginStatus("Finishing...");
    return result;
  } catch (err) {
    showBrowserLoginStatus("Unable to finish signing in.");
    try {
      globalThis.__MAUTRIX_TWITTER_CASTLE_IN_PROGRESS__ = false;
    } catch (_) {}
    throw err;
  }
})(__MAUTRIX_TWITTER_CASTLE_CONFIG__)
