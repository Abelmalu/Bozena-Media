import type { SessionUser } from '../types';

function base64UrlDecode(value: string) {
  const normalized = value.replace(/-/g, '+').replace(/_/g, '/');
  const padding = '='.repeat((4 - (normalized.length % 4)) % 4);
  return atob(normalized + padding);
}

export function decodeSessionUser(token: string | null): SessionUser {
  if (!token) {
    return { userId: null, role: null };
  }

  const parts = token.split('.');
  if (parts.length < 2) {
    return { userId: null, role: null };
  }

  try {
    const payload = JSON.parse(base64UrlDecode(parts[1])) as Record<string, unknown>;
    const rawUserId = payload.user_id;
    const rawRole = payload.userRole;

    return {
      userId: typeof rawUserId === 'number' ? rawUserId : null,
      role: typeof rawRole === 'string' ? rawRole : null,
    };
  } catch {
    return { userId: null, role: null };
  }
}
