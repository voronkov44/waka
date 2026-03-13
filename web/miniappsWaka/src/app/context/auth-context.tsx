import { createContext, useCallback, useContext, useEffect, useMemo, useState, type ReactNode } from 'react';
import { apiClient } from '../api/client';
import { ApiError } from '../api/http';
import { clearAuthToken, readAuthToken, writeAuthToken } from '../api/auth-storage';
import type { MeResponseDTO, TelegramAuthRequestDTO } from '../api/types';
import type { TelegramUser } from '../types/telegram';

interface AuthContextValue {
  token: string | null;
  user: MeResponseDTO | null;
  telegramUser: TelegramUser | null;
  hasTelegramContext: boolean;
  isLoading: boolean;
  isAuthenticated: boolean;
  error: string | null;
  refreshUser: () => Promise<void>;
}

const AuthContext = createContext<AuthContextValue | null>(null);

function readTelegramUser(): TelegramUser | null {
  return window.Telegram?.WebApp?.initDataUnsafe?.user ?? null;
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

export function AuthProvider({ children }: { children: ReactNode }) {
  const [token, setToken] = useState<string | null>(null);
  const [user, setUser] = useState<MeResponseDTO | null>(null);
  const [telegramUser, setTelegramUser] = useState<TelegramUser | null>(null);
  const [hasTelegramContext, setHasTelegramContext] = useState(false);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

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

      const hasContext = Boolean(window.Telegram?.WebApp);
      const tgUser = readTelegramUser();

      if (!cancelled) {
        setHasTelegramContext(hasContext);
        setTelegramUser(tgUser);
      }

      const existingToken = readAuthToken();
      if (existingToken) {
        setToken(existingToken);
        try {
          const me = await loadCurrentUserOrNull();
          if (cancelled) {
            return;
          }
          if (me) {
            setUser(me);
            setIsLoading(false);
            return;
          }
          clearAuthToken();
          setToken(null);
          setUser(null);
        } catch (err) {
          if (cancelled) {
            return;
          }
          setError(err instanceof Error ? err.message : 'Failed to authorize');
          setIsLoading(false);
          return;
        }
      }

      if (tgUser) {
        try {
          const authResponse = await apiClient.loginTelegram(mapTelegramUserToAuthPayload(tgUser));
          if (cancelled) {
            return;
          }

          writeAuthToken(authResponse.token);
          setToken(authResponse.token);

          const me = await loadCurrentUserOrNull();
          if (cancelled) {
            return;
          }

          if (me) {
            setUser(me);
          } else {
            clearAuthToken();
            setToken(null);
            setUser(null);
          }
        } catch (err) {
          if (cancelled) {
            return;
          }
          clearAuthToken();
          setToken(null);
          setUser(null);
          setError(err instanceof Error ? err.message : 'Failed to authorize');
        }
      }

      if (!cancelled) {
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
      refreshUser,
    }),
    [token, user, telegramUser, hasTelegramContext, isLoading, error, refreshUser],
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
