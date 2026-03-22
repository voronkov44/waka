import { brandAssets } from '../assets';
import { useI18n } from '../../shared/i18n';

interface WakaIconProps {
  className?: string;
  size?: number;
}

interface WakaFullLogoProps {
  className?: string;
  height?: number;
}

const THEME_ADAPTIVE_CLASS = 'brightness-0 dark:brightness-100';

export function WakaIcon({ className = '', size = 48 }: WakaIconProps) {
  const { t } = useI18n();

  return (
    <div
      className={`inline-flex items-center overflow-hidden ${className}`}
      style={{ width: size, height: size }}
    >
      <img
        src={brandAssets.logoIcon}
        alt={t('media.wakaAlt')}
        className={`h-full w-full object-contain ${THEME_ADAPTIVE_CLASS}`}
      />
    </div>
  );
}

export function WakaFullLogo({ className = '', height = 48 }: WakaFullLogoProps) {
  const { t } = useI18n();

  return (
    <img
      src={brandAssets.logoBig}
      alt={t('media.wakaAlt')}
      style={{ height }}
      className={`h-auto w-auto max-w-full object-contain object-left ${THEME_ADAPTIVE_CLASS} ${className}`}
    />
  );
}
