import { LucideIcon } from 'lucide-react';

interface EmptyStateProps {
  icon: LucideIcon;
  title: string;
  description: string;
  action?: {
    label: string;
    onClick: () => void;
  };
}

export function EmptyState({ icon: Icon, title, description, action }: EmptyStateProps) {
  return (
    <div className="flex flex-col items-center justify-center text-center px-6 py-20 bg-card border border-border/50 rounded-[36px] shadow-sm dark:shadow-none relative overflow-hidden">
      <div className="absolute top-0 right-0 w-48 h-48 bg-foreground/5 rounded-full blur-3xl -translate-y-1/2 translate-x-1/4 pointer-events-none" />
      <div className="w-24 h-24 rounded-full bg-background border border-border/50 shadow-sm flex items-center justify-center mb-8 relative z-10">
        <Icon className="w-10 h-10 text-foreground opacity-80" />
      </div>
      <h3 className="text-2xl font-bold tracking-tight mb-3 text-foreground relative z-10">{title}</h3>
      <p className="text-[13px] font-medium leading-relaxed text-muted-foreground mb-8 max-w-[240px] relative z-10">{description}</p>
      {action && (
        <button
          onClick={action.onClick}
          className="relative z-10 px-8 py-4 bg-foreground text-background rounded-full text-[11px] font-bold tracking-[0.2em] uppercase hover:scale-105 transition-all duration-500 shadow-md hover:shadow-lg"
        >
          {action.label}
        </button>
      )}
    </div>
  );
}
