import { useMemo } from 'react';
import { useAuth } from './useAuth';
import { i18nText } from '../../shared/i18n';

function withImageVersion(url: string, version: string): string {
  if (!url || !version) {
    return url;
  }

  const [base, hashPart] = url.split('#');
  const [path, query = ''] = base.split('?');
  const params = new URLSearchParams(query);
  params.set('v', version);
  const queryString = params.toString();

  return `${path}${queryString ? `?${queryString}` : ''}${hashPart ? `#${hashPart}` : ''}`;
}

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
      const photoUrl = isBackendProfile
        ? user?.photo_url
          ? withImageVersion(user.photo_url, user.updated_at)
          : null
        : sourceUser.photo_url ?? null;

      return {
        displayName,
        handle,
        photoUrl,
        source: isBackendProfile ? 'backend' : 'telegram',
      };
    }

    if (!sourceUser) {
      return {
        displayName: i18nText('common.guest'),
        handle: hasTelegramContext
          ? i18nText('common.telegramProfileUnavailable')
          : i18nText('common.openFromTelegramToSyncProfile'),
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
