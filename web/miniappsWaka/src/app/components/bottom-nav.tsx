import { Home, Heart, HelpCircle, User } from 'lucide-react';
import { Link, useLocation } from 'react-router';
import { VapeDeviceIcon } from './icons/vape-device-icon';

const navItems = [
  { path: '/', icon: Home, label: 'Home' },
  { path: '/catalog', icon: VapeDeviceIcon, label: 'Catalog' },
  { path: '/favorites', icon: Heart, label: 'Favorites' },
  { path: '/faq', icon: HelpCircle, label: 'FAQ' },
  { path: '/profile', icon: User, label: 'Profile' },
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

  return (
    <div className="fixed bottom-0 left-0 right-0 z-50 pb-safe pointer-events-none">
      <div className="absolute inset-0 bg-gradient-to-t from-background via-background/90 to-transparent pointer-events-none -z-10 h-32 bottom-0 top-auto" />
      <nav className="mx-4 mb-5 mt-2 bg-card/60 backdrop-blur-[32px] border border-border/50 rounded-full shadow-lg dark:shadow-[0_8px_32px_rgba(0,0,0,0.5)] pointer-events-auto relative overflow-hidden">
        {/* Soft edge highlight */}
        <div className="absolute inset-0 rounded-full ring-1 ring-inset ring-foreground/5 dark:ring-white/5 pointer-events-none mix-blend-overlay" />
        
        <div className="flex items-center justify-around px-2 py-3 max-w-lg mx-auto relative z-10">
          {navItems.map((item) => {
            const isActive = isNavItemActive(location.pathname, item.path);
            const Icon = item.icon;
            
            return (
              <Link
                key={item.path}
                to={item.path}
                className={`relative flex flex-col items-center gap-1.5 px-4 py-2.5 rounded-full transition-all duration-500 ease-out group ${
                  isActive
                    ? 'text-foreground'
                    : 'text-muted-foreground hover:text-foreground'
                }`}
              >
                {isActive && (
                  <div className="absolute inset-0 bg-foreground/10 rounded-full -z-10 animate-in fade-in zoom-in duration-300" />
                )}
                <Icon className={`w-5 h-5 transition-transform duration-500 ${isActive ? 'scale-110 drop-shadow-md' : 'group-hover:scale-110'}`} />
                <span className={`text-[9px] font-bold tracking-[0.1em] uppercase transition-all duration-500 ${isActive ? 'opacity-100' : 'opacity-70'}`}>
                  {item.label}
                </span>
              </Link>
            );
          })}
        </div>
      </nav>
    </div>
  );
}
