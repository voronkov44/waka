const DEFAULT_API_BASE_URL = 'http://localhost:28080';

function trimTrailingSlash(value: string): string {
  return value.replace(/\/+$/, '');
}

function resolveApiBaseURL(): string {
  const value = (import.meta.env.VITE_API_BASE_URL as string | undefined)?.trim();
  if (!value) {
    return DEFAULT_API_BASE_URL;
  }
  return trimTrailingSlash(value);
}

export const API_BASE_URL = resolveApiBaseURL();
