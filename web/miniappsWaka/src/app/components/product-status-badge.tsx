import type { ProductStatus, ProductTag } from '../types/domain';

interface ProductStatusBadgeProps {
  status: ProductStatus;
  tag?: ProductTag;
}

const statusClassNames: Record<ProductStatus, string> = {
  new: 'bg-foreground text-background border border-foreground/50',
  popular: 'bg-background/90 text-foreground border border-border/50',
  limited: 'bg-background/90 text-foreground border border-border/50',
};

function withAlpha(hexColor: string, alphaHex = '55') {
  if (!/^#[0-9A-Fa-f]{6}$/.test(hexColor)) {
    return undefined;
  }
  return `${hexColor}${alphaHex}`;
}

export function ProductStatusBadge({ status, tag }: ProductStatusBadgeProps) {
  if (tag) {
    return (
      <span
        className="px-4 py-1.5 rounded-full text-[9px] font-bold tracking-[0.2em] uppercase shadow-sm backdrop-blur-md"
        style={{
          backgroundColor: tag.bgColor,
          color: tag.textColor,
          border: `1px solid ${withAlpha(tag.textColor) ?? 'transparent'}`,
        }}
      >
        {tag.label}
      </span>
    );
  }

  return (
    <span
      className={`px-4 py-1.5 rounded-full text-[9px] font-bold tracking-[0.25em] uppercase shadow-sm backdrop-blur-md ${statusClassNames[status]}`}
    >
      {status}
    </span>
  );
}
