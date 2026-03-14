# Waka Telegram Mini App

## Environment
Copy `.env.example` to `.env`.

- `VITE_API_BASE_URL`:
  - leave empty for recommended dev mode (`/api` same-origin requests + Vite proxy)
  - set to full HTTPS API URL only if you intentionally want direct browser -> API calls
- `VITE_API_PROXY_TARGET`:
  - local backend target for Vite proxy (default `http://localhost:28080`)
- `VITE_DEV_ALLOWED_HOSTS`:
  - comma-separated tunnel hosts allowed by Vite dev server
  - set this when using ngrok, for example `my-app.ngrok-free.app`

## Development run
1. Start backend API on `http://localhost:28080`.
2. Start frontend:
   - `npm i`
   - `npm run dev`
3. Expose frontend dev server with ngrok (or similar tunnel), for example:
   - `ngrok http 5173`
4. Add the generated ngrok host to `VITE_DEV_ALLOWED_HOSTS` and restart `npm run dev`.

## How API routing works now
- Frontend runtime calls `/api/...` by default.
- Vite dev server proxies `/api/...` to `VITE_API_PROXY_TARGET`.
- Result: Telegram/ngrok clients never call `http://localhost:28080` directly from the browser.

## Telegram auth flow
On app open inside Telegram:
1. Frontend waits briefly for `window.Telegram.WebApp`.
2. Calls `WebApp.ready()` and `WebApp.expand()` when available.
3. Resolves launch/auth data from:
   - `WebApp.initDataUnsafe.user`
   - `WebApp.initData`
   - `tgWebAppData` URL/query/hash fallback
4. Sends user payload to `POST /api/auth/telegram`.
5. Stores returned JWT.
6. Calls `GET /api/auth/me`.
7. Renders profile from backend user data.

Outside Telegram, app falls back to guest mode and keeps public screens usable.

## Dev diagnostics
- In development, Profile shows a compact "Telegram auth debug" panel.
- The latest bootstrap/auth diagnostics are also exposed as `window.__wakaTelegramBootstrapDebug`.

## Bot / Mini App URL for testing
Use the HTTPS ngrok frontend URL (the tunneled Vite URL), for example:
- `https://<your-ngrok-subdomain>.ngrok-free.app`
