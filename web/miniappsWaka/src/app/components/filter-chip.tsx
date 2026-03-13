interface FilterChipProps {
  label: string;
  active: boolean;
  onClick: () => void;
}

export function FilterChip({ label, active, onClick }: FilterChipProps) {
  return (
    <button
      type="button"
      onClick={onClick}
      className={`px-6 py-3 rounded-full text-[10px] font-bold tracking-[0.2em] uppercase whitespace-nowrap transition-all duration-500 border ${
        active
          ? 'bg-foreground border-foreground text-background shadow-md scale-105'
          : 'bg-card border-border/50 text-muted-foreground hover:border-foreground/50 hover:text-foreground hover:scale-105'
      }`}
    >
      {label}
    </button>
  );
}
