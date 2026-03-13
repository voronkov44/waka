import { productAssets } from '../assets';
import type { ContentBlock, FAQArticleDetail, FAQArticleSummary, FAQTopic, FAQTopicIcon, Product, ProductStatus } from '../types/domain';
import type { FAQArticleDetailDTO, FAQArticleSummaryDTO, FAQBlockDTO, FAQTopicDTO, ModelDTO, ModelTagDTO, PublicModelDTO } from './types';

const fallbackProductPhotos = [productAssets.pulse, productAssets.ultra];

function isObject(value: unknown): value is Record<string, unknown> {
  return Boolean(value) && typeof value === 'object' && !Array.isArray(value);
}

function readString(data: Record<string, unknown>, keys: string[]): string | undefined {
  for (const key of keys) {
    const value = data[key];
    if (typeof value === 'string' && value.trim().length > 0) {
      return value.trim();
    }
  }
  return undefined;
}

function readStringArray(data: Record<string, unknown>, keys: string[]): string[] | undefined {
  for (const key of keys) {
    const value = data[key];
    if (Array.isArray(value)) {
      const items = value.filter((item): item is string => typeof item === 'string' && item.trim().length > 0);
      if (items.length > 0) {
        return items.map((item) => item.trim());
      }
    }
  }
  return undefined;
}

function normalizeStatusFromTag(tag?: ModelTagDTO): ProductStatus {
  const text = `${tag?.key ?? ''} ${tag?.label ?? ''}`.toLowerCase();
  if (text.includes('new')) {
    return 'new';
  }
  if (text.includes('limit') || text.includes('exclusive')) {
    return 'limited';
  }
  return 'popular';
}

function mapTag(tag?: ModelTagDTO) {
  if (!tag) {
    return undefined;
  }

  return {
    key: tag.key,
    label: tag.label,
    bgColor: tag.bg_color,
    textColor: tag.text_color,
  };
}

function toDescription(value?: string): string {
  if (!value || value.trim().length === 0) {
    return 'Premium Waka model with signature performance.';
  }
  return value;
}

function toPhotoURL(value: string | undefined, id: number): string {
  if (value && value.trim().length > 0) {
    return value;
  }

  return fallbackProductPhotos[id % fallbackProductPhotos.length];
}

export function mapProduct(dto: PublicModelDTO | ModelDTO): Product {
  return {
    id: dto.id,
    name: dto.name,
    status: normalizeStatusFromTag(dto.tag),
    description: toDescription(dto.description),
    photoUrl: toPhotoURL(dto.photo_url, dto.id),
    puffsMax: dto.puffs_max,
    flavors: Array.isArray(dto.flavors) ? dto.flavors : [],
    priceCents: dto.price_cents ?? null,
    tag: mapTag(dto.tag),
  };
}

const iconFallbackCycle: FAQTopicIcon[] = ['Rocket', 'Info', 'Sparkles', 'Wrench', 'AlertCircle', 'Shield'];

function resolveTopicIcon(title: string, index: number): FAQTopicIcon {
  const normalized = title.toLowerCase();
  if (normalized.includes('start')) {
    return 'Rocket';
  }
  if (normalized.includes('device') || normalized.includes('info') || normalized.includes('spec')) {
    return 'Info';
  }
  if (normalized.includes('flavor')) {
    return 'Sparkles';
  }
  if (normalized.includes('maint') || normalized.includes('care')) {
    return 'Wrench';
  }
  if (normalized.includes('trouble') || normalized.includes('issue')) {
    return 'AlertCircle';
  }
  if (normalized.includes('safe') || normalized.includes('guideline') || normalized.includes('policy')) {
    return 'Shield';
  }
  return iconFallbackCycle[index % iconFallbackCycle.length];
}

export function mapFAQTopic(topic: FAQTopicDTO, articleCount: number, index: number): FAQTopic {
  return {
    id: topic.id,
    title: topic.title,
    description: articleCount > 0 ? `${articleCount} article${articleCount === 1 ? '' : 's'} available` : 'Explore answers and guides',
    icon: resolveTopicIcon(topic.title, index),
    articleCount,
  };
}

export function mapFAQArticleSummary(dto: FAQArticleSummaryDTO): FAQArticleSummary {
  return {
    id: dto.id,
    topicId: dto.topic_id,
    slug: dto.slug,
    title: dto.title,
    updatedAt: dto.updated_at,
  };
}

function normalizeCalloutVariant(value: string | undefined): ContentBlock['variant'] {
  const normalized = value?.toLowerCase();
  if (!normalized) {
    return 'info';
  }
  if (normalized.includes('warn') || normalized.includes('alert') || normalized.includes('danger')) {
    return 'warning';
  }
  if (normalized.includes('success') || normalized.includes('ok')) {
    return 'success';
  }
  return 'info';
}

function mapFAQBlock(block: FAQBlockDTO): ContentBlock | null {
  const type = block.type.toLowerCase();
  const data = isObject(block.data) ? block.data : {};

  if (type === 'divider') {
    return { type: 'divider' };
  }

  if (type === 'text') {
    const content = readString(data, ['content', 'text', 'body', 'title']);
    return content ? { type: 'text', content } : null;
  }

  if (type === 'image') {
    const url = readString(data, ['url', 'src', 'image_url', 'image']);
    return url ? { type: 'image', url } : null;
  }

  if (type === 'link') {
    const url = readString(data, ['url', 'href', 'link']);
    if (!url) {
      return null;
    }
    const content = readString(data, ['content', 'text', 'label', 'title']) ?? url;
    return { type: 'link', url, content };
  }

  if (type === 'bullets') {
    const items = readStringArray(data, ['items', 'list', 'bullets', 'points']);
    return items && items.length > 0 ? { type: 'bullets', items } : null;
  }

  if (type === 'callout') {
    const content = readString(data, ['content', 'text', 'body', 'title']);
    if (!content) {
      return null;
    }
    const variant = normalizeCalloutVariant(readString(data, ['variant', 'tone', 'style']));
    return { type: 'callout', content, variant };
  }

  // Unknown blocks are rendered as plain text when possible.
  const fallbackContent = readString(data, ['content', 'text', 'body']);
  if (fallbackContent) {
    return { type: 'text', content: fallbackContent };
  }

  return null;
}

export function mapFAQArticleDetail(dto: FAQArticleDetailDTO): FAQArticleDetail {
  const blocks = Array.isArray(dto.blocks)
    ? dto.blocks
        .slice()
        .sort((a, b) => a.sort - b.sort || a.id - b.id)
        .map(mapFAQBlock)
        .filter((block): block is ContentBlock => Boolean(block))
    : [];

  return {
    id: dto.id,
    topicId: dto.topic_id,
    slug: dto.slug,
    title: dto.title,
    updatedAt: dto.updated_at,
    contentBlocks: blocks,
  };
}
