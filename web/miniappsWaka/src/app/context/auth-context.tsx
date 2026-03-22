import { createContext, useCallback, useContext, useEffect, useMemo, useState, type ReactNode } from 'react';
import { apiClient } from '../api/client';
import { ApiError } from '../api/http';
import { clearAuthToken, readAuthToken, writeAuthToken } from '../api/auth-storage';
import { bootstrapTelegramContext, type TelegramBootstrapDiagnostics } from '../telegram/bootstrap';
import type { MeResponseDTO, TelegramAuthRequestDTO } from '../api/types';
import type { TelegramUser } from '../types/telegram';
import { i18nText } from '../../shared/i18n';

export interface AuthDebugState extends TelegramBootstrapDiagnostics {
  usedStoredToken: boolean;
  authRequestSent: boolean;
  authSucceeded: boolean;
  tokenStored: boolean;
  fallbackReason: string | null;
}

interface AuthContextValue {
  token: string | null;
  user: MeResponseDTO | null;
  telegramUser: TelegramUser | null;
  hasTelegramContext: boolean;
  isLoading: boolean;
  isAuthenticated: boolean;
  error: string | null;
  debug: AuthDebugState | null;
  refreshUser: () => Promise<void>;
}

const AuthContext = createContext<AuthContextValue | null>(null);

function emptyBootstrapDiagnostics(): TelegramBootstrapDiagnostics {
  return {
    webAppDetected: false,
    initDataPresent: false,
    initDataSource: 'none',
    userPresent: false,
    userSource: 'none',
    attempts: 0,
    readyCalled: false,
    expandCalled: false,
  };
}

function mapTelegramUserToAuthPayload(user: TelegramUser): TelegramAuthRequestDTO {
  return {
    tg_id: user.id,
    username: user.username,
    first_name: user.first_name,
    last_name: user.last_name,
    photo_url: user.photo_url,
  };
}

async function loadCurrentUserOrNull(): Promise<MeResponseDTO | null> {
  try {
    return await apiClient.getCurrentUser();
  } catch (error) {
    if (error instanceof ApiError && error.status === 401) {
      return null;
    }
    throw error;
  }
}

function publishDebug(debug: AuthDebugState) {
  window.__wakaTelegramBootstrapDebug = debug;
  if (import.meta.env.DEV) {
    console.debug('[waka:tma] auth diagnostics', debug);
  }
}

export function AuthProvider({ children }: { children: ReactNode }) {
  const [token, setToken] = useState<string | null>(null);
  const [user, setUser] = useState<MeResponseDTO | null>(null);
  const [telegramUser, setTelegramUser] = useState<TelegramUser | null>(null);
  const [hasTelegramContext, setHasTelegramContext] = useState(false);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [debug, setDebug] = useState<AuthDebugState | null>(null);

  const refreshUser = useCallback(async () => {
    if (!readAuthToken()) {
      setUser(null);
      return;
    }

    const nextUser = await loadCurrentUserOrNull();
    if (!nextUser) {
      clearAuthToken();
      setToken(null);
      setUser(null);
      return;
    }

    setUser(nextUser);
  }, []);

  useEffect(() => {
    let cancelled = false;

    const bootstrap = async () => {
      setIsLoading(true);
      setError(null);

      let telegram: Awaited<ReturnType<typeof bootstrapTelegramContext>> | null = null;
      try {
        telegram = await bootstrapTelegramContext();
      } catch (err) {
        const diagnostics: AuthDebugState = {
          ...emptyBootstrapDiagnostics(),
          usedStoredToken: false,
          authRequestSent: false,
          authSucceeded: false,
          tokenStored: false,
          fallbackReason: 'telegram_bootstrap_failed',
        };
        if (!cancelled) {
          setHasTelegramContext(false);
          setTelegramUser(null);
          setDebug(diagnostics);
          publishDebug(diagnostics);
          setError(err instanceof Error ? err.message : i18nText('errors.telegramBootstrapFailed'));
          setIsLoading(false);
        }
        return;
      }

      if (cancelled) {
        return;
      }

      if (!telegram) {
        setIsLoading(false);
        return;
      }

      setHasTelegramContext(telegram.hasTelegramContext);
      setTelegramUser(telegram.user);

      const diagnostics: AuthDebugState = {
        ...telegram.diagnostics,
        usedStoredToken: false,
        authRequestSent: false,
        authSucceeded: false,
        tokenStored: false,
        fallbackReason: null,
      };

      let pendingError: string | null = null;

      const existingToken = readAuthToken();
      if (existingToken) {
        diagnostics.usedStoredToken = true;
        setToken(existingToken);
        if (import.meta.env.DEV) {
          console.debug('[waka:tma] validating stored token via /api/auth/me');
        }

        try {
          const me = await loadCurrentUserOrNull();
          if (cancelled) {
            return;
          }

          if (me) {
            setUser(me);
            diagnostics.authSucceeded = true;
          } else {
            clearAuthToken();
            setToken(null);
            setUser(null);
            diagnostics.fallbackReason = 'stored_token_invalid_or_expired';
          }
        } catch (err) {
          clearAuthToken();
          setToken(null);
          setUser(null);
          diagnostics.fallbackReason = 'stored_token_me_request_failed';
          pendingError = err instanceof Error ? err.message : i18nText('errors.storedTokenValidationFailed');
        }
      }

      const hadStoredSession = diagnostics.authSucceeded;

      if (telegram.user) {
        diagnostics.authRequestSent = true;
        if (import.meta.env.DEV) {
          console.debug('[waka:tma] sending Telegram auth request to /api/auth/telegram');
        }
        try {
          const authResponse = await apiClient.loginTelegram(mapTelegramUserToAuthPayload(telegram.user));
          if (cancelled) {
            return;
          }

          writeAuthToken(authResponse.token);
          diagnostics.tokenStored = true;
          setToken(authResponse.token);
          if (import.meta.env.DEV) {
            console.debug('[waka:tma] token stored, fetching /api/auth/me');
          }

          const me = await loadCurrentUserOrNull();
          if (cancelled) {
            return;
          }

          if (me) {
            setUser(me);
            diagnostics.authSucceeded = true;
            diagnostics.fallbackReason = null;
            pendingError = null;
          } else if (hadStoredSession && existingToken) {
            writeAuthToken(existingToken);
            setToken(existingToken);
            diagnostics.authSucceeded = true;
            diagnostics.fallbackReason = 'telegram_refresh_me_failed_used_stored_session';
            pendingError = null;
          } else {
            clearAuthToken();
            setToken(null);
            setUser(null);
            diagnostics.fallbackReason = 'telegram_auth_succeeded_but_me_missing';
            pendingError = pendingError ?? i18nText('errors.telegramLoginSucceededProfileMissing');
          }
        } catch (err) {
          if (hadStoredSession && existingToken) {
            writeAuthToken(existingToken);
            setToken(existingToken);
            diagnostics.authSucceeded = true;
            diagnostics.fallbackReason = 'telegram_refresh_request_failed_used_stored_session';
            pendingError = null;
          } else {
            clearAuthToken();
            setToken(null);
            setUser(null);
            diagnostics.fallbackReason = 'telegram_auth_request_failed';
            pendingError = err instanceof Error ? err.message : i18nText('errors.telegramAuthenticationFailed');
          }
        }
      }

      if (!diagnostics.authSucceeded && !telegram.user && !diagnostics.fallbackReason) {
        diagnostics.fallbackReason = 'telegram_user_not_available';
      }

      if (!cancelled) {
        setDebug(diagnostics);
        publishDebug(diagnostics);
        setError(diagnostics.authSucceeded ? null : pendingError);
        setIsLoading(false);
      }
    };

    void bootstrap();

    return () => {
      cancelled = true;
    };
  }, []);

  const contextValue = useMemo<AuthContextValue>(
    () => ({
      token,
      user,
      telegramUser,
      hasTelegramContext,
      isLoading,
      isAuthenticated: Boolean(token && user),
      error,
      debug,
      refreshUser,
    }),
    [token, user, telegramUser, hasTelegramContext, isLoading, error, debug, refreshUser],
  );

  return <AuthContext.Provider value={contextValue}>{children}</AuthContext.Provider>;
}

export function useAuthContext(): AuthContextValue {
  const context = useContext(AuthContext);
  if (!context) {
    throw new Error('useAuthContext must be used inside AuthProvider');
  }
  return context;
}
