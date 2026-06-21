export type ApiErrorPayload = {
  code?: string;
  message?: string;
  details?: unknown;
};

export type ApiEnvelope<T> = {
  success?: boolean;
  data?: T;
  request_id?: string;
  timestamp?: string;
  error?: ApiErrorPayload;
  message?: string;
};

export type AuthResponse = {
  access_token: string;
  refresh_token?: string;
  message?: string;
};

export type AuthPayload = {
  access_token: string;
  refresh_token?: string;
};

export type SearchUserItem = {
  name: string;
  username: string;
};

export type SearchUserResponse = {
  users: SearchUserItem[];
  cursor?: string;
  has_next?: boolean;
};

export type Post = {
  id: number;
  title: string;
  content: string;
  user_id: number;
};

export type PostsResponse = {
  posts: Post[];
};

export type CursorListResponse<T> = {
  cursor?: string;
  has_next?: boolean;
} & Record<string, T[]>;

export type FollowUser = {
  name: string;
  username: string;
};

export type FollowersResponse = {
  followers: FollowUser[];
  cursor?: string;
  has_next?: boolean;
  limit?: number;
};

export type FollowingsResponse = {
  followings: FollowUser[];
  cursor?: string;
  has_next?: boolean;
  limit?: number;
};

export type LikesResponse = {
  users: FollowUser[];
  cursor?: string;
  has_next?: boolean;
};

export type PostDraft = {
  title: string;
  content: string;
};

export type SessionUser = {
  userId: number | null;
  role: string | null;
};
