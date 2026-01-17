## LOCI

Geo-messaging platform for dropping short messages (“loci”) pinned to map coordinates, with real-time visibility, replies, and view tracking.

### Why Things Are Structured This Way
- **Map-first experience:** App router (Next.js) renders a live map + feed so users can see density and open threads quickly.
- **PostGIS-ready backend:** `internal/store` targets a Postgres/PostGIS database for geo queries (bounding boxes) and persistence.
- **JWT auth & rate limiting:** Auth middleware gates protected routes; Redis-backed limiter (`internal/limitter`) prevents locus spam per user.
- **Websocket hub ready:** `internal/socket` provides a hub for future real-time fan-out of new loci/replies without reloading.
- **sqlc for safety:** Queries in `internal/store/queries` generate typed data access (`sqlc/`), reducing runtime SQL errors.
- **Separation of concerns:** Handlers (`internal/transport/http`) stay thin; services (`internal/service`) orchestrate business rules; storage isolated in `internal/store`.

### Endpoints (current)
| Method | Path | Auth | Purpose |
| --- | --- | --- | --- |
| POST | `/register` | none | Create a user account. |
| POST | `/login` | none | Login with email/username + password, returns JWT pair. |
| GET | `/api/get/loci/` | none | Fetch loci within a bounding box (`SouthWestLat`, `SouthWestLong`, `NorthEastLat`, `NorthEastLong`). |
| POST | `/api/post/loci` | JWT | Create a new locus message (body includes message + coords). |
| GET | `/api/loci/{id}/replies` | none | List replies for a locus (newest first). |
| POST | `/api/loci/{id}/reply` | JWT | Create a reply for a locus. |
| POST | `/api/loci/{id}/view` | JWT | Record a view for a locus (increments view_count). |
| GET | `/ws` | JWT | WebSocket endpoint (hub broadcasts for loci/replies/views).

### Configuration
Set environment variables (or `.env` loaded by `internal/config`):
- `DB_HOST`, `DB_PORT`, `DB_NAME`, `DB_USER`, `DB_PASSWORD`, `DB_SSLMODE`
- `JWTSecret` (signing secret)
- `REDIS_ADDR`, `REDIS_PASSWORD`

Defaults (dev): Postgres on localhost:5432, Redis on localhost:6379, JWT secret `mydogiscalledrufus`.

### Getting Started (Backend)
1. Ensure Postgres (ideally PostGIS) and Redis are running.
2. Export env vars or create `.env` with the keys above.
3. Run database migrations from `internal/store/migrations/` (e.g., `psql -f 20251003152017_create_user_table.sql ...`).
4. Generate sqlc code if queries change: `sqlc generate` (already vendored in `sqlc/`).
5. Start the API: `go run cmd/main.go` (listens on `:3000`).

### Getting Started (Frontend)
- Location: `web/`
- Install: `cd web && npm install`
- Run dev server: `npm run dev` (defaults to http://localhost:3000)
- Build: `npm run build`
- Start production build: `npm run start`
- Configure API base: set `NEXT_PUBLIC_API_BASE` to your backend (e.g., `http://localhost:3000`).

App routes: live map + feed at `/`, thread view at `/loci/[id]`, auth pages at `/login` and `/signup`.

### Project Layout (key parts)
- `cmd/main.go` — bootstraps server with config, Redis limiter, router.
- `api/` — server wiring and route registration (mux + CORS).
- `internal/transport/http` — HTTP handlers.
- `internal/service` — business logic (auth, loci creation, geo queries, replies, views, broadcasts).
- `internal/store` — DB setup, migrations, SQL queries.
- `pkg/middleware` — auth middleware.
- `internal/socket` — websocket hub for broadcast channels.
- `web/` — Next.js app router frontend.
- `docs/` — detailed architecture, geolocation/PostGIS, replies/views, and frontend flow notes.

### Notes & Next Steps
- Frontend needs a valid bearer token to post loci/replies/views; public reads work unauthenticated.
- Configure `NEXT_PUBLIC_API_BASE` to point at the Go API (defaults to `http://localhost:3000`).
- Consider enabling HTTPS/production configs and containerizing DB/Redis for deployment.
- See `docs/` for deep dives on architecture, geolocation/PostGIS, replies/views, and frontend flows.