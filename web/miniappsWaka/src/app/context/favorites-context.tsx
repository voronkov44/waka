import { createContext, useCallback, useContext, useEffect, useMemo, useState, type ReactNode } from 'react';
import { apiClient } from '../api/client';
import { ApiError } from '../api/http';
import { mapProduct } from '../api/mappers';
import { useAuthContext } from './auth-context';
import type { Product } from '../types/domain';

interface FavoritesContextValue {
  favorites: number[];
  favoriteProducts: Product[];
  isLoading: boolean;
  error: string | null;
  isFavorite: (productID: number) => boolean;
  toggleFavorite: (productID: number) => Promise<void>;
  refreshFavorites: () => Promise<void>;
}

const FavoritesContext = createContext<FavoritesContextValue | null>(null);

export function FavoritesProvider({ children }: { children: ReactNode }) {
  const { isLoading: authLoading, token } = useAuthContext();
  const [favoriteProducts, setFavoriteProducts] = useState<Product[]>([]);
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [pendingIDs, setPendingIDs] = useState<Record<number, boolean>>({});

  const favorites = useMemo(() => favoriteProducts.map((product) => product.id), [favoriteProducts]);

  const refreshFavorites = useCallback(async () => {
    if (!token) {
      setFavoriteProducts([]);
      setError(null);
      return;
    }

    setIsLoading(true);
    setError(null);
    try {
      const response = await apiClient.listFavorites();
      const nextFavorites = response.items.map(mapProduct);
      setFavoriteProducts(nextFavorites);
    } catch (err) {
      if (err instanceof ApiError && err.status === 401) {
        setFavoriteProducts([]);
        setError('Please reopen the mini app to refresh your session.');
        return;
      }
      setError(err instanceof Error ? err.message : 'Failed to load favorites');
    } finally {
      setIsLoading(false);
    }
  }, [token]);

  useEffect(() => {
    if (authLoading) {
      return;
    }

    if (!token) {
      setFavoriteProducts([]);
      setError(null);
      return;
    }

    void refreshFavorites();
  }, [authLoading, token, refreshFavorites]);

  const isFavorite = useCallback(
    (productID: number) => favoriteProducts.some((product) => product.id === productID),
    [favoriteProducts],
  );

  const toggleFavorite = useCallback(
    async (productID: number) => {
      if (!token) {
        return;
      }
      if (pendingIDs[productID]) {
        return;
      }

      const currentlyFavorite = favoriteProducts.some((product) => product.id === productID);
      const previousFavorites = favoriteProducts;

      setPendingIDs((prev) => ({ ...prev, [productID]: true }));
      setFavoriteProducts((prev) =>
        currentlyFavorite ? prev.filter((product) => product.id !== productID) : prev,
      );

      try {
        if (currentlyFavorite) {
          await apiClient.removeFavorite(productID);
        } else {
          await apiClient.addFavorite(productID);
        }
        await refreshFavorites();
      } catch (err) {
        setFavoriteProducts(previousFavorites);
        setError(err instanceof Error ? err.message : 'Failed to update favorites');
      } finally {
        setPendingIDs((prev) => {
          const next = { ...prev };
          delete next[productID];
          return next;
        });
      }
    },
    [token, pendingIDs, favoriteProducts, refreshFavorites],
  );

  const contextValue = useMemo<FavoritesContextValue>(
    () => ({
      favorites,
      favoriteProducts,
      isLoading,
      error,
      isFavorite,
      toggleFavorite,
      refreshFavorites,
    }),
    [favorites, favoriteProducts, isLoading, error, isFavorite, toggleFavorite, refreshFavorites],
  );

  return <FavoritesContext.Provider value={contextValue}>{children}</FavoritesContext.Provider>;
}

export function useFavoritesContext(): FavoritesContextValue {
  const context = useContext(FavoritesContext);
  if (!context) {
    throw new Error('useFavoritesContext must be used inside FavoritesProvider');
  }
  return context;
}
