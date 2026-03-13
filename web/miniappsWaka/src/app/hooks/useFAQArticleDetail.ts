import { useCallback, useEffect, useState } from 'react';
import { ApiError } from '../api/http';
import { apiClient } from '../api/client';
import { mapFAQArticleDetail } from '../api/mappers';
import type { FAQArticleDetail } from '../types/domain';

export function useFAQArticleDetail(articleID?: number, topicID?: number) {
  const [article, setArticle] = useState<FAQArticleDetail | null>(null);
  const [isLoading, setIsLoading] = useState(Boolean(articleID));
  const [error, setError] = useState<string | null>(null);
  const [notFound, setNotFound] = useState(false);

  const refresh = useCallback(async () => {
    if (!articleID || Number.isNaN(articleID)) {
      setArticle(null);
      setNotFound(true);
      setIsLoading(false);
      return;
    }

    setIsLoading(true);
    setError(null);
    setNotFound(false);

    try {
      const response = await apiClient.getFAQArticle(articleID);
      const mapped = mapFAQArticleDetail(response);

      if (topicID && mapped.topicId !== topicID) {
        setNotFound(true);
        setArticle(null);
      } else {
        setArticle(mapped);
      }
    } catch (err) {
      if (err instanceof ApiError && err.status === 404) {
        setNotFound(true);
      } else {
        setError(err instanceof Error ? err.message : 'Failed to load article');
      }
    } finally {
      setIsLoading(false);
    }
  }, [articleID, topicID]);

  useEffect(() => {
    void refresh();
  }, [refresh]);

  return {
    article,
    isLoading,
    error,
    notFound,
    refresh,
  };
}
