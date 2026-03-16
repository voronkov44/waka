import { useMemo } from 'react';
import { useParams, Link, useNavigate } from 'react-router';
import { ChevronLeft, Heart, Zap } from 'lucide-react';
import { useFavorites } from '../hooks/useFavorites';
import { FlavorChip } from '../components/flavor-chip';
import { ProductStatusBadge } from '../components/product-status-badge';
import { useCatalogModel } from '../hooks/useCatalogModel';
import { useCatalogModels } from '../hooks/useCatalogModels';

function formatPrice(cents: number | null) {
  if (cents === null || Number.isNaN(cents)) {
    return '—';
  }
  return `$${(cents / 100).toFixed(2)}`;
}

export function ProductDetail() {
  const { id } = useParams<{ id: string }>();
  const productID = Number(id);
  const navigate = useNavigate();
  const { isFavorite, toggleFavorite } = useFavorites();
  const { product, isLoading, error, notFound } = useCatalogModel(productID);
  const { products } = useCatalogModels();

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

  if (isLoading) {
    return (
      <div className="min-h-screen flex items-center justify-center text-muted-foreground">
        Loading product...
      </div>
    );
  }

  if (error) {
    return (
      <div className="min-h-screen flex items-center justify-center p-6">
        <div className="rounded-2xl border border-destructive/30 bg-destructive/10 p-4 text-sm text-destructive">
          {error}
        </div>
      </div>
    );
  }

  if (notFound || !product) {
    return (
      <div className="min-h-screen flex items-center justify-center">
        <div className="text-center">
          <p className="text-xl mb-4">Product not found</p>
          <Link to="/catalog" className="text-foreground underline">
            Back to catalog
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
          <div className="text-sm font-bold tracking-[0.2em] uppercase text-foreground">Details</div>
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
          <img src={product.photoUrl} alt={product.name} className="w-3/4 h-3/4 object-contain filter drop-shadow-2xl" />
          <div className="absolute top-6 left-6 z-10">
            <ProductStatusBadge status={product.status} tag={product.tag} />
          </div>
        </div>
      </div>

      <div className="px-6 pt-10">
        <h1 className="text-4xl font-extrabold tracking-tighter mb-4 leading-none">{product.name}</h1>
        <p className="text-[13px] font-medium leading-relaxed text-muted-foreground mb-10">{product.description}</p>

        <div className="bg-card border border-border/50 rounded-[32px] p-8 mb-8 shadow-sm">
          <div className="flex items-center gap-3 mb-6">
            <div className="w-10 h-10 rounded-[14px] bg-foreground/10 flex items-center justify-center">
              <Zap className="w-5 h-5 text-foreground" />
            </div>
            <h3 className="font-bold tracking-tight text-lg">Specifications</h3>
          </div>
          <div className="grid grid-cols-2 gap-6">
            <div>
              <p className="text-[10px] font-bold tracking-[0.2em] uppercase text-muted-foreground mb-2">Max Puffs</p>
              <p className="font-extrabold text-xl tracking-tight text-foreground">{product.puffsMax.toLocaleString()}</p>
            </div>
            <div>
              <p className="text-[10px] font-bold tracking-[0.2em] uppercase text-muted-foreground mb-2">Flavors</p>
              <p className="font-extrabold text-xl tracking-tight text-foreground">
                {product.flavors.length} <span className="text-sm font-medium text-muted-foreground">opts</span>
              </p>
            </div>
          </div>
        </div>

        <div className="mb-10">
          <h3 className="text-[11px] font-bold tracking-[0.2em] uppercase text-muted-foreground mb-4 pl-2">Available Flavors</h3>
          <div className="flex flex-wrap gap-2.5">
            {product.flavors.map((flavor) => (
              <FlavorChip key={flavor} flavor={flavor} />
            ))}
          </div>
        </div>

        <div className="bg-card border border-border/50 rounded-[32px] p-8 mb-10 shadow-sm relative overflow-hidden">
          <div className="absolute top-0 right-0 w-32 h-32 bg-foreground/5 rounded-full blur-2xl -translate-y-1/2 translate-x-1/2 pointer-events-none" />
          <div className="relative z-10">
            <p className="text-[10px] font-bold tracking-[0.2em] uppercase text-muted-foreground mb-2">Price</p>
            <p className="text-4xl font-extrabold tracking-tighter text-foreground">{formatPrice(product.priceCents)}</p>
          </div>
        </div>

        {relatedProducts.length > 0 && (
          <div>
            <div className="flex items-center justify-between mb-6 pl-2">
              <h3 className="text-[11px] font-bold tracking-[0.2em] uppercase text-muted-foreground">Related Models</h3>
            </div>
            <div className="flex gap-4 overflow-x-auto pb-6 -mx-6 px-6 scrollbar-hide snap-x">
              {relatedProducts.map((related) => (
                <Link
                  key={related.id}
                  to={`/product/${related.id}`}
                  className="flex-shrink-0 w-[240px] bg-card border border-border/50 rounded-[28px] overflow-hidden hover:border-foreground/30 hover:shadow-lg transition-all duration-500 snap-start shadow-sm"
                >
                  <div className="aspect-[4/3] bg-gradient-to-b from-foreground/5 to-transparent relative overflow-hidden flex items-center justify-center p-6">
                    <img src={related.photoUrl} alt={related.name} className="w-full h-full object-contain filter drop-shadow-xl" />
                  </div>
                  <div className="p-5 border-t border-border/40">
                    <p className="font-bold text-lg tracking-tight mb-1 truncate">{related.name}</p>
                    <p className="text-sm font-semibold text-foreground/80">{formatPrice(related.priceCents)}</p>
                  </div>
                </Link>
              ))}
            </div>
          </div>
        )}
      </div>
    </div>
  );
}
