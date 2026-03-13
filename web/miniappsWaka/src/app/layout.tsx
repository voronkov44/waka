import { Outlet } from 'react-router';
import { useEffect } from 'react';
import { BottomNav } from './components/bottom-nav';
import { useTheme } from './hooks/useTheme';
import { AuthProvider } from './context/auth-context';
import { FavoritesProvider } from './context/favorites-context';

export function Layout() {
  const { theme } = useTheme();

  useEffect(() => {
    const webApp = window.Telegram?.WebApp;
    if (!webApp) {
      return;
    }
    webApp.ready?.();
    webApp.expand?.();
  }, []);

  return (
    <div className="min-h-screen bg-background text-foreground transition-colors duration-300" data-theme={theme}>
      <AuthProvider>
        <FavoritesProvider>
          <Outlet />
        </FavoritesProvider>
      </AuthProvider>
      <BottomNav />
    </div>
  );
}
