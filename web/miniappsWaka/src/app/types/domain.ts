export type ProductStatus = 'new' | 'popular' | 'limited';

export type FAQTopicIcon = 'Rocket' | 'Info' | 'Sparkles' | 'Wrench' | 'AlertCircle' | 'Shield';

export interface ProductTag {
  key?: string;
  label: string;
  bgColor: string;
  textColor: string;
}

export interface Product {
  id: number;
  name: string;
  status: ProductStatus;
  description: string;
  photoUrl: string;
  puffsMax: number;
  flavors: string[];
  priceCents: number | null;
  tag?: ProductTag;
}

export interface FAQTopic {
  id: number;
  title: string;
  description: string;
  icon: FAQTopicIcon;
  articleCount: number;
}

export interface FAQArticleSummary {
  id: number;
  topicId: number;
  slug: string;
  title: string;
  updatedAt: string;
}

export type ContentBlockType = 'text' | 'image' | 'link' | 'bullets' | 'divider' | 'callout';

export interface ContentBlock {
  type: ContentBlockType;
  content?: string;
  items?: string[];
  url?: string;
  variant?: 'info' | 'warning' | 'success';
}

export interface FAQArticleDetail {
  id: number;
  topicId: number;
  slug: string;
  title: string;
  updatedAt: string;
  contentBlocks: ContentBlock[];
}
