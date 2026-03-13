interface VapeDeviceIconProps {
  className?: string;
}

export function VapeDeviceIcon({ className }: VapeDeviceIconProps) {
  return (
    <svg
      className={className}
      viewBox="0 0 24 24"
      fill="none"
      stroke="currentColor"
      strokeWidth="1.8"
      strokeLinecap="round"
      strokeLinejoin="round"
      aria-hidden="true"
    >
      <rect x="7" y="3" width="10" height="18" rx="3" />
      <line x1="12" y1="3" x2="12" y2="1" />
      <circle cx="12" cy="17" r="1" fill="currentColor" />
      <line x1="10" y1="8" x2="14" y2="8" />
    </svg>
  );
}
