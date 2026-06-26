const ACCESS_TOKEN_KEY = 'bozena_media_access_token';
const USERNAME_KEY = 'bozena_media_username';

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
