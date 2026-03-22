import { useMemo, useState } from 'react';
import { SearchBar } from '../components/search-bar';
import { FilterChip } from '../components/filter-chip';
import { ProductCard } from '../components/product-card';
import { useFavorites } from '../hooks/useFavorites';
import { useCatalogModels } from '../hooks/useCatalogModels';
import { resolveI18nText, useI18n } from '../../shared/i18n';

type PuffFilter = 'all' | 'up-to-10k' | '10k-30k' | '30k+';

export function Catalog() {
  const { t, tp, localeCode } = useI18n();
  const [searchQuery, setSearchQuery] = useState('');
  const [puffFilter, setPuffFilter] = useState<PuffFilter>('all');
  const { isFavorite, toggleFavorite } = useFavorites();
  const { products, isLoading, error } = useCatalogModels();
  const localizedError = resolveI18nText(error, t);

  const filteredProducts = useMemo(() => {
    return products.filter((product) => {
      const normalizedQuery = searchQuery.toLowerCase();
      const matchesSearch =
        product.name.toLowerCase().includes(normalizedQuery) ||
        product.description.toLowerCase().includes(normalizedQuery) ||
        product.flavors.some((flavor) => flavor.toLowerCase().includes(normalizedQuery));

      let matchesPuff = true;
      if (puffFilter === 'up-to-10k') {
        matchesPuff = product.puffsMax < 10000;
      } else if (puffFilter === '10k-30k') {
        matchesPuff = product.puffsMax >= 10000 && product.puffsMax < 30000;
      } else if (puffFilter === '30k+') {
        matchesPuff = product.puffsMax >= 30000;
      }

      return matchesSearch && matchesPuff;
    });
  }, [products, searchQuery, puffFilter]);

  return (
    <div className="min-h-screen pb-32">
      <div className="sticky top-0 z-40 bg-background/80 backdrop-blur-[32px] border-b border-border/40 pt-safe">
        <div className="px-6 py-6 pb-4">
          <h1 className="text-4xl font-extrabold mb-6 tracking-tighter leading-none">{t('catalog.title')}</h1>
          <SearchBar value={searchQuery} onChange={setSearchQuery} placeholder={t('search.catalogPlaceholder')} />
        </div>

        <div className="px-6 pb-5">
          <div className="flex gap-2.5 overflow-x-auto scrollbar-hide py-1 pl-1 pr-2">
            <FilterChip label={t('catalog.filterAll')} active={puffFilter === 'all'} onClick={() => setPuffFilter('all')} />
            <FilterChip
              label={t('catalog.filterUpTo10k')}
              active={puffFilter === 'up-to-10k'}
              onClick={() => setPuffFilter('up-to-10k')}
            />
            <FilterChip
              label={t('catalog.filter10kTo30k')}
              active={puffFilter === '10k-30k'}
              onClick={() => setPuffFilter('10k-30k')}
            />
            <FilterChip label={t('catalog.filter30kPlus')} active={puffFilter === '30k+'} onClick={() => setPuffFilter('30k+')} />
          </div>
        </div>
      </div>

      <div className="px-6 pt-6">
        {isLoading && <p className="text-sm text-muted-foreground py-8 text-center">{t('catalog.loading')}</p>}

        {localizedError && (
          <div className="rounded-2xl border border-destructive/30 bg-destructive/10 p-4 text-sm text-destructive">
            {localizedError}
          </div>
        )}

        {!isLoading && !localizedError && filteredProducts.length > 0 && (
          <div className="flex flex-col gap-5">
            {filteredProducts.map((product) => (
              <ProductCard
                key={product.id}
                product={product}
                isFavorite={isFavorite(product.id)}
                onToggleFavorite={() => {
                  void toggleFavorite(product.id);
                }}
              />
            ))}
          </div>
        )}

        {!isLoading && !localizedError && filteredProducts.length === 0 && (
          <div className="text-center py-20">
            <p className="text-[11px] font-bold tracking-[0.1em] uppercase text-muted-foreground mb-2">
              {t('catalog.emptyTitle')}
            </p>
            <p className="text-sm font-medium text-muted-foreground/70">{t('catalog.emptyDescription')}</p>
          </div>
        )}
      </div>

      {!isLoading && !localizedError && filteredProducts.length > 0 && (
        <div className="px-6 pt-8 pb-4">
          <p className="text-[10px] font-bold tracking-[0.2em] uppercase text-muted-foreground text-center">
            {t('catalog.showingCount', {
              count: filteredProducts.length.toLocaleString(localeCode),
              unit: tp('nouns.device', filteredProducts.length),
            })}
          </p>
        </div>
      )}
    </div>
  );
}
