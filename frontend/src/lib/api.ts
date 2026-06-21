import type {
  ApiEnvelope,
  AuthPayload,
  AuthResponse,
  FollowersResponse,
  FollowingsResponse,
  LikesResponse,
  PostDraft,
  PostsResponse,
  SearchUserResponse,
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

const API_BASE_URL = (import.meta.env.VITE_API_BASE_URL as string | undefined)?.replace(/\/$/, '') || 'http://localhost:8080';
const CLIENT_TYPE_HEADER = 'web';

let authHooks: AuthHooks | null = null;

export function setAuthHooks(hooks: AuthHooks) {
  authHooks = hooks;
}

function buildHeaders(init?: HeadersInit, includeJson = true) {
  const headers = new Headers(init);
  headers.set('X-Client-Type', CLIENT_TYPE_HEADER);
  headers.set('Accept', 'application/json');
  if (includeJson && !headers.has('Content-Type')) {
    headers.set('Content-Type', 'application/json');
  }
  return headers;
}

function extractErrorMessage(payload: unknown, fallback: string) {
  if (!payload || typeof payload !== 'object') {
    return fallback;
  }

  const maybeEnvelope = payload as ApiEnvelope<unknown> & { error?: { message?: string } };
  if (typeof maybeEnvelope.message === 'string' && maybeEnvelope.message.trim()) {
    return maybeEnvelope.message;
  }
  if (maybeEnvelope.error && typeof maybeEnvelope.error.message === 'string') {
    return maybeEnvelope.error.message;
  }
  return fallback;
}

async function parseResponse<T>(response: Response): Promise<T> {
  const text = await response.text();
  const payload = text ? (JSON.parse(text) as ApiEnvelope<T> | T) : null;

  if (payload && typeof payload === 'object' && 'success' in payload) {
    const envelope = payload as ApiEnvelope<T>;
    if (envelope.success === false) {
      throw new Error(extractErrorMessage(envelope, 'Request failed'));
    }
    return (envelope.data ?? (envelope as unknown as T)) as T;
  }

  if (payload === null) {
    return undefined as T;
  }

  return payload as T;
}

async function rawRequest<T>(path: string, options: RequestOptions = {}): Promise<T> {
  const { method = 'GET', body, auth = true } = options;
  const headers = buildHeaders();

  if (auth) {
    const token = authHooks?.getAccessToken() ?? null;
    if (token) {
      headers.set('Authorization', `Bearer ${token}`);
    }
  }

  const response = await fetch(`${API_BASE_URL}${path}`, {
    method,
    credentials: 'include',
    headers,
    body: body === undefined ? undefined : JSON.stringify(body),
  });

  if (!response.ok) {
    throw new Error(await readFailureMessage(response));
  }

  return parseResponse<T>(response);
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
      const retryResponse = await fetchWithAuth(path, options);
      if (retryResponse.ok) {
        return parseResponse<T>(retryResponse);
      }

      throw new Error(await readFailureMessage(retryResponse));
    }
  }

  throw new Error(await readFailureMessage(response));
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

async function readFailureMessage(response: Response) {
  const text = await response.text();
  if (!text) {
    return `Request failed with status ${response.status}`;
  }

  try {
    const payload = JSON.parse(text) as ApiEnvelope<unknown> | { message?: string; error?: { message?: string } };
    return extractErrorMessage(payload, `Request failed with status ${response.status}`);
  } catch {
    return text;
  }
}

export async function refreshSession() {
  return rawRequest<AuthResponse>('/api/auth/refresh', {
    method: 'POST',
    auth: false,
    retryOnUnauthorized: false,
  });
}

export function registerApiHooks(hooks: AuthHooks) {
  setAuthHooks(hooks);
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
  if (cursor) {
    params.set('cursor', cursor);
  }
  params.set('limit', String(limit));
  return request<SearchUserResponse>(`/api/auth/search?${params.toString()}`);
}

export async function listPosts() {
  return request<PostsResponse>('/api/posts/');
}

export async function createPost(payload: PostDraft) {
  return request<{ post_id: string; title: string; content: string }>('/api/posts/', {
    method: 'POST',
    body: payload,
  });
}

export async function updatePost(postId: number, payload: PostDraft) {
  return request<{ status: string; message: string }>(`/api/posts/update/${postId}`, {
    method: 'PUT',
    body: payload,
  });
}

export async function deletePost(postId: number) {
  return request<{ status: string; message: string }>(`/api/posts/delete/${postId}`, {
    method: 'DELETE',
  });
}

export async function getUserPosts(userId: number, cursor = '', limit = 10) {
  const params = new URLSearchParams();
  if (cursor) {
    params.set('cursor', cursor);
  }
  params.set('limit', String(limit));
  const query = params.toString();
  return request<{ posts: Array<{ id: number; title: string; content: string; user_id: number }>; cursor?: string; has_next?: boolean }>(
    `/api/posts/user/${userId}${query ? `?${query}` : ''}`,
  );
}

export async function toggleLike(postId: number, state: boolean) {
  return request<{ message: string }>(`/api/post/like/${postId}`, {
    method: 'POST',
    body: { like: state },
  });
}

export async function getPostLikes(postId: number, cursor = '', limit = 10) {
  const params = new URLSearchParams();
  if (cursor) {
    params.set('cursor', cursor);
  }
  params.set('limit', String(limit));
  const query = params.toString();
  return request<LikesResponse>(`/api/post/likes/${postId}${query ? `?${query}` : ''}`);
}

export async function toggleFollow(userId: number, follow: boolean) {
  return request<{ message: string }>(`/api/follow/${userId}`, {
    method: 'POST',
    body: { follow },
  });
}

export async function getFollowers(userId: number, cursor = '', limit = 10) {
  const params = new URLSearchParams();
  if (cursor) {
    params.set('cursor', cursor);
  }
  params.set('limit', String(limit));
  const query = params.toString();
  return request<FollowersResponse>(`/api/follow/followers/${userId}${query ? `?${query}` : ''}`);
}

export async function getFollowings(userId: number, cursor = '', limit = 10) {
  const params = new URLSearchParams();
  if (cursor) {
    params.set('cursor', cursor);
  }
  params.set('limit', String(limit));
  const query = params.toString();
  return request<FollowingsResponse>(`/api/follow/followings/${userId}${query ? `?${query}` : ''}`);
}
