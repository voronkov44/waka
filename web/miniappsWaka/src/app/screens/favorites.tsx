import { useNavigate } from 'react-router';
import { Heart } from 'lucide-react';
import { useFavorites } from '../hooks/useFavorites';
import { ProductCard } from '../components/product-card';
import { EmptyState } from '../components/empty-state';
import { useAuth } from '../hooks/useAuth';

export function Favorites() {
  const navigate = useNavigate();
  const { isLoading: authLoading, isAuthenticated, hasTelegramContext } = useAuth();
  const { favoriteProducts, toggleFavorite, isFavorite, isLoading, error } = useFavorites();

  if (authLoading || isLoading) {
    return (
      <div className="min-h-screen flex items-center justify-center text-muted-foreground">
        Loading favorites...
      </div>
    );
  }

  if (!isAuthenticated) {
    return (
      <div className="min-h-screen px-6 pt-16">
        <EmptyState
          icon={Heart}
          title="Favorites require Telegram login"
          description={
            hasTelegramContext
              ? 'Please wait while your Telegram profile is being linked.'
              : 'Open this mini app from Telegram to sync your favorites.'
          }
          action={{
            label: 'Go Home',
            onClick: () => navigate('/'),
          }}
        />
      </div>
    );
  }

  return (
    <div className="min-h-screen pb-32">
      <div className="sticky top-0 z-40 bg-background/80 backdrop-blur-[32px] border-b border-border/40 pt-safe">
        <div className="px-6 py-6 pb-4">
          <div className="flex items-center gap-3 mb-1">
            <Heart className="w-6 h-6 text-foreground fill-current drop-shadow-sm" />
            <h1 className="text-4xl font-extrabold tracking-tighter leading-none">Favorites</h1>
          </div>
          {favoriteProducts.length > 0 && (
            <p className="text-[11px] font-bold tracking-[0.1em] uppercase text-muted-foreground mt-3">
              {favoriteProducts.length} {favoriteProducts.length === 1 ? 'device' : 'devices'} saved
            </p>
          )}
        </div>
      </div>

      {error && (
        <div className="px-6 pt-6">
          <div className="rounded-2xl border border-destructive/30 bg-destructive/10 p-4 text-sm text-destructive">
            {error}
          </div>
        </div>
      )}

      {favoriteProducts.length > 0 ? (
        <div className="px-6 pt-6">
          <div className="flex flex-col items-center gap-4">
            {favoriteProducts.map((product) => (
              <ProductCard
                key={product.id}
                product={product}
                isFavorite={isFavorite(product.id)}
                compact
                className="w-full max-w-[300px]"
                onToggleFavorite={() => {
                  void toggleFavorite(product.id);
                }}
              />
            ))}
          </div>
        </div>
      ) : (
        <div className="px-6 pt-16">
          <EmptyState
            icon={Heart}
            title="No favorites yet"
            description="Start exploring our catalog and save your premium devices here"
            action={{
              label: 'Browse Catalog',
              onClick: () => navigate('/catalog'),
            }}
          />
        </div>
      )}
    </div>
  );
}
