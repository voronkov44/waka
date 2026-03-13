import { useMemo } from 'react';
import { useAuth } from './useAuth';

export function useTelegramUser() {
  const { user, telegramUser, hasTelegramContext, isAuthenticated } = useAuth();

  const profile = useMemo(() => {
    const sourceUser = user
      ? {
          id: user.tg_id,
          first_name: user.first_name ?? '',
          last_name: user.last_name,
          username: user.username,
          photo_url: user.photo_url,
        }
      : telegramUser;

    const isBackendProfile = Boolean(user);

    if (sourceUser) {
      const fullName = `${sourceUser.first_name ?? ''} ${sourceUser.last_name ?? ''}`.trim();
      const fallbackName = sourceUser.username ? `@${sourceUser.username}` : `User ${sourceUser.id}`;
      const displayName = fullName || fallbackName;
      const handle = sourceUser.username ? `@${sourceUser.username}` : `ID ${sourceUser.id}`;

      return {
        displayName,
        handle,
        photoUrl: sourceUser.photo_url ?? null,
        source: isBackendProfile ? 'backend' : 'telegram',
      };
    }

    if (!sourceUser) {
      return {
        displayName: 'Guest',
        handle: hasTelegramContext
          ? 'Telegram profile is unavailable for this session'
          : 'Open from Telegram to sync profile',
        photoUrl: null,
        source: 'fallback',
      };
    }
  }, [user, telegramUser, hasTelegramContext]);

  return {
    user,
    telegramUser,
    profile,
    hasTelegramContext,
    isAuthenticated,
  };
}
