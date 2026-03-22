import { useCallback, useEffect, useState } from 'react';
import { apiClient } from '../api/client';
import { mapFAQArticleSummary } from '../api/mappers';
import { ApiError } from '../api/http';
import type { FAQArticleSummary, FAQTopic } from '../types/domain';
import { i18nText } from '../../shared/i18n';

interface UseFAQArticlesResult {
  topic: FAQTopic | null;
  articles: FAQArticleSummary[];
  isLoading: boolean;
  error: string | null;
  notFound: boolean;
  refresh: () => Promise<void>;
}

export function useFAQArticles(topicID?: number): UseFAQArticlesResult {
  const [topic, setTopic] = useState<FAQTopic | null>(null);
  const [articles, setArticles] = useState<FAQArticleSummary[]>([]);
  const [isLoading, setIsLoading] = useState(Boolean(topicID));
  const [error, setError] = useState<string | null>(null);
  const [notFound, setNotFound] = useState(false);

  const refresh = useCallback(async () => {
    if (!topicID || Number.isNaN(topicID)) {
      setTopic(null);
      setArticles([]);
      setIsLoading(false);
      setNotFound(true);
      return;
    }

    setIsLoading(true);
    setError(null);
    setNotFound(false);

    try {
      const [topicsResponse, articlesResponse] = await Promise.all([
        apiClient.listFAQTopics(),
        apiClient.listFAQArticlesByTopic(topicID),
      ]);

      const currentTopic = topicsResponse.find((item) => item.id === topicID);
      if (!currentTopic) {
        setTopic(null);
        setArticles([]);
        setNotFound(true);
        return;
      }

      setTopic({
        id: currentTopic.id,
        title: currentTopic.title,
        description: '',
        icon: 'Info',
        articleCount: articlesResponse.length,
      });
      setArticles(articlesResponse.map(mapFAQArticleSummary));
    } catch (err) {
      if (err instanceof ApiError && err.status === 404) {
        setNotFound(true);
      } else {
        setError(err instanceof Error ? err.message : i18nText('errors.loadFaqArticles'));
      }
    } finally {
      setIsLoading(false);
    }
  }, [topicID]);

  useEffect(() => {
    void refresh();
  }, [refresh]);

  return { topic, articles, isLoading, error, notFound, refresh };
}
