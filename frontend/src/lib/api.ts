import type {
  ApiEnvelope,
  AuthPayload,
  AuthResponse,
  FeedResponse,
  FollowersResponse,
  FollowingsResponse,
  LikesResponse,
  PostDraft,
  PostResponse,
  SearchUserResponse,
  UserPostsResponse,
} from '../types';

type AuthHooks = {
  getAccessToken: () => string | null;
  setAccessToken: (token: string) => void;
  clearAccessToken: () => void;
  refreshSession: () => Promise<boolean>;
};

type RequestOptions = {
  method?: string;
  body?: unknown;
  auth?: boolean;
  retryOnUnauthorized?: boolean;
};

const API_BASE_URL = (import.meta.env.VITE_API_BASE_URL as string | undefined)?.replace(/\/$/, '') || 'http://localhost:8082';
const CLIENT_TYPE_HEADER = 'web';

let authHooks: AuthHooks | null = null;

export function registerAuthHooks(hooks: AuthHooks) {
  authHooks = hooks;
}

function buildHeaders(init?: HeadersInit) {
  const headers = new Headers(init);
  headers.set('X-Client-Type', CLIENT_TYPE_HEADER);
  headers.set('Accept', 'application/json');
  return headers;
}

function extractErrorMessage(payload: unknown, fallback: string) {
  if (!payload || typeof payload !== 'object') {
    return fallback;
  }

  const envelope = payload as ApiEnvelope<unknown>;
  if (typeof envelope.message === 'string' && envelope.message.trim()) {
    return envelope.message;
  }
  if (envelope.error?.message) {
    return envelope.error.message;
  }
  return fallback;
}

async function readJson(response: Response) {
  const text = await response.text();
  return text ? JSON.parse(text) : null;
}

async function readFailureMessage(response: Response) {
  const payload = await readJson(response).catch(() => null);
  return extractErrorMessage(payload, `Request failed with status ${response.status}`);
}

async function parseResponse<T>(response: Response): Promise<T> {
  const payload = await readJson(response);

  if (payload && typeof payload === 'object' && 'success' in payload) {
    const envelope = payload as ApiEnvelope<T>;
    if (envelope.success === false) {
      throw new Error(extractErrorMessage(envelope, 'Request failed'));
    }
    return (envelope.data ?? (envelope as unknown as T)) as T;
  }

  return payload as T;
}

async function fetchWithAuth(path: string, options: RequestOptions = {}) {
  const { method = 'GET', body, auth = true } = options;
  const headers = buildHeaders();

  if (auth) {
    const token = authHooks?.getAccessToken() ?? null;
    if (token) {
      headers.set('Authorization', `Bearer ${token}`);
    }
  }

  return fetch(`${API_BASE_URL}${path}`, {
    method,
    credentials: 'include',
    headers,
    body: body === undefined ? undefined : JSON.stringify(body),
  });
}

export async function request<T>(path: string, options: RequestOptions = {}): Promise<T> {
  const { auth = true, retryOnUnauthorized = true } = options;
  const response = await fetchWithAuth(path, options);

  if (response.ok) {
    return parseResponse<T>(response);
  }

  if (response.status === 401 && auth && retryOnUnauthorized && authHooks) {
    const refreshed = await authHooks.refreshSession();
    if (refreshed) {
      const retry = await fetchWithAuth(path, options);
      if (retry.ok) {
        return parseResponse<T>(retry);
      }
      throw new Error(await readFailureMessage(retry));
    }
  }

  throw new Error(await readFailureMessage(response));
}

export async function refreshSession() {
  const response = await fetch(`${API_BASE_URL}/api/auth/refresh`, {
    method: 'POST',
    credentials: 'include',
    headers: buildHeaders(),
  });

  if (!response.ok) {
    throw new Error(await readFailureMessage(response));
  }

  return (await readJson(response)) as AuthResponse;
}

export async function login(payload: { username: string; password: string }) {
  return request<AuthPayload>('/api/auth/login', {
    method: 'POST',
    auth: false,
    body: payload,
  });
}

export async function registerUser(payload: {
  name: string;
  username: string;
  email: string;
  password: string;
}) {
  return request<AuthPayload>('/api/auth/register', {
    method: 'POST',
    auth: false,
    body: payload,
  });
}

export async function logout() {
  return request<{ message: string }>('/api/auth/logout', {
    method: 'POST',
  });
}

export async function searchUsers(search: string, cursor = '', limit = 10) {
  const params = new URLSearchParams();
  params.set('search', search);
  params.set('limit', String(limit));
  if (cursor) {
    params.set('cursor', cursor);
  }
  return request<SearchUserResponse>(`/api/auth/search?${params.toString()}`);
}

export async function getFeed(cursor = '', limit = 10) {
  const params = new URLSearchParams();
  params.set('limit', String(limit));
  if (cursor) {
    params.set('cursor', cursor);
  }
  return request<FeedResponse>(`/api/feed/?${params.toString()}`);
}

export async function createPost(payload: PostDraft) {
  return request<PostResponse>('/api/posts/', {
    method: 'POST',
    body: payload,
  });
}

export async function toggleLike(postId: number, state: boolean) {
  return request<{ message: string }>(`/api/post/like/${postId}`, {
    method: 'POST',
    body: { like: state },
  });
}

export async function getPostLikes(postId: number, cursor = '', limit = 10) {
  const params = new URLSearchParams();
  params.set('limit', String(limit));
  if (cursor) {
    params.set('cursor', cursor);
  }
  return request<LikesResponse>(`/api/post/likes/${postId}?${params.toString()}`);
}

export async function getUserFollowers(userId: number, cursor = '', limit = 10) {
  const params = new URLSearchParams();
  params.set('limit', String(limit));
  if (cursor) {
    params.set('cursor', cursor);
  }
  return request<FollowersResponse>(`/api/follow/followers/${userId}?${params.toString()}`);
}

export async function getUserFollowings(userId: number, cursor = '', limit = 10) {
  const params = new URLSearchParams();
  params.set('limit', String(limit));
  if (cursor) {
    params.set('cursor', cursor);
  }
  return request<FollowingsResponse>(`/api/follow/followings/${userId}?${params.toString()}`);
}

export async function getUserPosts(userId: number, cursor = '', limit = 10) {
  const params = new URLSearchParams();
  params.set('limit', String(limit));
  if (cursor) {
    params.set('cursor', cursor);
  }
  return request<UserPostsResponse>(`/api/posts/user/${userId}?${params.toString()}`);
}

export async function toggleFollow(userId: number) {
  return request<{ message: string }>(`/api/follow/${userId}`, {
    method: 'POST',
  });
}
