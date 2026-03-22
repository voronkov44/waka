import { Home, Heart, HelpCircle, User } from 'lucide-react';
import { Link, useLocation } from 'react-router';
import { VapeDeviceIcon } from './icons/vape-device-icon';
import { useI18n } from '../../shared/i18n';

const navItems = [
  { path: '/', icon: Home, labelKey: 'nav.home' },
  { path: '/catalog', icon: VapeDeviceIcon, labelKey: 'nav.catalog' },
  { path: '/favorites', icon: Heart, labelKey: 'nav.favorites' },
  { path: '/faq', icon: HelpCircle, labelKey: 'nav.faq' },
  { path: '/profile', icon: User, labelKey: 'nav.profile' },
];

function isNavItemActive(currentPath: string, itemPath: string) {
  if (itemPath === '/') {
    return currentPath === '/';
  }

  if (itemPath === '/catalog') {
    return currentPath === '/catalog' || currentPath.startsWith('/product/');
  }

  return currentPath === itemPath || currentPath.startsWith(`${itemPath}/`);
}

export function BottomNav() {
  const location = useLocation();
  const { t } = useI18n();

  return (
    <div className="fixed bottom-0 left-0 right-0 z-50 pb-safe pointer-events-none">
      <div className="absolute inset-0 bg-gradient-to-t from-background via-background/90 to-transparent pointer-events-none -z-10 h-32 bottom-0 top-auto" />
      <nav className="mx-3 mb-4 mt-2 bg-card/60 backdrop-blur-[32px] border border-border/50 rounded-full shadow-lg dark:shadow-[0_8px_32px_rgba(0,0,0,0.5)] pointer-events-auto relative overflow-hidden">
        {/* Soft edge highlight */}
        <div className="absolute inset-0 rounded-full ring-1 ring-inset ring-foreground/5 dark:ring-white/5 pointer-events-none mix-blend-overlay" />
        
        <div className="relative z-10 mx-auto flex max-w-lg items-stretch justify-between px-1 py-2.5">
          {navItems.map((item) => {
            const isActive = isNavItemActive(location.pathname, item.path);
            const Icon = item.icon;
            
            return (
              <Link
                key={item.path}
                to={item.path}
                className={`group relative flex min-w-0 flex-1 flex-col items-center gap-1 px-1 py-2 rounded-[18px] transition-all duration-500 ease-out ${
                  isActive
                    ? 'text-foreground'
                    : 'text-muted-foreground hover:text-foreground'
                }`}
              >
                {isActive && (
                  <div className="absolute inset-0 bg-foreground/10 rounded-full -z-10 animate-in fade-in zoom-in duration-300" />
                )}
                <Icon className={`h-[18px] w-[18px] transition-transform duration-500 ${isActive ? 'scale-110 drop-shadow-md' : 'group-hover:scale-110'}`} />
                <span
                  className={`block w-full truncate px-1 text-center text-[8px] font-bold tracking-[0.04em] leading-none transition-all duration-500 ${
                    isActive ? 'opacity-100' : 'opacity-75'
                  }`}
                >
                  {t(item.labelKey)}
                </span>
              </Link>
            );
          })}
        </div>
      </nav>
    </div>
  );
}
