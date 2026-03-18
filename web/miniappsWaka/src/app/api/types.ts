export interface TelegramAuthRequestDTO {
  tg_id: number;
  username?: string;
  first_name?: string;
  last_name?: string;
  photo_url?: string;
}

export interface TokenResponseDTO {
  token: string;
}

export interface MeResponseDTO {
  id: number;
  tg_id: number;
  username?: string;
  first_name?: string;
  last_name?: string;
  photo_url?: string;
  created_at: string;
  updated_at: string;
}

export interface ModelTagDTO {
  key?: string;
  label: string;
  bg_color: string;
  text_color: string;
}

export interface ShowcaseTagDTO {
  label: string;
  bg_color: string;
  text_color: string;
  outlined: boolean;
}

export interface ShowcaseItemDTO {
  id: number;
  tag: ShowcaseTagDTO;
  title: string;
  description?: string;
  model_id: number;
  photo_url?: string;
  sort: number;
  created_at: string;
  updated_at: string;
}

export interface PublicModelDTO {
  id: number;
  name: string;
  description?: string;
  photo_url?: string;
  tag?: ModelTagDTO;
  puffs_max: number;
  flavors: string[];
  price_cents?: number;
}

export interface ModelDTO {
  id: number;
  name: string;
  status: string;
  description?: string;
  photo_key?: string;
  photo_url?: string;
  tag?: ModelTagDTO;
  puffs_max: number;
  flavors: string[];
  price_cents?: number;
  created_at: string;
  updated_at: string;
}

export interface ListModelsResponseDTO<TModel> {
  items: TModel[];
  limit: number;
  offset: number;
}

export interface ListShowcaseItemsResponseDTO {
  items: ShowcaseItemDTO[];
  limit: number;
  offset: number;
}

export interface FAQTopicDTO {
  id: number;
  title: string;
  sort: number;
  is_active: boolean;
  created_at: string;
  updated_at: string;
}

export interface FAQArticleSummaryDTO {
  id: number;
  topic_id: number;
  slug: string;
  title: string;
  updated_at: string;
}

export interface FAQBlockDTO {
  id: number;
  article_id: number;
  sort: number;
  type: string;
  data: unknown;
  created_at: string;
  updated_at: string;
}

export interface FAQArticleDTO {
  id: number;
  topic_id: number;
  slug: string;
  title: string;
  status: string;
  channel: string;
  search_text?: string;
  published_at?: string;
  created_at: string;
  updated_at: string;
}

export interface FAQArticleDetailDTO {
  id: number;
  topic_id: number;
  slug: string;
  title: string;
  status: string;
  channel: string;
  search_text?: string;
  published_at?: string;
  created_at: string;
  updated_at: string;
  blocks: FAQBlockDTO[];
}
