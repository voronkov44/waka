import { ChevronLeft } from 'lucide-react';
import { Link } from 'react-router';

export function TermsOfService() {
  return (
    <div className="min-h-screen pb-32">
      <div className="sticky top-0 z-40 border-b border-border/40 bg-background/80 pt-safe backdrop-blur-[32px]">
        <div className="flex items-center justify-between px-6 py-4">
          <Link
            to="/profile"
            className="flex h-12 w-12 items-center justify-center rounded-[18px] border border-border/50 bg-card shadow-sm transition-all duration-300 hover:scale-105"
          >
            <ChevronLeft className="h-6 w-6 text-foreground" />
          </Link>
          <div className="text-sm font-bold uppercase tracking-[0.2em] text-foreground">Legal</div>
          <div className="h-12 w-12" />
        </div>
      </div>

      <div className="px-6 pt-8">
        <div className="rounded-[32px] border border-border/50 bg-card p-8 shadow-sm">
          <h1 className="mb-4 text-3xl font-extrabold tracking-tighter text-foreground">Terms of Service</h1>
          <p className="mb-6 text-sm font-medium leading-relaxed text-muted-foreground">
            This section will contain the usage terms for the Waka mini app.
          </p>
          <div className="rounded-2xl border border-border/50 bg-background p-5">
            <p className="text-[11px] font-bold uppercase tracking-[0.14em] text-muted-foreground">Content coming soon</p>
          </div>
        </div>
      </div>
    </div>
  );
}
