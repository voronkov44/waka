import { http } from './http';
import type {
  FAQArticleDetailDTO,
  FAQArticleSummaryDTO,
  FAQTopicDTO,
  ListModelsResponseDTO,
  MeResponseDTO,
  ModelDTO,
  PublicModelDTO,
  TelegramAuthRequestDTO,
  TokenResponseDTO,
} from './types';

function withQuery(path: string, params: Record<string, string | number | undefined>): string {
  const query = new URLSearchParams();
  Object.entries(params).forEach(([key, value]) => {
    if (value === undefined || value === null || value === '') {
      return;
    }
    query.set(key, String(value));
  });
  const queryString = query.toString();
  return queryString ? `${path}?${queryString}` : path;
}

export const apiClient = {
  loginTelegram(payload: TelegramAuthRequestDTO) {
    return http.post<TokenResponseDTO>('/api/auth/telegram', payload, { auth: false });
  },

  getCurrentUser() {
    return http.get<MeResponseDTO>('/api/auth/me');
  },

  listCatalogModels(limit = 100, offset = 0) {
    return http.get<ListModelsResponseDTO<PublicModelDTO>>(withQuery('/api/catalog/models', { limit, offset }), {
      auth: false,
    });
  },

  getCatalogModel(id: number) {
    return http.get<PublicModelDTO>(`/api/catalog/models/${id}`, { auth: false });
  },

  listFavorites(limit = 100, offset = 0) {
    return http.get<ListModelsResponseDTO<ModelDTO>>(withQuery('/api/favorites', { limit, offset }));
  },

  addFavorite(modelID: number) {
    return http.post<null>(`/api/favorites/${modelID}`);
  },

  removeFavorite(modelID: number) {
    return http.delete<null>(`/api/favorites/${modelID}`);
  },

  listFAQTopics() {
    return http.get<FAQTopicDTO[]>('/api/faq/topics', { auth: false });
  },

  listFAQArticlesByTopic(topicID: number, channel = 'miniapp') {
    return http.get<FAQArticleSummaryDTO[]>(
      withQuery(`/api/faq/topics/${topicID}/articles`, { channel }),
      { auth: false },
    );
  },

  getFAQArticle(id: number) {
    return http.get<FAQArticleDetailDTO>(`/api/faq/articles/${id}`, { auth: false });
  },
};
