import type { ProductStatus, ProductTag } from '../types/domain';

interface ProductStatusBadgeProps {
  status: ProductStatus;
  tag?: ProductTag;
}

function withAlpha(hexColor: string, alphaHex = '55') {
  if (!/^#[0-9A-Fa-f]{6}$/.test(hexColor)) {
    return undefined;
  }
  return `${hexColor}${alphaHex}`;
}

export function ProductStatusBadge({ tag }: ProductStatusBadgeProps) {
  const label = typeof tag?.label === 'string' ? tag.label.trim() : '';
  if (label.length === 0) {
    return null;
  }

  return (
    <span
      className="px-4 py-1.5 rounded-full text-[9px] font-bold tracking-[0.2em] uppercase shadow-sm backdrop-blur-md"
      style={{
        backgroundColor: tag?.bgColor,
        color: tag?.textColor,
        border: `1px solid ${withAlpha(tag?.textColor ?? '') ?? 'transparent'}`,
      }}
    >
      {label}
    </span>
  );
}
