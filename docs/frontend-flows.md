# Frontend Flows and UX

## Routing
- `/` Live Map + Feed + Drop Locus composer.
- `/loci/[id]` Thread view with map preview and replies.
- `/login`, `/signup` Auth pages.

## Data fetching
- **fetchLoci:** GET `/api/get/loci/?SouthWestLat=-90&SouthWestLong=-180&NorthEastLat=90&NorthEastLong=180` (public). Normalizes ids, coords, counts.
- **fetchReplies:** GET `/api/loci/{id}/replies` (public). Returns newest-first.
- **fetchLocus:** Convenience that calls `fetchLoci` then `fetchReplies`.

## Posting flows
- **DropLocus component:**
  - Uses geolocation when available; user can submit a message to `createLocus` helper.
  - Requires bearer token; on success, new locus is prepended to state so map/feed update instantly.
- **Composer in ReplyThread:** Posts replies via `createReply` helper and prepends to thread.

## Auth handling
- `login/register` helpers persist `accessToken`/`refreshToken` and user info in `localStorage`.
- `TopBar` shows Login/Signup or Logout depending on stored token; Logout clears storage and redirects.
- API helpers attach `Authorization: Bearer <token>` when present.

## Maps
- Leaflet client component (`MapView`) renders circles sized by view count and tooltips with view/reply counts.
- Centers on first locus with non-zero coords; falls back to default coordinates otherwise.

## View reporting
- `viewReporter` batches loci seen on screen (IntersectionObserver) and calls `registerView` in bursts.
- Calls are skipped without a token to avoid noisy 4xx.

## Error handling
- Loci fetch returns `[]` on non-OK to avoid crashing the page; errors surface as inline banners.
- Reply creation throws on non-OK so the UI can surface an error.
- Map fallback prevents 0/0 centering when data is missing.
