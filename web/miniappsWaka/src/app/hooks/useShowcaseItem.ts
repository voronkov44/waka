import { useCallback, useEffect, useState } from 'react';
import { apiClient } from '../api/client';
import { mapShowcaseItem } from '../api/mappers';
import type { ShowcaseItem } from '../types/domain';

export function useShowcaseItem() {
  const [showcaseItem, setShowcaseItem] = useState<ShowcaseItem | null>(null);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  const refresh = useCallback(async () => {
    setIsLoading(true);
    setError(null);

    try {
      const response = await apiClient.listShowcaseItems(1, 0);
      const firstItem = response.items[0];
      setShowcaseItem(firstItem ? mapShowcaseItem(firstItem) : null);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to load showcase');
      setShowcaseItem(null);
    } finally {
      setIsLoading(false);
    }
  }, []);

  useEffect(() => {
    void refresh();
  }, [refresh]);

  return {
    showcaseItem,
    isLoading,
    error,
    refresh,
  };
}

