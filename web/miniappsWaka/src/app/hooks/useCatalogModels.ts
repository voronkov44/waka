import { useCallback, useEffect, useState } from 'react';
import { apiClient } from '../api/client';
import { mapProduct } from '../api/mappers';
import type { Product } from '../types/domain';
import { i18nText } from '../../shared/i18n';

export function useCatalogModels() {
  const [products, setProducts] = useState<Product[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  const refresh = useCallback(async () => {
    setIsLoading(true);
    setError(null);
    try {
      const response = await apiClient.listCatalogModels();
      setProducts(response.items.map(mapProduct));
    } catch (err) {
      setError(err instanceof Error ? err.message : i18nText('errors.loadCatalog'));
    } finally {
      setIsLoading(false);
    }
  }, []);

  useEffect(() => {
    void refresh();
  }, [refresh]);

  return {
    products,
    isLoading,
    error,
    refresh,
  };
}
