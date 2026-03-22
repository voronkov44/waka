import { useEffect, useMemo, useState } from 'react';
import { useParams, Link, useNavigate } from 'react-router';
import { ChevronDown, ChevronLeft, Heart, Zap } from 'lucide-react';
import { useFavorites } from '../hooks/useFavorites';
import { FlavorChip } from '../components/flavor-chip';
import { ProductStatusBadge } from '../components/product-status-badge';
import { useCatalogModel } from '../hooks/useCatalogModel';
import { useCatalogModels } from '../hooks/useCatalogModels';
import { ModelImage } from '../components/model-image';
import { resolveI18nText, useI18n } from '../../shared/i18n';

function normalizePriceCents(value: unknown): number | null {
  if (typeof value === 'number') {
    return Number.isFinite(value) ? value : null;
  }

  if (typeof value === 'string') {
    const trimmed = value.trim();
    if (trimmed.length === 0 || trimmed === '-' || trimmed === '—') {
      return null;
    }

    const parsed = Number(trimmed);
    return Number.isFinite(parsed) ? parsed : null;
  }

  return null;
}

function formatPrice(cents: number) {
  return `$${(cents / 100).toFixed(2)}`;
}

export function ProductDetail() {
  const { t, tp, localeCode } = useI18n();
  const { id } = useParams<{ id: string }>();
  const productID = Number(id);
  const [isFlavorsOpen, setIsFlavorsOpen] = useState(false);
  const navigate = useNavigate();
  const { isFavorite, toggleFavorite } = useFavorites();
  const { product, isLoading, error, notFound } = useCatalogModel(productID);
  const { products } = useCatalogModels();
  const localizedError = resolveI18nText(error, t);

  const handleGoBack = () => {
    if (window.history.length > 1) {
      navigate(-1);
      return;
    }
    navigate('/catalog');
  };

  const relatedProducts = useMemo(() => {
    if (!product) {
      return [];
    }
    return products.filter((item) => item.id !== product.id).slice(0, 3);
  }, [products, product]);

  useEffect(() => {
    setIsFlavorsOpen(false);
  }, [product?.id]);

  const flavorCount = Array.isArray(product?.flavors) ? product.flavors.length : 0;
  const hasFlavors = flavorCount > 0;
  const normalizedPriceCents = normalizePriceCents(product?.priceCents);
  const priceLabel = normalizedPriceCents !== null ? formatPrice(normalizedPriceCents) : null;
  const hasPrice = priceLabel !== null;
  const productDescription =
    typeof product?.description === 'string' && product.description.trim().length > 0
      ? product.description
      : t('common.defaultProductDescription');
  const flavorCountLabel = tp('nouns.option', flavorCount);

  if (isLoading) {
    return (
      <div className="min-h-screen flex items-center justify-center text-muted-foreground">
        {t('product.loading')}
      </div>
    );
  }

  if (localizedError) {
    return (
      <div className="min-h-screen flex items-center justify-center p-6">
        <div className="rounded-2xl border border-destructive/30 bg-destructive/10 p-4 text-sm text-destructive">
          {localizedError}
        </div>
      </div>
    );
  }

  if (notFound || !product) {
    return (
      <div className="min-h-screen flex items-center justify-center">
        <div className="text-center">
          <p className="text-xl mb-4">{t('product.notFound')}</p>
          <Link to="/catalog" className="text-foreground underline">
            {t('actions.backToCatalog')}
          </Link>
        </div>
      </div>
    );
  }

  return (
    <div className="min-h-screen pb-32">
      <div className="sticky top-0 z-40 bg-background/80 backdrop-blur-[32px] border-b border-border/40 pt-safe">
        <div className="flex items-center justify-between px-6 py-4">
          <button
            type="button"
            onClick={handleGoBack}
            className="w-12 h-12 rounded-[18px] bg-card border border-border/50 shadow-sm flex items-center justify-center hover:scale-105 transition-all duration-300 group"
          >
            <ChevronLeft className="w-6 h-6 text-foreground group-hover:text-foreground/80 transition-colors" />
          </button>
          <div className="text-sm font-bold tracking-[0.2em] uppercase text-foreground">{t('product.detailsTitle')}</div>
          <button
            type="button"
            onClick={() => {
              void toggleFavorite(product.id);
            }}
            className={`w-12 h-12 rounded-[18px] flex items-center justify-center border transition-all duration-300 shadow-sm ${
              isFavorite(product.id)
                ? 'bg-foreground border-foreground text-background scale-105 shadow-md'
                : 'bg-card border-border/50 text-foreground hover:border-foreground/30 hover:scale-105'
            }`}
          >
            <Heart className={`w-6 h-6 ${isFavorite(product.id) ? 'fill-current' : ''}`} />
          </button>
        </div>
      </div>

      <div className="px-6 pt-6">
        <div className="relative aspect-[4/5] rounded-[40px] overflow-hidden bg-card border border-border/50 shadow-sm dark:shadow-none flex items-center justify-center">
          <div className="absolute inset-0 bg-[radial-gradient(circle_at_50%_0%,var(--glow-primary),transparent_70%)] pointer-events-none mix-blend-overlay" />
          <ModelImage
            preset="details"
            src={product.photoUrl}
            alt={product.name}
            className="p-4"
            imageClassName="translate-y-[11%]"
          />
          <div className="absolute top-6 left-6 z-10">
            <ProductStatusBadge status={product.status} tag={product.tag} />
          </div>
        </div>
      </div>

      <div className="px-6 pt-10">
        <h1 className="text-4xl font-extrabold tracking-tighter mb-4 leading-none">{product.name}</h1>
        <p className="text-[13px] font-medium leading-relaxed text-muted-foreground mb-10">{productDescription}</p>

        <div className="bg-card border border-border/50 rounded-[32px] p-8 mb-8 shadow-sm">
          <div className="flex items-center gap-3 mb-6">
            <div className="w-10 h-10 rounded-[14px] bg-foreground/10 flex items-center justify-center">
              <Zap className="w-5 h-5 text-foreground" />
            </div>
            <h3 className="font-bold tracking-tight text-lg">{t('product.specifications')}</h3>
          </div>
          <div className="grid grid-cols-2 gap-6">
            <div>
              <p className="text-[10px] font-bold tracking-[0.2em] uppercase text-muted-foreground mb-2">{t('product.maxPuffs')}</p>
              <p className="font-extrabold text-xl tracking-tight text-foreground">
                {product.puffsMax.toLocaleString(localeCode)}
              </p>
            </div>
            <div>
              <p className="text-[10px] font-bold tracking-[0.2em] uppercase text-muted-foreground mb-2">{t('common.flavors')}</p>
              <p className="font-extrabold text-xl tracking-tight text-foreground">
                {t('product.flavorCountWithUnit', {
                  count: flavorCount.toLocaleString(localeCode),
                  unit: flavorCountLabel,
                })}
              </p>
            </div>
          </div>
        </div>

        {hasFlavors && (
          <div className="mb-10">
            <button
              type="button"
              aria-expanded={isFlavorsOpen}
              aria-controls={`flavors-panel-${product.id}`}
              onClick={() => setIsFlavorsOpen((prev) => !prev)}
              className="w-full rounded-[28px] border border-border/50 bg-card px-6 py-4 shadow-sm flex items-center justify-between transition-all duration-300 hover:border-foreground/30"
            >
              <div className="flex items-center gap-3">
                <span className="text-[11px] font-bold tracking-[0.2em] uppercase text-muted-foreground">{t('product.availableFlavors')}</span>
                <span className="min-w-8 rounded-full border border-border/60 bg-background px-2 py-0.5 text-[11px] font-bold text-foreground">
                  {flavorCount.toLocaleString(localeCode)}
                </span>
              </div>
              <ChevronDown
                className={`h-4 w-4 text-muted-foreground transition-transform duration-300 ${
                  isFlavorsOpen ? 'rotate-180' : 'rotate-0'
                }`}
              />
            </button>

            <div
              id={`flavors-panel-${product.id}`}
              aria-hidden={!isFlavorsOpen}
              className={`grid overflow-hidden transition-all duration-300 ease-[cubic-bezier(0.2,0.8,0.2,1)] ${
                isFlavorsOpen ? 'mt-4 grid-rows-[1fr] opacity-100' : 'grid-rows-[0fr] opacity-0'
              }`}
            >
              <div className="min-h-0">
                <div className="flex flex-wrap gap-2.5">
                  {product.flavors.map((flavor) => (
                    <FlavorChip key={flavor} flavor={flavor} />
                  ))}
                </div>
              </div>
            </div>
          </div>
        )}

        {hasPrice && (
          <div className="bg-card border border-border/50 rounded-[32px] p-8 mb-10 shadow-sm relative overflow-hidden">
            <div className="absolute top-0 right-0 w-32 h-32 bg-foreground/5 rounded-full blur-2xl -translate-y-1/2 translate-x-1/2 pointer-events-none" />
            <div className="relative z-10">
              <p className="text-[10px] font-bold tracking-[0.2em] uppercase text-muted-foreground mb-2">{t('product.price')}</p>
              <p className="text-4xl font-extrabold tracking-tighter text-foreground">{priceLabel}</p>
            </div>
          </div>
        )}

        {relatedProducts.length > 0 && (
          <div>
            <div className="flex items-center justify-between mb-6 pl-2">
              <h3 className="text-[11px] font-bold tracking-[0.2em] uppercase text-muted-foreground">{t('product.relatedModels')}</h3>
            </div>
            <div className="flex gap-4 overflow-x-auto pb-6 -mx-6 px-6 scrollbar-hide snap-x">
              {relatedProducts.map((related) => {
                const flavorCount = Array.isArray(related.flavors) ? related.flavors.length : 0;
                const description =
                  typeof related.description === 'string' && related.description.trim().length > 0
                    ? related.description.trim()
                    : t('common.defaultProductDescription');

                return (
                  <Link
                    key={related.id}
                    to={`/product/${related.id}`}
                    className="group flex-shrink-0 w-[240px] bg-card border border-border/50 rounded-[28px] overflow-hidden hover:border-foreground/30 hover:shadow-lg transition-all duration-500 snap-start shadow-sm"
                  >
                    <div className="aspect-[11/10] bg-gradient-to-b from-foreground/5 to-transparent relative overflow-hidden flex items-center justify-center p-3">
                      <ModelImage preset="related" src={related.photoUrl} alt={related.name} />
                    </div>
                    <div className="p-5 border-t border-border/40">
                      <p className="text-[18px] font-bold tracking-tight leading-tight text-foreground break-words line-clamp-2">
                        {related.name}
                      </p>
                      {description.length > 0 && (
                        <p className="mb-2 mt-2 text-[12px] font-medium leading-relaxed text-muted-foreground line-clamp-2">
                          {description}
                        </p>
                      )}

                      <div className="mt-2 flex items-center gap-1.5">
                        <div className="flex flex-col">
                          <span className="mb-1 text-[8px] font-bold uppercase tracking-[0.2em] text-muted-foreground">
                            {t('common.capacity')}
                          </span>
                          <span className="text-[11px] font-bold tracking-tight text-foreground">
                            {related.puffsMax.toLocaleString(localeCode)} {t('common.puffs')}
                          </span>
                        </div>
                        <div className="h-8 w-[1px] bg-border/50"></div>
                        <div className="flex flex-col">
                          <span className="mb-1 text-[8px] font-bold uppercase tracking-[0.2em] text-muted-foreground">
                            {t('common.flavors')}
                          </span>
                          <span className="text-[11px] font-bold tracking-tight text-foreground">
                            {flavorCount.toLocaleString(localeCode)} {tp('nouns.option', flavorCount)}
                          </span>
                        </div>
                      </div>
                    </div>
                  </Link>
                );
              })}
            </div>
          </div>
        )}
      </div>
    </div>
  );
}
