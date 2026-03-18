import type { ImgHTMLAttributes } from 'react';

export type ModelImagePreset = 'showcase' | 'featured' | 'catalog' | 'related' | 'details';

const presetClassMap: Record<ModelImagePreset, string> = {
  showcase:
    'h-full w-full object-contain object-center scale-[1.52] translate-y-[6%] filter drop-shadow-[0_34px_44px_rgba(0,0,0,0.36)]',
  featured:
    'h-full w-full object-contain object-center scale-[1.56] translate-y-[7%] filter drop-shadow-2xl transition-transform duration-700 ease-[cubic-bezier(0.2,0.8,0.2,1)] group-hover:scale-[1.66]',
  catalog:
    'h-full w-full object-contain object-center scale-[1.38] translate-y-[6%] filter drop-shadow-2xl transition-transform duration-700 ease-[cubic-bezier(0.2,0.8,0.2,1)] group-hover:scale-[1.46]',
  related:
    'h-full w-full object-contain object-center scale-[1.56] translate-y-[7%] filter drop-shadow-xl transition-transform duration-700 ease-[cubic-bezier(0.2,0.8,0.2,1)] group-hover:scale-[1.66]',
  details:
    'h-full w-full object-contain object-center scale-[1.34] translate-y-[5%] filter drop-shadow-[0_34px_46px_rgba(0,0,0,0.35)]',
};

interface ModelImageProps extends Omit<ImgHTMLAttributes<HTMLImageElement>, 'className'> {
  preset: ModelImagePreset;
  className?: string;
  imageClassName?: string;
}

export function ModelImage({ preset, className = '', imageClassName = '', ...imgProps }: ModelImageProps) {
  return (
    <div className={`h-full w-full overflow-hidden ${className}`}>
      <img {...imgProps} className={`${presetClassMap[preset]} ${imageClassName}`} />
    </div>
  );
}
