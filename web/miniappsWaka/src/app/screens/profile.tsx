import { useState } from 'react';
import { Link } from 'react-router';
import { Heart, HelpCircle, ChevronRight, Bell, Shield, FileText, Sun, Moon, Globe } from 'lucide-react';
import type { LucideIcon } from 'lucide-react';
import { useFavorites } from '../hooks/useFavorites';
import { useTheme } from '../hooks/useTheme';
import { useTelegramUser } from '../hooks/useTelegramUser';
import { WakaFullLogo, WakaIcon } from '../components/waka-brand';
import { ImageWithFallback } from '../components/image-with-fallback';
import { LanguagePicker, LANGUAGE_OPTIONS } from '../components/language-picker';
import { useAuth } from '../hooks/useAuth';
import { resolveI18nText, useI18n } from '../../shared/i18n';

type SettingsItem = {
  icon: LucideIcon;
  labelKey: string;
  helperTextKey?: string;
  path?: string;
};

type SettingsSection = {
  titleKey: string;
  items: SettingsItem[];
};

const settingsSections: SettingsSection[] = [
  {
    titleKey: 'common.account',
    items: [{ icon: Bell, labelKey: 'common.notifications', helperTextKey: 'common.notificationsManagedTelegram' }],
  },
  {
    titleKey: 'common.support',
    items: [{ icon: HelpCircle, labelKey: 'common.helpCenter', path: '/faq' }],
  },
  {
    titleKey: 'common.legal',
    items: [
      { icon: FileText, labelKey: 'common.termsOfService', helperTextKey: 'common.contentComingSoon', path: '/legal/terms' },
      { icon: Shield, labelKey: 'common.privacyPolicy', helperTextKey: 'common.contentComingSoon', path: '/legal/privacy' },
    ],
  },
];

export function Profile() {
  const { locale, localeCode, setLocale, t } = useI18n();
  const [langPickerOpen, setLangPickerOpen] = useState(false);
  const { favorites } = useFavorites();
  const { theme, toggleTheme } = useTheme();
  const { profile, telegramUser } = useTelegramUser();
  const { user, error, debug } = useAuth();
  const localizedError = resolveI18nText(error, t);
  const firstName = user?.first_name ?? telegramUser?.first_name ?? '';
  const lastName = user?.last_name ?? telegramUser?.last_name ?? '';
  const hasNameLines = Boolean(firstName || lastName);
  const tgID = user?.tg_id ?? telegramUser?.id ?? null;
  const currentLang = LANGUAGE_OPTIONS.find((item) => item.code === locale) ?? LANGUAGE_OPTIONS[0];
  const localizedDisplayName = resolveI18nText(profile.displayName, t) ?? profile.displayName;
  const localizedHandle = resolveI18nText(profile.handle, t) ?? profile.handle;

  return (
    <div className="min-h-screen pb-32">
      <div className="px-6 pb-8 pt-14">
        <h1 className="text-4xl font-extrabold leading-none tracking-tighter">{t('profile.title')}</h1>
      </div>

      <div className="mb-8 px-6">
        <div className="relative overflow-hidden rounded-[36px] border border-border/50 bg-card p-8 shadow-sm dark:shadow-none">
          <div className="pointer-events-none absolute right-0 top-0 h-48 w-48 -translate-y-1/2 translate-x-1/4 rounded-full bg-foreground/5 blur-3xl" />
          <div className="relative z-10 flex items-center gap-6">
            <div className="flex h-20 w-20 items-center justify-center overflow-hidden rounded-[24px] border border-border/50 bg-background shadow-sm">
              {profile.photoUrl ? (
                <ImageWithFallback src={profile.photoUrl} alt={profile.displayName} className="h-full w-full object-cover" />
              ) : (
                <WakaIcon size={80} className="rounded-[22px]" />
              )}
            </div>
            <div className="min-w-0 flex-1">
              <h2 className="mb-1 text-[28px] font-bold tracking-tight leading-[1.05] text-foreground">
                {hasNameLines ? (
                  <>
                    {firstName && <span className="block break-words">{firstName}</span>}
                    {lastName && <span className="block break-words">{lastName}</span>}
                  </>
                ) : (
                  <span className="block break-words">{localizedDisplayName}</span>
                )}
              </h2>
              <p className="truncate text-[11px] font-bold uppercase tracking-[0.1em] text-foreground/70">{localizedHandle}</p>
              {tgID && (
                <p className="mt-2 truncate text-[10px] font-semibold uppercase tracking-[0.12em] text-muted-foreground">
                  {t('common.tgId')}: {tgID}
                </p>
              )}
            </div>
          </div>
        </div>
      </div>

      {localizedError && (
        <div className="px-6 mb-8">
          <div className="rounded-2xl border border-destructive/30 bg-destructive/10 p-4 text-sm text-destructive">
            {localizedError}
          </div>
        </div>
      )}

      {import.meta.env.DEV && debug && (
        <div className="px-6 mb-8">
          <details className="rounded-2xl border border-border/50 bg-card p-4">
            <summary className="text-[11px] font-bold uppercase tracking-[0.12em] text-muted-foreground cursor-pointer">
              {t('common.telegramAuthDebug')}
            </summary>
            <pre className="mt-3 overflow-auto text-[11px] leading-relaxed text-muted-foreground">
              {JSON.stringify(debug, null, 2)}
            </pre>
          </details>
        </div>
      )}

      <div className="mb-10 px-6">
        <div className="grid grid-cols-2 gap-4">
          <Link
            to="/favorites"
            className="group relative overflow-hidden rounded-[32px] border border-border/50 bg-card p-6 shadow-sm transition-all duration-500 hover:shadow-lg dark:shadow-none"
          >
            <div className="absolute right-0 top-0 h-24 w-24 -translate-y-1/2 translate-x-1/2 rounded-full bg-gradient-to-bl from-foreground/5 to-transparent opacity-0 blur-xl transition-opacity duration-500 group-hover:opacity-100" />
            <div className="mb-4 flex h-12 w-12 items-center justify-center rounded-[18px] border border-border/50 bg-background shadow-sm transition-transform duration-500 group-hover:scale-110 group-hover:border-foreground/30">
              <Heart className="h-5 w-5 text-foreground" />
            </div>
            <p className="mb-1 text-3xl font-extrabold tracking-tighter text-foreground">{favorites.length.toLocaleString(localeCode)}</p>
            <p className="text-[10px] font-bold uppercase tracking-[0.15em] text-muted-foreground">{t('common.savedModels')}</p>
          </Link>
          <Link
            to="/faq"
            className="group relative overflow-hidden rounded-[32px] border border-border/50 bg-card p-6 shadow-sm transition-all duration-500 hover:shadow-lg dark:shadow-none"
          >
            <div className="absolute right-0 top-0 h-24 w-24 -translate-y-1/2 translate-x-1/2 rounded-full bg-gradient-to-bl from-foreground/5 to-transparent opacity-0 blur-xl transition-opacity duration-500 group-hover:opacity-100" />
            <div className="mb-4 flex h-12 w-12 items-center justify-center rounded-[18px] border border-border/50 bg-background shadow-sm transition-transform duration-500 group-hover:scale-110 group-hover:border-foreground/30">
              <HelpCircle className="h-5 w-5 text-foreground" />
            </div>
            <p className="mb-1 text-3xl font-extrabold tracking-tighter text-foreground">24/7</p>
            <p className="text-[10px] font-bold uppercase tracking-[0.15em] text-muted-foreground">{t('profile.supportTile')}</p>
          </Link>
        </div>
      </div>

      <div className="px-6">
        <div className="mb-10">
          <h3 className="mb-4 px-2 text-[10px] font-bold uppercase tracking-[0.2em] text-muted-foreground">{t('common.settings')}</h3>
          <div className="overflow-hidden rounded-[32px] border border-border/50 bg-card shadow-sm">
            <button
              type="button"
              onClick={toggleTheme}
              className="group flex w-full min-w-0 items-center gap-4 px-5 py-5 text-left transition-all hover:bg-foreground/5 border-b border-border/50"
            >
              <div className="flex min-w-0 flex-1 items-center gap-4">
                <div className="flex h-10 w-10 shrink-0 items-center justify-center rounded-[14px] border border-border/50 bg-background shadow-sm transition-transform duration-500 group-hover:scale-105 group-hover:border-foreground/30">
                  {theme === 'dark' ? <Moon className="h-4 w-4 text-foreground" /> : <Sun className="h-4 w-4 text-foreground" />}
                </div>
                <span className="truncate font-bold tracking-tight text-foreground">{t('common.appearance')}</span>
              </div>
              <div className="ml-2 flex shrink-0 items-center gap-2">
                <span className="max-w-[100px] truncate text-right text-[10px] font-bold uppercase tracking-[0.06em] text-foreground">
                  {theme === 'dark' ? t('common.dark') : t('common.light')}
                </span>
                <div
                  className={`h-8 w-14 rounded-full p-1 transition-colors duration-500 ${
                    theme === 'dark' ? 'bg-foreground' : 'bg-border/80'
                  }`}
                >
                  <div
                    className={`h-6 w-6 rounded-full bg-background shadow-md transition-transform duration-500 ease-[cubic-bezier(0.2,0.8,0.2,1)] ${
                      theme === 'dark' ? 'translate-x-6' : 'translate-x-0'
                    }`}
                  />
                </div>
              </div>
            </button>

            <button
              type="button"
              onClick={() => setLangPickerOpen(true)}
              className="group flex w-full min-w-0 items-center gap-4 px-5 py-5 text-left transition-all hover:bg-foreground/5"
            >
              <div className="flex min-w-0 flex-1 items-center gap-4">
                <div className="flex h-10 w-10 shrink-0 items-center justify-center rounded-[14px] border border-border/50 bg-background shadow-sm transition-transform duration-500 group-hover:scale-105 group-hover:border-foreground/30">
                  <Globe className="h-4 w-4 text-foreground" />
                </div>
                <span className="truncate font-bold tracking-tight text-foreground">{t('common.language')}</span>
              </div>
              <div className="ml-2 flex shrink-0 items-center gap-2">
                <span className="text-xl leading-none">{currentLang.flag}</span>
                <span className="max-w-[88px] truncate text-[11px] font-bold uppercase tracking-[0.06em] text-muted-foreground">
                  {currentLang.labelEn}
                </span>
                <ChevronRight className="h-4 w-4 text-muted-foreground/50 transition-colors group-hover:text-foreground" />
              </div>
            </button>
          </div>
        </div>

        {settingsSections.map((section, sectionIndex) => (
          <div key={section.titleKey} className={sectionIndex > 0 ? 'mt-10' : ''}>
            <h3 className="mb-4 px-2 text-[10px] font-bold uppercase tracking-[0.2em] text-muted-foreground">{t(section.titleKey)}</h3>
            <div className="overflow-hidden rounded-[32px] border border-border/50 bg-card shadow-sm">
              {section.items.map((item, itemIndex) => {
                const Icon = item.icon;
                const isLast = itemIndex === section.items.length - 1;
                const rowClassName = `flex items-center gap-5 px-6 py-5 text-left ${!isLast ? 'border-b border-border/50' : ''}`;
                const label = t(item.labelKey);
                const helperText = item.helperTextKey ? t(item.helperTextKey) : null;

                const content = (
                  <>
                    <div className="flex h-10 w-10 items-center justify-center rounded-[14px] border border-border/50 bg-background shadow-sm">
                      <Icon className="h-4 w-4 text-foreground" />
                    </div>
                    <div className="flex min-w-0 flex-1 flex-col">
                      <span className="font-bold tracking-tight text-foreground">{label}</span>
                      {helperText && (
                        <span className="truncate text-[10px] font-semibold uppercase tracking-[0.12em] text-muted-foreground">
                          {helperText}
                        </span>
                      )}
                    </div>
                  </>
                );

                if (item.path) {
                  return (
                    <Link key={`${section.titleKey}:${item.path}`} to={item.path} className={`${rowClassName} group transition-all hover:bg-foreground/5`}>
                      {content}
                      <ChevronRight className="h-4 w-4 text-muted-foreground/50 transition-colors group-hover:text-foreground" />
                    </Link>
                  );
                }

                return (
                  <div key={`${section.titleKey}:${item.labelKey}`} className={rowClassName}>
                    {content}
                  </div>
                );
              })}
            </div>
          </div>
        ))}

        <div className="mb-8 mt-12 text-center">
          <WakaFullLogo height={50} className="mx-auto mb-4 opacity-20" />
          <p className="text-[10px] font-bold uppercase tracking-[0.2em] text-muted-foreground">{t('profile.wakaOs')}</p>
          <p className="mt-1.5 text-[9px] font-bold tracking-[0.1em] text-muted-foreground/50">{t('common.versionLabel', { version: '1.0.0' })}</p>
        </div>
      </div>

      <LanguagePicker
        open={langPickerOpen}
        current={locale}
        onSelect={setLocale}
        onClose={() => setLangPickerOpen(false)}
      />
    </div>
  );
}
