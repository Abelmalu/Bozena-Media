const ACCESS_TOKEN_KEY = 'bozena_media_access_token';

export function readAccessToken() {
  return sessionStorage.getItem(ACCESS_TOKEN_KEY);
}

export function writeAccessToken(token: string) {
  sessionStorage.setItem(ACCESS_TOKEN_KEY, token);
}

export function clearAccessToken() {
  sessionStorage.removeItem(ACCESS_TOKEN_KEY);
}
