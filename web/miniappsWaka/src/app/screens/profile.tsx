import { Link } from 'react-router';
import { Heart, HelpCircle, ChevronRight, Bell, Shield, FileText, Sun, Moon, MessageCircle } from 'lucide-react';
import type { LucideIcon } from 'lucide-react';
import { useFavorites } from '../hooks/useFavorites';
import { useTheme } from '../hooks/useTheme';
import { useTelegramUser } from '../hooks/useTelegramUser';
import { WakaFullLogo, WakaIcon } from '../components/waka-brand';
import { ImageWithFallback } from '../components/image-with-fallback';
import { useAuth } from '../hooks/useAuth';

type SettingsItem = {
  icon: LucideIcon;
  label: string;
  helperText?: string;
  path?: string;
};

type SettingsSection = {
  title: string;
  items: SettingsItem[];
};

const settingsSections: SettingsSection[] = [
  {
    title: 'Account',
    items: [{ icon: Bell, label: 'Notifications', helperText: 'Managed in Telegram' }],
  },
  {
    title: 'Support',
    items: [{ icon: HelpCircle, label: 'Help Center', path: '/faq' }],
  },
  {
    title: 'Legal',
    items: [
      { icon: FileText, label: 'Terms of Service', helperText: 'Provided in the official bot' },
      { icon: Shield, label: 'Privacy Policy', helperText: 'Provided in the official bot' },
    ],
  },
];

export function Profile() {
  const { favorites } = useFavorites();
  const { theme, toggleTheme } = useTheme();
  const { profile, hasTelegramContext } = useTelegramUser();
  const { user, isLoading, isAuthenticated, error, debug } = useAuth();

  return (
    <div className="min-h-screen pb-32">
      <div className="px-6 pb-8 pt-14">
        <h1 className="text-4xl font-extrabold leading-none tracking-tighter">Profile</h1>
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
              <h2 className="mb-1 truncate text-2xl font-bold tracking-tight text-foreground">{profile.displayName}</h2>
              <p className="truncate text-[11px] font-bold uppercase tracking-[0.1em] text-foreground/70">{profile.handle}</p>
              {user && (
                <p className="truncate text-[10px] font-semibold uppercase tracking-[0.12em] text-muted-foreground mt-2">
                  Backend user ID: {user.id} · TG ID: {user.tg_id}
                </p>
              )}
            </div>
          </div>
        </div>
      </div>

      {error && (
        <div className="px-6 mb-8">
          <div className="rounded-2xl border border-destructive/30 bg-destructive/10 p-4 text-sm text-destructive">
            {error}
          </div>
        </div>
      )}

      {import.meta.env.DEV && debug && (
        <div className="px-6 mb-8">
          <details className="rounded-2xl border border-border/50 bg-card p-4">
            <summary className="text-[11px] font-bold uppercase tracking-[0.12em] text-muted-foreground cursor-pointer">
              Telegram auth debug
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
            <p className="mb-1 text-3xl font-extrabold tracking-tighter text-foreground">{favorites.length}</p>
            <p className="text-[10px] font-bold uppercase tracking-[0.15em] text-muted-foreground">Saved Models</p>
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
            <p className="text-[10px] font-bold uppercase tracking-[0.15em] text-muted-foreground">Support</p>
          </Link>
        </div>
      </div>

      <div className="px-6">
        <div className="mb-10">
          <h3 className="mb-4 px-2 text-[10px] font-bold uppercase tracking-[0.2em] text-muted-foreground">Settings</h3>
          <div className="overflow-hidden rounded-[32px] border border-border/50 bg-card shadow-sm">
            <button
              type="button"
              onClick={toggleTheme}
              className="group flex w-full items-center gap-5 px-6 py-5 text-left transition-all hover:bg-foreground/5"
            >
              <div className="flex h-10 w-10 items-center justify-center rounded-[14px] border border-border/50 bg-background shadow-sm transition-transform duration-500 group-hover:scale-105 group-hover:border-foreground/30">
                {theme === 'dark' ? <Moon className="h-4 w-4 text-foreground" /> : <Sun className="h-4 w-4 text-foreground" />}
              </div>
              <span className="flex-1 font-bold tracking-tight text-foreground">Appearance</span>
              <span className="mr-3 text-[10px] font-bold uppercase tracking-[0.1em] text-foreground">
                {theme === 'dark' ? 'Dark' : 'Light'}
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
            </button>
          </div>
        </div>

        {settingsSections.map((section, sectionIndex) => (
          <div key={section.title} className={sectionIndex > 0 ? 'mt-10' : ''}>
            <h3 className="mb-4 px-2 text-[10px] font-bold uppercase tracking-[0.2em] text-muted-foreground">{section.title}</h3>
            <div className="overflow-hidden rounded-[32px] border border-border/50 bg-card shadow-sm">
              {section.items.map((item, itemIndex) => {
                const Icon = item.icon;
                const isLast = itemIndex === section.items.length - 1;
                const rowClassName = `flex items-center gap-5 px-6 py-5 text-left ${!isLast ? 'border-b border-border/50' : ''}`;

                const content = (
                  <>
                    <div className="flex h-10 w-10 items-center justify-center rounded-[14px] border border-border/50 bg-background shadow-sm">
                      <Icon className="h-4 w-4 text-foreground" />
                    </div>
                    <div className="flex min-w-0 flex-1 flex-col">
                      <span className="font-bold tracking-tight text-foreground">{item.label}</span>
                      {item.helperText && (
                        <span className="truncate text-[10px] font-semibold uppercase tracking-[0.12em] text-muted-foreground">
                          {item.helperText}
                        </span>
                      )}
                    </div>
                  </>
                );

                if (item.path) {
                  return (
                    <Link key={item.label} to={item.path} className={`${rowClassName} group transition-all hover:bg-foreground/5`}>
                      {content}
                      <ChevronRight className="h-4 w-4 text-muted-foreground/50 transition-colors group-hover:text-foreground" />
                    </Link>
                  );
                }

                return (
                  <div key={item.label} className={rowClassName}>
                    {content}
                  </div>
                );
              })}
            </div>
          </div>
        ))}

        <div className="mt-10 rounded-[32px] border border-border/50 bg-card p-6 shadow-sm">
          <div className="flex items-start gap-4">
            <div className="flex h-10 w-10 items-center justify-center rounded-[14px] border border-border/50 bg-background shadow-sm">
              <MessageCircle className="h-4 w-4 text-foreground" />
            </div>
            <div>
              <p className="mb-1 font-bold tracking-tight text-foreground">Telegram account is your profile source</p>
              <p className="text-sm font-medium leading-relaxed text-muted-foreground">
                {hasTelegramContext
                  ? isLoading
                    ? 'Syncing your Telegram identity with backend profile...'
                    : isAuthenticated
                      ? 'Your profile is synced with the backend and Telegram.'
                      : 'Trying to authorize your Telegram identity.'
                  : 'Open this mini app from Telegram to load your account details automatically.'}
              </p>
            </div>
          </div>
        </div>

        <div className="mb-8 mt-12 text-center">
          <WakaFullLogo height={50} className="mx-auto mb-4 opacity-20" />
          <p className="text-[10px] font-bold uppercase tracking-[0.2em] text-muted-foreground">Waka OS</p>
          <p className="mt-1.5 text-[9px] font-bold tracking-[0.1em] text-muted-foreground/50">Version 1.0.0</p>
        </div>
      </div>
    </div>
  );
}
