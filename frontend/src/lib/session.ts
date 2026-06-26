const ACCESS_TOKEN_KEY = 'bozena_media_access_token';
const USERNAME_KEY = 'bozena_media_username';
const LIKED_POSTS_KEY = 'bozena_media_liked_posts';

export function readAccessToken() {
  return sessionStorage.getItem(ACCESS_TOKEN_KEY);
}

export function writeAccessToken(token: string) {
  sessionStorage.setItem(ACCESS_TOKEN_KEY, token);
}

export function clearAccessToken() {
  sessionStorage.removeItem(ACCESS_TOKEN_KEY);
}

export function readUsername() {
  return sessionStorage.getItem(USERNAME_KEY);
}

export function writeUsername(username: string) {
  sessionStorage.setItem(USERNAME_KEY, username);
}

export function clearUsername() {
  sessionStorage.removeItem(USERNAME_KEY);
}

export function readLikedPosts() {
  const raw = sessionStorage.getItem(LIKED_POSTS_KEY);
  if (!raw) {
    return {} as Record<number, boolean>;
  }

  try {
    const parsed = JSON.parse(raw) as number[];
    return parsed.reduce<Record<number, boolean>>((acc, postId) => {
      acc[postId] = true;
      return acc;
    }, {});
  } catch {
    return {} as Record<number, boolean>;
  }
}

export function writeLikedPosts(likedMap: Record<number, boolean>) {
  const likedPostIds = Object.entries(likedMap)
    .filter(([, liked]) => liked)
    .map(([postId]) => Number(postId))
    .filter((postId) => Number.isFinite(postId));

  sessionStorage.setItem(LIKED_POSTS_KEY, JSON.stringify(likedPostIds));
}

export function clearLikedPosts() {
  sessionStorage.removeItem(LIKED_POSTS_KEY);
}
