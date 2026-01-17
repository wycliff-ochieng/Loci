# Architecture and Decisions

## High-level shape
- **Go backend** with Gorilla mux for routing, sqlc for typed DB access, and Postgres/PostGIS storage. Services orchestrate business logic; handlers stay thin.
- **Next.js frontend (App Router)** for the map-first UX, using client components where browser APIs are needed (geolocation, Leaflet).
- **Redis rate limiter** (per-user) to throttle locus posting.
- **WebSocket hub** scaffolded for real-time fan-out of loci, replies, and views.

## Key decisions
- **Separation of concerns:**
  - `internal/transport/http`: request parsing/validation only.
  - `internal/service`: business rules, transactions, and hub broadcasts.
  - `internal/store/queries` + `sqlc/`: SQL lives in `.sql` files; generated Go code is type-safe.
- **PostGIS-native storage:** `location` is stored as `geography`; queries use `ST_MakeEnvelope` (bounds) and `ST_MakePoint` (insert) to stay in WGS84.
- **JWT auth middleware:** bearer tokens guard write paths; context carries the user UUID to services.
- **Rate limiting:** Redis sliding-window limiter per user for posting loci to prevent spam.
- **Event hub:** `internal/socket.Hub` holds broadcast channels (`BroadcastLocus`, `BroadcastReply`, `BroadcastView`). Handlers/services push events; a WS endpoint can fan-out to clients.

## Design patterns
- **Thin controllers, fat services:** handlers validate/parse; services perform transactions and domain logic (reply count increment, view counting, broadcasts).
- **Transactional writes:** replies and views run inside pgx transactions to keep counts consistent.
- **Data mappers:** sqlc rows are mapped into `models` before leaving the service layer; models carry JSON tags so the API payloads are stable.
- **Idempotent view recording:** view creation + count increment happens in one transaction; hub broadcast is best-effort after commit.
- **Client-side resilience:** frontend fetch helpers normalize varying field shapes and avoid throwing on non-200 where user experience should continue (e.g., loci fetch returns `[]` on 401).

## Current surface area
- **Endpoints** (see README for table): register/login, list loci (public), post locus (auth), list replies (public), post reply (auth), record view (auth), websocket endpoint (auth).
- **Frontend routes:** `/` live map + feed + drop-locus composer; `/loci/[id]` thread view; `/login`, `/signup` auth forms.

## Notable constraints
- **Auth required for writes:** posting loci, replies, and view events require a bearer token.
- **Map rendering:** Leaflet needs client components; SSR falls back to client-only fetches to avoid `window` errors.
- **Legacy data:** Some rows may lack lat/long JSON fields; frontend mapping tolerates both uppercase/lowercase and nested location objects.
