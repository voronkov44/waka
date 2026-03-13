import { Heart, ChevronRight } from 'lucide-react';
import { Link } from 'react-router';
import type { Product } from '../types/domain';
import { ProductStatusBadge } from './product-status-badge';

interface ProductCardProps {
  product: Product;
  isFavorite: boolean;
  onToggleFavorite: () => void;
}

export function ProductCard({ product, isFavorite, onToggleFavorite }: ProductCardProps) {
  const formatPrice = (cents: number | null) => {
    if (cents === null || Number.isNaN(cents)) {
      return '—';
    }
    return `$${(cents / 100).toFixed(2)}`;
  };

  return (
    <div className="group relative bg-card rounded-[36px] overflow-hidden border border-border/40 hover:border-foreground/30 transition-all duration-500 shadow-sm hover:shadow-2xl dark:shadow-none">
      <Link to={`/product/${product.id}`} className="block relative">
        <div className="relative aspect-[4/5] sm:aspect-square overflow-hidden bg-gradient-to-b from-foreground/5 to-transparent flex items-center justify-center p-14">
          <div className="absolute inset-0 bg-[radial-gradient(circle_at_50%_0%,var(--glow-primary),transparent_70%)] pointer-events-none mix-blend-overlay" />
          
          <img
            src={product.photoUrl}
            alt={product.name}
            className="w-full h-full object-contain filter drop-shadow-2xl transition-transform duration-700 ease-[cubic-bezier(0.2,0.8,0.2,1)] group-hover:scale-110"
          />
          
          <div className="absolute top-6 left-6 z-10">
            <ProductStatusBadge status={product.status} tag={product.tag} />
          </div>
          <button
            type="button"
            onClick={(e) => {
              e.preventDefault();
              onToggleFavorite();
            }}
            className={`absolute top-6 right-6 w-11 h-11 rounded-full flex items-center justify-center backdrop-blur-2xl transition-all duration-300 z-10 ${
              isFavorite
                ? 'bg-foreground border border-foreground text-background shadow-lg scale-110'
                : 'bg-background/50 border border-border/50 text-foreground hover:bg-background shadow-sm hover:scale-105'
            }`}
          >
            <Heart className={`w-5 h-5 ${isFavorite ? 'fill-current' : 'opacity-80'}`} />
          </button>
        </div>
      </Link>

      <div className="relative p-8 pt-6 bg-card z-20 border-t border-border/40">
        <Link to={`/product/${product.id}`}>
          <div className="flex justify-between items-baseline mb-4">
            <h3 className="text-3xl font-bold tracking-tighter truncate pr-4 text-foreground">{product.name}</h3>
            <span className="text-xl font-semibold tracking-tight text-foreground shrink-0">
              {formatPrice(product.priceCents)}
            </span>
          </div>
          <p className="text-sm text-muted-foreground mb-8 line-clamp-2 leading-relaxed font-medium">{product.description}</p>

          <div className="flex items-center justify-between">
            <div className="flex items-center gap-6">
              <div className="flex flex-col">
                <span className="text-[9px] font-bold tracking-[0.2em] uppercase text-muted-foreground mb-1.5">Capacity</span>
                <span className="text-sm font-bold tracking-tight text-foreground">{product.puffsMax.toLocaleString()} Puffs</span>
              </div>
              <div className="w-[1px] h-8 bg-border/50"></div>
              <div className="flex flex-col">
                <span className="text-[9px] font-bold tracking-[0.2em] uppercase text-muted-foreground mb-1.5">Flavors</span>
                <span className="text-sm font-bold tracking-tight text-foreground">{product.flavors.length} Options</span>
              </div>
            </div>
            <div className="w-12 h-12 rounded-full bg-background border border-border/50 flex items-center justify-center group-hover:bg-foreground group-hover:text-background group-hover:border-foreground transition-all duration-500 shadow-sm">
              <ChevronRight className="w-5 h-5 text-current" />
            </div>
          </div>
        </Link>
      </div>
    </div>
  );
}
