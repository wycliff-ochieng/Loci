# Replies and Views

## Replies pipeline
- **Endpoint:** `POST /api/loci/{id}/reply` (auth required).
- **Flow:**
  1. Handler parses `content`, pulls `userID` from JWT context.
  2. Service runs a transaction:
     - Insert into `replies` (`CreateReply`).
     - `IncrementReplyCount` on the parent locus.
     - Fetch reply + locus lat/long for broadcast (`GetReplyForBroadcast`).
  3. Commit; broadcast `ReplyEvent` on hub (id, locusId, username, locus lat/long, createdAt).
  4. Returns `models.Reply` with JSON fields (`replyid`, `locusid`, `userid`, `content`, `createdat`).
- **Frontend:**
  - `createReply` posts JSON `{content}` with bearer token.
  - Response is normalized to `Reply` shape and prepended to the thread.
  - Reply counts come from the locus payload (`replies_count`).

## View pipeline
- **Endpoint:** `POST /api/loci/{id}/view` (auth required).
- **Flow:**
  1. Handler parses `locusID`, gets `userID` from JWT.
  2. Service transaction:
     - Insert into `views` (`CreateView`).
     - `IncrementViewCount` on loci.
  3. Commit; fetch locus lat/long; broadcast `ViewEvent` (userId, locusId, lat/long, viewedAt).
  4. Returns a `View` payload.
- **Frontend:**
  - `registerView` posts with bearer token. If no token, it no-ops to avoid 424 spam.
  - `viewReporter` batches on-screen loci and flushes every 5s (best-effort, errors swallowed).
  - View counts arrive with each locus (`view_count`).

## Thread loading
- `GET /api/loci/{id}/replies` (public) returns replies ordered newest-first.
- `ReplyThread` loads replies on mount and appends newly created replies optimistically.

## Consistency and correctness
- Reply and view counters are updated inside the same transaction as the insert to avoid skew.
- Broadcast happens after commit; if the hub is unavailable, data is still persisted.
- Frontend tolerates missing auth on view calls and missing fields in older payloads via normalization.

## Common failure modes
- **401/424 on view:** happens when no bearer token; mitigated by skipping calls when unauthenticated.
- **Empty content:** handler rejects empty replies.
- **Invalid locus id:** handlers respond 400 on bad UUIDs.
