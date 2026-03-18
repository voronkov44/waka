import { ChevronRight, Heart, HelpCircle, User } from 'lucide-react';
import { Link } from 'react-router';
import { useFavorites } from '../hooks/useFavorites';
import { useCatalogModels } from '../hooks/useCatalogModels';
import { WakaFullLogo } from '../components/waka-brand';
import { VapeDeviceIcon } from '../components/icons/vape-device-icon';
import { ProductStatusBadge } from '../components/product-status-badge';

const quickActions = [
  { icon: VapeDeviceIcon, label: 'Browse Catalog', path: '/catalog' },
  { icon: Heart, label: 'My Favorites', path: '/favorites' },
  { icon: HelpCircle, label: 'Help & FAQ', path: '/faq' },
  { icon: User, label: 'My Profile', path: '/profile' },
];

export function Home() {
  const { favorites } = useFavorites();
  const { products, isLoading, error } = useCatalogModels();
  const heroProduct = products[0] ?? null;
  const featuredProducts = products.slice(0, 3);

  return (
    <div className="min-h-screen pb-24">
      <div className="pl-[26px] pr-[24px] pt-[32px] pb-[16px]">
        <div className="mb-2">
          <WakaFullLogo height={124} />
        </div>
      </div>

      <div className="px-6 mb-14 mt-4">
        <div className="relative overflow-hidden rounded-[40px] bg-foreground text-background p-10 shadow-[0_30px_60px_-15px_var(--glow-primary)]">
          <div className="absolute top-0 right-0 w-[400px] h-[400px] bg-background/10 rounded-full blur-3xl mix-blend-overlay -translate-y-1/2 translate-x-1/3" />
          <div className="absolute bottom-0 left-0 w-[300px] h-[300px] bg-background/5 rounded-full blur-3xl mix-blend-overlay translate-y-1/3 -translate-x-1/3" />

          <div className="relative flex items-center gap-8 z-10">
            <div className="flex-1">
              <div className="inline-flex items-center gap-2 border border-background/20 bg-background/10 backdrop-blur-md px-4 py-1.5 rounded-full text-[9px] font-bold mb-8 uppercase tracking-[0.3em] text-background">
                Featured
              </div>
              <h2 className="text-5xl font-extrabold mb-5 tracking-tighter leading-[0.9]">
                {heroProduct ? heroProduct.name : 'Waka'}
              </h2>
              <p className="text-background/70 mb-10 text-[13px] leading-relaxed max-w-[220px] font-medium tracking-wide">
                {heroProduct?.description ?? 'Explore our latest models and curated flavors.'}
              </p>
              <Link
                to={heroProduct ? `/product/${heroProduct.id}` : '/catalog'}
                className="inline-flex items-center justify-center text-foreground bg-background px-8 py-4 rounded-full font-bold text-[11px] tracking-[0.2em] uppercase hover:scale-105 transition-transform duration-500 shadow-xl"
              >
                Discover
              </Link>
            </div>
            <div className="w-40 h-56 flex-shrink-0 relative">
              {heroProduct && (
                <img
                  src={heroProduct.photoUrl}
                  alt={heroProduct.name}
                  className="w-full h-full object-contain scale-[1.8] origin-center drop-shadow-2xl"
                />
              )}
            </div>
          </div>
        </div>
      </div>

      <div className="px-6 mb-14">
        <h2 className="text-xl font-bold mb-6 tracking-tighter uppercase text-[11px] tracking-[0.2em] text-muted-foreground">
          Quick Actions
        </h2>
        <div className="grid grid-cols-2 gap-4">
          {quickActions.map((action) => {
            const Icon = action.icon;
            return (
              <Link
                key={action.path}
                to={action.path}
                className="group relative bg-card rounded-[32px] p-6 border border-border/40 hover:border-foreground/30 transition-all duration-500 shadow-sm hover:shadow-xl dark:shadow-none overflow-hidden"
              >
                <div className="absolute top-0 right-0 w-32 h-32 bg-gradient-to-bl from-foreground/5 to-transparent rounded-full blur-2xl -translate-y-1/2 translate-x-1/2 opacity-0 group-hover:opacity-100 transition-opacity duration-500" />
                <div className="w-14 h-14 bg-background rounded-2xl flex items-center justify-center mb-6 shadow-sm border border-border/50 group-hover:scale-110 transition-transform duration-500 ease-[cubic-bezier(0.2,0.8,0.2,1)]">
                  <Icon className="w-6 h-6 text-foreground" />
                </div>
                <p className="font-bold tracking-tight text-lg relative z-10 text-foreground">{action.label}</p>
              </Link>
            );
          })}
        </div>
      </div>

      <div className="px-6 mb-10">
        <div className="flex items-center justify-between mb-8">
          <h2 className="text-xl font-bold tracking-tighter uppercase text-[11px] tracking-[0.2em] text-muted-foreground">
            Featured Models
          </h2>
          <Link
            to="/catalog"
            className="text-foreground text-[10px] font-bold uppercase tracking-[0.2em] flex items-center gap-1 transition-all hover:opacity-70"
          >
            See all
            <ChevronRight className="w-3.5 h-3.5" />
          </Link>
        </div>

        {isLoading && <p className="text-sm text-muted-foreground">Loading models...</p>}

        {error && (
          <div className="rounded-2xl border border-destructive/30 bg-destructive/10 p-4 text-sm text-destructive">
            {error}
          </div>
        )}

        {!isLoading && !error && featuredProducts.length > 0 && (
          <div className="flex gap-5 overflow-x-auto pb-10 -mx-6 px-6 scrollbar-hide snap-x">
            {featuredProducts.map((product) => {
              const flavorCount = Array.isArray(product.flavors) ? product.flavors.length : 0;
              const description = typeof product.description === 'string' ? product.description.trim() : '';

              return (
                <Link
                  key={product.id}
                  to={`/product/${product.id}`}
                  className="group relative flex-shrink-0 w-[300px] bg-card rounded-[36px] overflow-hidden border border-border/40 hover:border-foreground/30 transition-all duration-500 shadow-sm hover:shadow-2xl dark:shadow-none snap-start"
                >
                  <div className="aspect-[4/3] bg-gradient-to-b from-foreground/5 to-transparent relative overflow-hidden flex items-center justify-center p-10">
                    <div className="absolute inset-0 bg-[radial-gradient(circle_at_50%_0%,var(--glow-primary),transparent_70%)] pointer-events-none mix-blend-overlay" />
                    <img
                      src={product.photoUrl}
                      alt={product.name}
                      className="w-full h-full object-contain filter drop-shadow-2xl transition-transform duration-700 ease-[cubic-bezier(0.2,0.8,0.2,1)] group-hover:scale-110"
                    />
                    <div className="absolute top-6 left-6 z-10">
                      <ProductStatusBadge status={product.status} tag={product.tag} />
                    </div>
                  </div>
                  <div className="relative p-7 pt-5 bg-card z-20 border-t border-border/40">
                    <h3 className="text-2xl font-bold tracking-tighter leading-tight text-foreground break-words">{product.name}</h3>
                    {description.length > 0 && (
                      <p className="mb-3 mt-2 text-sm font-medium leading-relaxed text-muted-foreground line-clamp-2">
                        {description}
                      </p>
                    )}

                    <div className="mt-3 flex items-center gap-1.5">
                      <div className="flex flex-col">
                        <span className="mb-1.5 text-[9px] font-bold uppercase tracking-[0.2em] text-muted-foreground">
                          Capacity
                        </span>
                        <span className="text-[13px] font-bold tracking-tight text-foreground">
                          {product.puffsMax.toLocaleString()} Puffs
                        </span>
                      </div>
                      <div className="h-8 w-[1px] bg-border/50"></div>
                      <div className="flex flex-col">
                        <span className="mb-1.5 text-[9px] font-bold uppercase tracking-[0.2em] text-muted-foreground">
                          Flavors
                        </span>
                        <span className="text-[13px] font-bold tracking-tight text-foreground">{flavorCount} Options</span>
                      </div>
                    </div>
                  </div>
                </Link>
              );
            })}
          </div>
        )}
      </div>

      {favorites.length > 0 && (
        <div className="px-6 mt-12">
          <div className="relative overflow-hidden bg-card rounded-[32px] p-8 border border-border/40 shadow-sm dark:shadow-none">
            <div className="absolute top-0 right-0 w-48 h-48 bg-foreground/5 rounded-full blur-2xl -translate-y-1/2 translate-x-1/4 pointer-events-none" />
            <div className="flex items-center justify-between relative z-10">
              <div>
                <p className="text-[10px] font-bold tracking-[0.2em] uppercase text-muted-foreground mb-2">
                  Your Favorites
                </p>
                <p className="text-4xl font-bold tracking-tighter text-foreground">{favorites.length}</p>
              </div>
              <div className="w-16 h-16 rounded-full bg-background shadow-sm border border-border/50 flex items-center justify-center">
                <Heart className="w-7 h-7 text-foreground fill-current drop-shadow-sm" />
              </div>
            </div>
          </div>
        </div>
      )}

      <div className="px-6 mt-10 mb-4">
        <div className="flex justify-center opacity-30">
          <WakaFullLogo height={72} />
        </div>
      </div>
    </div>
  );
}
