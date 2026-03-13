const AUTH_TOKEN_KEY = 'waka_auth_token';

export function readAuthToken(): string | null {
  if (typeof window === 'undefined') {
    return null;
  }

  try {
    const token = localStorage.getItem(AUTH_TOKEN_KEY);
    return token && token.trim().length > 0 ? token : null;
  } catch {
    return null;
  }
}

export function writeAuthToken(token: string): void {
  if (typeof window === 'undefined') {
    return;
  }

  try {
    localStorage.setItem(AUTH_TOKEN_KEY, token);
  } catch {
    // Ignore persistence issues in restricted environments.
  }
}

export function clearAuthToken(): void {
  if (typeof window === 'undefined') {
    return;
  }

  try {
    localStorage.removeItem(AUTH_TOKEN_KEY);
  } catch {
    // Ignore persistence issues in restricted environments.
  }
}
