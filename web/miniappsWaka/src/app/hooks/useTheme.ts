import { useEffect, useState } from 'react';
import type { TelegramWebApp } from '../types/telegram';

const THEME_KEY = 'waka_theme';
const THEME_EVENT = 'waka_theme_change';

export type Theme = 'dark' | 'light';

function isTheme(value: unknown): value is Theme {
  return value === 'dark' || value === 'light';
}

function getTelegramTheme(webApp?: TelegramWebApp): Theme | null {
  if (!webApp || !isTheme(webApp.colorScheme)) {
    return null;
  }
  return webApp.colorScheme;
}

function getStoredTheme(): Theme | null {
  try {
    const storedTheme = localStorage.getItem(THEME_KEY);
    return isTheme(storedTheme) ? storedTheme : null;
  } catch {
    return null;
  }
}

function getSystemTheme(): Theme {
  return window.matchMedia('(prefers-color-scheme: dark)').matches ? 'dark' : 'light';
}

function resolveInitialTheme(): Theme {
  const storedTheme = getStoredTheme();
  if (storedTheme) {
    return storedTheme;
  }

  const telegramTheme = getTelegramTheme(window.Telegram?.WebApp);
  if (telegramTheme) {
    return telegramTheme;
  }

  return getSystemTheme();
}

function applyTheme(theme: Theme) {
  const root = document.documentElement;
  root.classList.toggle('dark', theme === 'dark');
  root.setAttribute('data-theme', theme);
}

function persistTheme(theme: Theme) {
  try {
    localStorage.setItem(THEME_KEY, theme);
  } catch {
    // Ignore write errors (private mode, quota exceeded, etc.).
  }
  window.dispatchEvent(new CustomEvent<Theme>(THEME_EVENT, { detail: theme }));
}

export function initializeTheme() {
  if (typeof window === 'undefined') {
    return 'dark' as const;
  }
  const initialTheme = resolveInitialTheme();
  applyTheme(initialTheme);
  return initialTheme;
}

export function useTheme() {
  const [theme, setThemeState] = useState<Theme>(() => {
    if (typeof window === 'undefined') {
      return 'dark';
    }
    return resolveInitialTheme();
  });

  useEffect(() => {
    applyTheme(theme);
    persistTheme(theme);
  }, [theme]);

  useEffect(() => {
    const handleStorage = (event: StorageEvent) => {
      if (event.key !== THEME_KEY) {
        return;
      }
      const nextTheme = isTheme(event.newValue) ? event.newValue : resolveInitialTheme();
      setThemeState(nextTheme);
    };

    const handleThemeChange = (event: Event) => {
      const nextTheme = (event as CustomEvent<Theme>).detail;
      if (!isTheme(nextTheme)) {
        return;
      }
      setThemeState((previousTheme) => (previousTheme === nextTheme ? previousTheme : nextTheme));
    };

    window.addEventListener('storage', handleStorage);
    window.addEventListener(THEME_EVENT, handleThemeChange as EventListener);

    return () => {
      window.removeEventListener('storage', handleStorage);
      window.removeEventListener(THEME_EVENT, handleThemeChange as EventListener);
    };
  }, []);

  const setTheme = (nextTheme: Theme) => {
    setThemeState(nextTheme);
  };

  const toggleTheme = () => {
    setThemeState((previousTheme) => (previousTheme === 'dark' ? 'light' : 'dark'));
  };

  return { theme, setTheme, toggleTheme };
}
