import { useCallback, useEffect, useState } from 'react';
import { apiClient } from '../api/client';
import { mapFAQTopic } from '../api/mappers';
import type { FAQTopic } from '../types/domain';
import { i18nText } from '../../shared/i18n';

export function useFAQTopics() {
  const [topics, setTopics] = useState<FAQTopic[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  const refresh = useCallback(async () => {
    setIsLoading(true);
    setError(null);
    try {
      const topicsResponse = await apiClient.listFAQTopics();

      const articleCounts = await Promise.all(
        topicsResponse.map(async (topic) => {
          try {
            const articles = await apiClient.listFAQArticlesByTopic(topic.id);
            return { topicID: topic.id, count: articles.length };
          } catch {
            return { topicID: topic.id, count: 0 };
          }
        }),
      );

      const countByTopic = new Map(articleCounts.map((item) => [item.topicID, item.count]));
      const mappedTopics = topicsResponse.map((topic, index) =>
        mapFAQTopic(topic, countByTopic.get(topic.id) ?? 0, index),
      );

      setTopics(mappedTopics);
    } catch (err) {
      setError(err instanceof Error ? err.message : i18nText('errors.loadFaqTopics'));
    } finally {
      setIsLoading(false);
    }
  }, []);

  useEffect(() => {
    void refresh();
  }, [refresh]);

  return {
    topics,
    isLoading,
    error,
    refresh,
  };
}
