import {
  createContext,
  useCallback,
  useContext,
  useEffect,
  useMemo,
  useState,
  type ReactNode,
} from 'react';
import { en } from './locales/en';
import { ru } from './locales/ru';

const dictionaries = {
  en,
  ru,
} as const;

const DEFAULT_LOCALE = 'en' as const;
const LOCALE_STORAGE_KEY = 'waka_locale';
const LOCALE_EVENT = 'waka_locale_change';
const I18N_TEXT_PREFIX = '__i18n__:';

export type Locale = keyof typeof dictionaries;
export type TranslateParams = Record<string, string | number>;
export type TranslateFn = (key: string, params?: TranslateParams) => string;

type PluralCategory = 'zero' | 'one' | 'two' | 'few' | 'many' | 'other';

interface I18nContextValue {
  locale: Locale;
  localeCode: string;
  setLocale: (locale: Locale) => void;
  t: TranslateFn;
  tp: (baseKey: string, count: number, params?: TranslateParams) => string;
}

const localeCodeMap: Record<Locale, string> = {
  en: 'en-US',
  ru: 'ru-RU',
};

const pluralRules: Record<Locale, Intl.PluralRules> = {
  en: new Intl.PluralRules('en-US'),
  ru: new Intl.PluralRules('ru-RU'),
};

const I18nContext = createContext<I18nContextValue | null>(null);

function isObject(value: unknown): value is Record<string, unknown> {
  return Boolean(value) && typeof value === 'object' && !Array.isArray(value);
}

function isLocale(value: unknown): value is Locale {
  return value === 'en' || value === 'ru';
}

function getValueByPath(dictionary: Record<string, unknown>, path: string): unknown {
  return path.split('.').reduce<unknown>((acc, key) => {
    if (!isObject(acc)) {
      return undefined;
    }
    return acc[key];
  }, dictionary);
}

function interpolate(template: string, params?: TranslateParams): string {
  if (!params) {
    return template;
  }

  return template.replace(/\{(\w+)\}/g, (_, key: string) => {
    const value = params[key];
    return value === undefined ? `{${key}}` : String(value);
  });
}

export function translate(locale: Locale, key: string, params?: TranslateParams): string {
  const fromLocale = getValueByPath(dictionaries[locale] as unknown as Record<string, unknown>, key);
  if (typeof fromLocale === 'string') {
    return interpolate(fromLocale, params);
  }

  const fromDefault = getValueByPath(dictionaries[DEFAULT_LOCALE] as unknown as Record<string, unknown>, key);
  if (typeof fromDefault === 'string') {
    return interpolate(fromDefault, params);
  }

  return key;
}

function getStoredLocale(): Locale | null {
  if (typeof window === 'undefined') {
    return null;
  }

  try {
    const storedLocale = localStorage.getItem(LOCALE_STORAGE_KEY);
    return isLocale(storedLocale) ? storedLocale : null;
  } catch {
    return null;
  }
}

function resolveLocaleFromString(value: string | undefined | null): Locale | null {
  if (!value) {
    return null;
  }

  const normalized = value.trim().toLowerCase();
  if (!normalized) {
    return null;
  }

  if (normalized.startsWith('ru')) {
    return 'ru';
  }

  if (normalized.startsWith('en')) {
    return 'en';
  }

  return null;
}

function resolveInitialLocale(): Locale {
  const storedLocale = getStoredLocale();
  if (storedLocale) {
    return storedLocale;
  }

  if (typeof window === 'undefined') {
    return DEFAULT_LOCALE;
  }

  const telegramLocale = resolveLocaleFromString(window.Telegram?.WebApp?.initDataUnsafe?.user?.language_code);
  if (telegramLocale) {
    return telegramLocale;
  }

  const navigatorLocale = resolveLocaleFromString(window.navigator.language);
  if (navigatorLocale) {
    return navigatorLocale;
  }

  return DEFAULT_LOCALE;
}

function persistLocale(locale: Locale) {
  try {
    localStorage.setItem(LOCALE_STORAGE_KEY, locale);
  } catch {
    // Ignore persistence issues in restricted environments.
  }

  window.dispatchEvent(new CustomEvent<Locale>(LOCALE_EVENT, { detail: locale }));
}

function getPluralCategory(locale: Locale, count: number): PluralCategory {
  const value = Number.isFinite(count) ? count : 0;
  return pluralRules[locale].select(Math.abs(value)) as PluralCategory;
}

export function i18nText(key: string): string {
  return `${I18N_TEXT_PREFIX}${key}`;
}

export function resolveI18nText(value: string | null | undefined, t: TranslateFn): string | null {
  if (typeof value !== 'string') {
    return value ?? null;
  }

  if (value.startsWith(I18N_TEXT_PREFIX)) {
    const key = value.slice(I18N_TEXT_PREFIX.length);
    return t(key);
  }

  const requestFailedStatusMatch = /^Request failed with status\s+(\d{3})$/i.exec(value.trim());
  if (requestFailedStatusMatch) {
    return t('errors.requestFailedWithStatus', { status: requestFailedStatusMatch[1] });
  }

  return value;
}

export function I18nProvider({ children }: { children: ReactNode }) {
  const [locale, setLocaleState] = useState<Locale>(() => {
    if (typeof window === 'undefined') {
      return DEFAULT_LOCALE;
    }
    return resolveInitialLocale();
  });

  useEffect(() => {
    if (typeof window === 'undefined') {
      return;
    }
    persistLocale(locale);
  }, [locale]);

  useEffect(() => {
    const handleStorage = (event: StorageEvent) => {
      if (event.key !== LOCALE_STORAGE_KEY) {
        return;
      }

      const nextLocale = isLocale(event.newValue) ? event.newValue : resolveInitialLocale();
      setLocaleState((prevLocale) => (prevLocale === nextLocale ? prevLocale : nextLocale));
    };

    const handleLocaleChange = (event: Event) => {
      const nextLocale = (event as CustomEvent<Locale>).detail;
      if (!isLocale(nextLocale)) {
        return;
      }

      setLocaleState((prevLocale) => (prevLocale === nextLocale ? prevLocale : nextLocale));
    };

    window.addEventListener('storage', handleStorage);
    window.addEventListener(LOCALE_EVENT, handleLocaleChange as EventListener);

    return () => {
      window.removeEventListener('storage', handleStorage);
      window.removeEventListener(LOCALE_EVENT, handleLocaleChange as EventListener);
    };
  }, []);

  const setLocale = useCallback((nextLocale: Locale) => {
    setLocaleState(nextLocale);
  }, []);

  const t = useCallback<TranslateFn>(
    (key: string, params?: TranslateParams) => {
      return translate(locale, key, params);
    },
    [locale],
  );

  const tp = useCallback(
    (baseKey: string, count: number, params?: TranslateParams) => {
      const category = getPluralCategory(locale, count);
      const keyedValue = translate(locale, `${baseKey}.${category}`, params);
      if (keyedValue !== `${baseKey}.${category}`) {
        return keyedValue;
      }

      const fallbackOther = translate(locale, `${baseKey}.other`, params);
      if (fallbackOther !== `${baseKey}.other`) {
        return fallbackOther;
      }

      return translate(DEFAULT_LOCALE, `${baseKey}.other`, params);
    },
    [locale],
  );

  const contextValue = useMemo<I18nContextValue>(
    () => ({
      locale,
      localeCode: localeCodeMap[locale],
      setLocale,
      t,
      tp,
    }),
    [locale, setLocale, t, tp],
  );

  return <I18nContext.Provider value={contextValue}>{children}</I18nContext.Provider>;
}

export function useI18n(): I18nContextValue {
  const context = useContext(I18nContext);
  if (!context) {
    throw new Error('useI18n must be used inside I18nProvider');
  }
  return context;
}
