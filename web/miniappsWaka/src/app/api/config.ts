function trimTrailingSlash(value: string): string {
  return value.replace(/\/+$/, '');
}

function resolveApiBaseURL(): string {
  const value = (import.meta.env.VITE_API_BASE_URL as string | undefined)?.trim();
  if (!value) {
    // Default to same-origin API calls (e.g. /api/*) so Telegram/ngrok
    // clients do not attempt loopback calls to localhost.
    return '';
  }
  return trimTrailingSlash(value);
}

export const API_BASE_URL = resolveApiBaseURL();
