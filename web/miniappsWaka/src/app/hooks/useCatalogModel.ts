import { useCallback, useEffect, useState } from 'react';
import { ApiError } from '../api/http';
import { apiClient } from '../api/client';
import { mapProduct } from '../api/mappers';
import type { Product } from '../types/domain';

export function useCatalogModel(modelID?: number) {
  const [product, setProduct] = useState<Product | null>(null);
  const [isLoading, setIsLoading] = useState(Boolean(modelID));
  const [error, setError] = useState<string | null>(null);
  const [notFound, setNotFound] = useState(false);

  const refresh = useCallback(async () => {
    if (!modelID || Number.isNaN(modelID)) {
      setIsLoading(false);
      setProduct(null);
      setNotFound(true);
      return;
    }

    setIsLoading(true);
    setError(null);
    setNotFound(false);
    try {
      const response = await apiClient.getCatalogModel(modelID);
      setProduct(mapProduct(response));
    } catch (err) {
      if (err instanceof ApiError && err.status === 404) {
        setNotFound(true);
        setProduct(null);
      } else {
        setError(err instanceof Error ? err.message : 'Failed to load model');
      }
    } finally {
      setIsLoading(false);
    }
  }, [modelID]);

  useEffect(() => {
    void refresh();
  }, [refresh]);

  return {
    product,
    isLoading,
    error,
    notFound,
    refresh,
  };
}
