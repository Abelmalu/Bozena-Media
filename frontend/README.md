# Frontend

React + TypeScript frontend for the Bozena Media gateway.

## Run locally

1. `cd frontend`
2. `npm install`
3. `npm run dev`

The app expects the API Gateway at `http://localhost:8082` by default.
Because the gateway sets the refresh cookie as `Secure`, browser refresh/logout behavior may require HTTPS locally depending on your setup.

## What it shows

- Feed first after login, loaded from `/api/feed/`
- Sidebar navigation for feed, profile, and search
- Profile counts based on the backend-returned follower/following page results
- User search using the current auth search endpoint
