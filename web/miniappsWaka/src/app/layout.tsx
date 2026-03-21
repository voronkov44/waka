import { Outlet } from 'react-router';
import { useEffect } from 'react';
import { BottomNav } from './components/bottom-nav';
import { useTheme } from './hooks/useTheme';
import { AuthProvider } from './context/auth-context';
import { FavoritesProvider } from './context/favorites-context';
import { ensureTelegramViewport } from './telegram/bootstrap';

export function Layout() {
  const { theme } = useTheme();

  useEffect(() => {
    void ensureTelegramViewport();
  }, []);

  return (
    <div className="tma-shell bg-background text-foreground transition-colors duration-300" data-theme={theme}>
      <AuthProvider>
        <FavoritesProvider>
          <div className="tma-scroll">
            <Outlet />
          </div>
        </FavoritesProvider>
      </AuthProvider>
      <BottomNav />
    </div>
  );
}
