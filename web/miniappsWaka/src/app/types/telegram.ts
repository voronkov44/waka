export interface TelegramUser {
  id: number;
  first_name: string;
  last_name?: string;
  username?: string;
  photo_url?: string;
  language_code?: string;
  is_premium?: boolean;
}

export interface TelegramWebApp {
  colorScheme?: 'light' | 'dark';
  initData?: string;
  viewportHeight?: number;
  viewportStableHeight?: number;
  initDataUnsafe?: {
    user?: TelegramUser;
  };
  ready?: () => void;
  expand?: () => void;
  onEvent?: (eventType: string, eventHandler: (...args: unknown[]) => void) => void;
  offEvent?: (eventType: string, eventHandler: (...args: unknown[]) => void) => void;
}

declare global {
  interface Window {
    Telegram?: {
      WebApp?: TelegramWebApp;
    };
    __wakaTelegramBootstrapDebug?: unknown;
  }
}
