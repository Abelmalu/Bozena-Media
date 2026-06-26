export type ApiEnvelope<T> = {
  success?: boolean;
  data?: T;
  request_id?: string;
  timestamp?: string;
  error?: {
    code?: string;
    message?: string;
    details?: unknown;
  };
  message?: string;
};

export type AuthPayload = {
  access_token: string;
  refresh_token?: string;
  username?: string;
  message?: string;
};

export type AuthResponse = AuthPayload;

export type SearchUser = {
  username: string;
  name: string;
};

export type SearchUserResponse = {
  users: SearchUser[];
  cursor?: string;
  has_next?: boolean;
};

export type FeedItem = {
  PostID: number;
  PostTitle: string;
  PostContent: string;
  PostOwnerID: number;
  UserName: string;
  Name: string;
  id?: number;
  post_id?: number;
  postId?: number;
  post_owner_id?: number;
  postOwnerId?: number;
  user_id?: number;
  userId?: number;
};

export type FeedResponse = {
  userfeeds: FeedItem[];
  cursor?: string;
  limit?: number;
};

export type ProfileUser = {
  name: string;
  username: string;
};

export type FollowersResponse = {
  followers: ProfileUser[];
  has_next?: boolean;
  limit?: number;
  cursor?: string;
};

export type FollowingsResponse = {
  Followings: ProfileUser[];
  has_next?: boolean;
  limit?: number;
  cursor?: string;
};

export type PostDraft = {
  title: string;
  content: string;
};

export type PostResponse = {
  post_id?: string;
  title?: string;
  content?: string;
  status?: string;
  message?: string;
};

export type UserPost = {
  id: number;
  title: string;
  content: string;
  user_id: number;
};

export type UserPostsResponse = {
  posts: UserPost[];
  cursor?: string;
  has_next?: boolean;
};

export type LikesResponse = {
  users: ProfileUser[];
  cursor?: string;
  has_next?: boolean;
};

export type SessionUser = {
  userId: number | null;
  role: string | null;
};
