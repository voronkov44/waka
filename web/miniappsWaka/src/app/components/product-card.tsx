import { Heart } from 'lucide-react';
import { Link } from 'react-router';
import type { Product } from '../types/domain';
import { ProductStatusBadge } from './product-status-badge';
import { ModelImage } from './model-image';

interface ProductCardProps {
  product: Product;
  isFavorite: boolean;
  onToggleFavorite: () => void;
  compact?: boolean;
  className?: string;
}

function formatPrice(cents: number | null): string | null {
  if (cents === null || Number.isNaN(cents)) {
    return null;
  }

  const rubles = cents / 100;
  const hasFraction = Math.abs(rubles % 1) > 0;
  const formatted = new Intl.NumberFormat('ru-RU', {
    minimumFractionDigits: hasFraction ? 2 : 0,
    maximumFractionDigits: hasFraction ? 2 : 0,
  }).format(rubles);

  return `${formatted} ₽`;
}

export function ProductCard({
  product,
  isFavorite,
  onToggleFavorite,
  compact = false,
  className = '',
}: ProductCardProps) {
  const price = formatPrice(product.priceCents);
  const rootRadiusClass = compact ? 'rounded-[20px]' : 'rounded-[22px]';
  const mediaAspectClass = compact ? 'aspect-[4/5]' : 'aspect-[11/12]';
  const mediaPaddingClass = compact ? 'p-2' : 'p-2.5';
  const imageClass = compact ? 'scale-[1.28] group-hover:scale-[1.36]' : '';
  const badgeOffsetClass = compact ? 'left-2 top-2' : 'left-2 top-2';
  const favoriteOffsetClass = compact ? 'right-2 top-2' : 'right-2 top-2';
  const favoriteSizeClass = compact ? 'h-7 w-7' : 'h-7 w-7';
  const favoriteIconClass = compact ? 'h-3.5 w-3.5' : 'h-3.5 w-3.5';
  const contentPaddingClass = compact ? 'p-2.5 pt-2' : 'p-2.5 pt-2';
  const titleSizeClass = compact ? 'text-[16px]' : 'text-[18px]';
  const priceSizeClass = compact ? 'text-[13px]' : 'text-[14px]';
  const descriptionClass = compact
    ? 'mb-1.5 text-[12px] font-medium leading-relaxed text-muted-foreground line-clamp-1'
    : 'mb-1.5 text-sm font-medium leading-relaxed text-muted-foreground line-clamp-2';
  const statsGapClass = compact ? 'gap-1.5' : 'gap-1.5';
  const statLabelClass = compact
    ? 'mb-1 text-[8px] font-bold uppercase tracking-[0.2em] text-muted-foreground'
    : 'mb-1.5 text-[9px] font-bold uppercase tracking-[0.2em] text-muted-foreground';
  const statValueClass = compact
    ? 'text-[11px] font-bold tracking-tight text-foreground'
    : 'text-[13px] font-bold tracking-tight text-foreground';

  return (
    <div
      className={`group relative overflow-hidden border border-border/40 bg-card shadow-sm transition-all duration-500 hover:border-foreground/30 hover:shadow-2xl dark:shadow-none ${rootRadiusClass} ${className}`}
    >
      <Link to={`/product/${product.id}`} className="block relative">
        <div
          className={`relative flex ${mediaAspectClass} items-center justify-center overflow-hidden bg-gradient-to-b from-foreground/5 to-transparent ${mediaPaddingClass}`}
        >
          <div className="absolute inset-0 bg-[radial-gradient(circle_at_50%_0%,var(--glow-primary),transparent_70%)] pointer-events-none mix-blend-overlay" />
          <ModelImage preset="catalog" src={product.photoUrl} alt={product.name} imageClassName={imageClass} />

          <div className={`absolute z-10 ${badgeOffsetClass}`}>
            <ProductStatusBadge status={product.status} tag={product.tag} />
          </div>
          <button
            type="button"
            onClick={(e) => {
              e.preventDefault();
              onToggleFavorite();
            }}
            className={`absolute z-10 flex items-center justify-center rounded-full backdrop-blur-2xl transition-all duration-300 ${favoriteOffsetClass} ${favoriteSizeClass} ${
              isFavorite
                ? 'bg-foreground border border-foreground text-background shadow-lg scale-110'
                : 'bg-background/50 border border-border/50 text-foreground hover:bg-background shadow-sm hover:scale-105'
            }`}
          >
            <Heart className={`${favoriteIconClass} ${isFavorite ? 'fill-current' : 'opacity-80'}`} />
          </button>
        </div>
      </Link>

      <div className={`relative z-20 border-t border-border/40 bg-card ${contentPaddingClass}`}>
        <Link to={`/product/${product.id}`}>
          <div className="mb-2 flex items-start justify-between gap-3">
            <h3 className={`${titleSizeClass} font-bold tracking-tight leading-tight text-foreground break-words`}>{product.name}</h3>
            {price && <span className={`shrink-0 ${priceSizeClass} font-semibold tracking-tight text-foreground`}>{price}</span>}
          </div>
          <p className={descriptionClass}>{product.description}</p>

          <div className={`flex items-center ${statsGapClass}`}>
            <div className="flex flex-col">
              <span className={statLabelClass}>Capacity</span>
              <span className={statValueClass}>{product.puffsMax.toLocaleString()} Puffs</span>
            </div>
            <div className="h-8 w-[1px] bg-border/50"></div>
            <div className="flex flex-col">
              <span className={statLabelClass}>Flavors</span>
              <span className={statValueClass}>{product.flavors.length} Options</span>
            </div>
          </div>
        </Link>
      </div>
    </div>
  );
}
