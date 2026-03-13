interface FlavorChipProps {
  flavor: string;
}

export function FlavorChip({ flavor }: FlavorChipProps) {
  return (
    <span className="inline-flex items-center px-4 py-2 bg-card border border-border/50 rounded-full text-xs font-bold tracking-wide text-foreground shadow-sm">
      {flavor}
    </span>
  );
}
