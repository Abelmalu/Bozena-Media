# Frontend

React + TypeScript frontend for the Bozena Media gateway.

## Run locally

1. Install dependencies.
2. Set `VITE_API_BASE_URL` if the API Gateway is not on `http://localhost:8080`.
3. Start the dev server with `npm run dev`.

## Notes

- The app uses `X-Client-Type: web` on every request.
- The gateway currently issues the refresh token as an HttpOnly cookie for web clients.
- The refresh cookie is marked `Secure`, so local HTTP dev may require HTTPS or a browser that accepts the cookie in your setup.
- Search results do not include user IDs, so profile navigation is only possible where the backend returns an ID.
