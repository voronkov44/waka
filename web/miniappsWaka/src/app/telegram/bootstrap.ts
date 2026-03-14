import type { TelegramUser, TelegramWebApp } from '../types/telegram';

type InitDataSource = 'webapp-initData' | 'url-search' | 'url-hash' | 'none';
type UserSource = 'initDataUnsafe' | 'initData' | 'none';

export interface TelegramBootstrapDiagnostics {
  webAppDetected: boolean;
  initDataPresent: boolean;
  initDataSource: InitDataSource;
  userPresent: boolean;
  userSource: UserSource;
  attempts: number;
  readyCalled: boolean;
  expandCalled: boolean;
}

export interface TelegramBootstrapContext {
  webApp: TelegramWebApp | null;
  hasTelegramContext: boolean;
  initData: string;
  user: TelegramUser | null;
  diagnostics: TelegramBootstrapDiagnostics;
}

interface WaitOptions {
  maxAttempts?: number;
  delayMs?: number;
}

const INIT_DATA_PARAM_KEYS = ['tgWebAppData', 'initData'];

function sleep(ms: number): Promise<void> {
  return new Promise((resolve) => window.setTimeout(resolve, ms));
}

function logDebug(message: string, payload?: unknown) {
  if (!import.meta.env.DEV) {
    return;
  }

  if (payload !== undefined) {
    console.debug(`[waka:tma] ${message}`, payload);
    return;
  }
  console.debug(`[waka:tma] ${message}`);
}

function normalizeTelegramUser(raw: unknown): TelegramUser | null {
  if (!raw || typeof raw !== 'object' || Array.isArray(raw)) {
    return null;
  }

  const data = raw as Record<string, unknown>;
  const idValue = data.id;
  const id =
    typeof idValue === 'number'
      ? idValue
      : typeof idValue === 'string'
        ? Number(idValue)
        : Number.NaN;

  if (!Number.isFinite(id) || id <= 0) {
    return null;
  }

  return {
    id,
    first_name: typeof data.first_name === 'string' ? data.first_name : '',
    last_name: typeof data.last_name === 'string' ? data.last_name : undefined,
    username: typeof data.username === 'string' ? data.username : undefined,
    photo_url: typeof data.photo_url === 'string' ? data.photo_url : undefined,
    language_code: typeof data.language_code === 'string' ? data.language_code : undefined,
    is_premium: typeof data.is_premium === 'boolean' ? data.is_premium : undefined,
  };
}

function decodeMaybe(value: string): string {
  if (!value.includes('%')) {
    return value;
  }

  try {
    return decodeURIComponent(value);
  } catch {
    return value;
  }
}

function parseInitDataUser(initData: string): TelegramUser | null {
  if (!initData) {
    return null;
  }

  const normalized = initData.startsWith('?') ? initData.slice(1) : initData;
  const params = new URLSearchParams(normalized);
  const userRaw = params.get('user');
  if (!userRaw) {
    return null;
  }

  const attempts = [userRaw, decodeMaybe(userRaw)];
  for (const attempt of attempts) {
    try {
      return normalizeTelegramUser(JSON.parse(attempt));
    } catch {
      // keep trying fallbacks
    }
  }

  return null;
}

function readSearchParam(name: string): string {
  return new URLSearchParams(window.location.search).get(name)?.trim() ?? '';
}

function extractParamFromHash(name: string): string {
  const hash = window.location.hash.replace(/^#/, '');
  if (!hash) {
    return '';
  }

  // #key=value&... OR #/route?key=value&...
  const queryPart = hash.includes('?') ? hash.slice(hash.indexOf('?') + 1) : hash;
  const params = new URLSearchParams(queryPart);
  return params.get(name) ?? '';
}

function resolveInitData(webApp: TelegramWebApp | null): { initData: string; source: InitDataSource } {
  const fromWebApp = typeof webApp?.initData === 'string' ? webApp.initData.trim() : '';
  if (fromWebApp) {
    return { initData: fromWebApp, source: 'webapp-initData' };
  }

  for (const key of INIT_DATA_PARAM_KEYS) {
    const fromSearch = readSearchParam(key);
    if (fromSearch) {
      return { initData: decodeMaybe(fromSearch), source: 'url-search' };
    }
  }

  for (const key of INIT_DATA_PARAM_KEYS) {
    const fromHash = extractParamFromHash(key).trim();
    if (fromHash) {
      return { initData: decodeMaybe(fromHash), source: 'url-hash' };
    }
  }

  return { initData: '', source: 'none' };
}

function getTelegramWebApp(): TelegramWebApp | null {
  return window.Telegram?.WebApp ?? null;
}

export async function waitForTelegramWebApp(options: WaitOptions = {}): Promise<TelegramWebApp | null> {
  const maxAttempts = options.maxAttempts ?? 8;
  const delayMs = options.delayMs ?? 100;

  for (let attempt = 0; attempt < maxAttempts; attempt += 1) {
    const webApp = getTelegramWebApp();
    if (webApp) {
      return webApp;
    }
    await new Promise((resolve) => window.setTimeout(resolve, delayMs));
  }

  return getTelegramWebApp();
}

function markReadyAndExpand(webApp: TelegramWebApp | null): { readyCalled: boolean; expandCalled: boolean } {
  let readyCalled = false;
  let expandCalled = false;

  if (!webApp) {
    return { readyCalled, expandCalled };
  }

  try {
    webApp.ready?.();
    readyCalled = true;
  } catch (error) {
    logDebug('WebApp.ready() failed', error);
  }

  try {
    webApp.expand?.();
    expandCalled = true;
  } catch (error) {
    logDebug('WebApp.expand() failed', error);
  }

  return { readyCalled, expandCalled };
}

export async function ensureTelegramViewport(options: WaitOptions = {}) {
  const webApp = await waitForTelegramWebApp({
    maxAttempts: options.maxAttempts ?? 4,
    delayMs: options.delayMs ?? 120,
  });
  markReadyAndExpand(webApp);
}

export async function bootstrapTelegramContext(options: WaitOptions = {}): Promise<TelegramBootstrapContext> {
  const maxAttempts = options.maxAttempts ?? 20;
  const delayMs = options.delayMs ?? 120;
  const maxAttemptsWithoutContext = Math.min(maxAttempts, 8);

  let webApp: TelegramWebApp | null = null;
  let initData = '';
  let initDataSource: InitDataSource = 'none';
  let user: TelegramUser | null = null;
  let userSource: UserSource = 'none';
  let readyCalled = false;
  let expandCalled = false;
  let attempts = 0;

  for (let attempt = 1; attempt <= maxAttempts; attempt += 1) {
    attempts = attempt;
    webApp = getTelegramWebApp();

    if (webApp && (!readyCalled || !expandCalled)) {
      const viewportDiagnostics = markReadyAndExpand(webApp);
      readyCalled = readyCalled || viewportDiagnostics.readyCalled;
      expandCalled = expandCalled || viewportDiagnostics.expandCalled;
    }

    const initDataResult = resolveInitData(webApp);
    initData = initDataResult.initData;
    initDataSource = initDataResult.source;

    const initDataUnsafeUser = normalizeTelegramUser(webApp?.initDataUnsafe?.user);
    const parsedUser = parseInitDataUser(initData);
    user = initDataUnsafeUser ?? parsedUser;
    userSource = initDataUnsafeUser ? 'initDataUnsafe' : parsedUser ? 'initData' : 'none';

    const hasContextSignal = Boolean(webApp || initData);
    logDebug(`Telegram bootstrap attempt ${attempt}/${maxAttempts}`, {
      webAppDetected: Boolean(webApp),
      initDataPresent: Boolean(initData),
      userPresent: Boolean(user),
      initDataSource,
      userSource,
    });

    if (user) {
      break;
    }

    if (!hasContextSignal && attempt >= maxAttemptsWithoutContext) {
      break;
    }

    if (attempt < maxAttempts) {
      // Telegram launch data can appear shortly after initial render.
      await sleep(delayMs);
    }
  }

  const hasTelegramContext = Boolean(webApp || initData);
  const diagnostics: TelegramBootstrapDiagnostics = {
    webAppDetected: Boolean(webApp),
    initDataPresent: initData.length > 0,
    initDataSource,
    userPresent: Boolean(user),
    userSource,
    attempts,
    readyCalled,
    expandCalled,
  };

  logDebug('Telegram bootstrap diagnostics', diagnostics);

  return {
    webApp,
    hasTelegramContext,
    initData,
    user,
    diagnostics,
  };
}
