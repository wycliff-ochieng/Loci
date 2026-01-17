# WebSockets, Hub, and Event Flows

## Overview
- A central hub (`internal/socket.Hub`) coordinates connections and fan-out of real-time events for loci, replies, and views.
- Services push events into hub channels after successful DB commits; hub relays to subscribed clients over a single `/ws` endpoint (auth-protected).
- Current frontend does not yet consume WS events; HTTP polling/fetch still works. The design keeps WS opt-in but ready.

## Hub design
- **Channels:**
  - `BroadcastLocus chan *models.Locus`
  - `BroadcastReply chan *models.ReplyEvent`
  - `BroadcastView chan *models.ViewEvent`
  - `Register/Unregister` for client connection management.
- **Responsibilities:**
  - Track active clients and their outbound channels.
  - Fan-out any incoming broadcast message to all connected clients.
- **Lifecycle:** hub is started in `api.Run()` and runs as a goroutine.

## Client model
- **Endpoint:** `GET /ws` (auth required). Middleware attaches user context before the WS upgrade.
- **Connection:** Each client is registered with the hub; when they disconnect, `Unregister` is called.
- **Outbound:** Clients receive whatever events are pushed into the hub channels; message format can be JSON-encoded structs for locus/reply/view payloads.

## Event producers
- **Create locus:** `UserService.CreateLoci` pushes a `models.Locus` to `BroadcastLocus` after the transaction commits.
- **Create reply:** `ReplyLoci` pushes `models.ReplyEvent` (reply id, locus id, username, locus lat/long, createdAt) to `BroadcastReply` after commit.
- **Record view:** `RecordView` pushes `models.ViewEvent` (userId, locusId, locus lat/long, viewedAt) to `BroadcastView` after commit.

## Delivery semantics and trade-offs
- **At-least-once broadcast:** Events are pushed after DB commit. If the hub channel is congested or the client disconnects, delivery may be dropped on the floor. Data is always persisted; real-time is best-effort.
- **No per-topic filtering yet:** All clients receive all events. Trade-off: simplicity vs. bandwidth. Future: add room/topic filtering (e.g., by bounding box or locus id) and per-client subscriptions.
- **Backpressure:** Hub uses buffered channels, but if clients stop reading, messages can queue. Future: drop-oldest or disconnect slow consumers.
- **Ordering:** Events are ordered per channel but not globally across channels (locus vs reply vs view). Consumers should tolerate out-of-order across different event types.
- **Reliability vs complexity:** Chose simple broadcast over durable queues; acceptable for MVP where page refresh can reconcile via HTTP fetch.

## Views logic (WS context)
- Views are counted via HTTP POST (`/api/loci/{id}/view`) with auth. The service transaction inserts a view row and increments `view_count`, then emits a `ViewEvent` for real-time counters. Clients without WS still see updated counts via subsequent fetches. Batching on the frontend (`viewReporter`) reduces HTTP noise; WS would remove the need for client-side polling for other users' views.

## Replies logic (WS context)
- Replies are created via HTTP POST (`/api/loci/{id}/reply`) with auth. Counts are incremented in-transaction; a `ReplyEvent` is emitted so other clients can append the new reply without refetching. Without WS, the creating client updates optimistically; others would refetch to see it.

## Why this approach
- **Simplicity first:** A single hub and broadcast channels reduce moving parts while the product stabilizes.
- **DB as source of truth:** Events only emit after commit, ensuring clients never see phantom writes.
- **Extensible:** Channels are separated per event type; adding filters (bounding boxes, locus rooms) or payload enrichments is straightforward.
- **Auth consistency:** Reuses the same JWT middleware for `/ws` so server trusts the user context before registering the client.

## Future improvements
- Topic/room subscriptions (e.g., per locus id, or per map viewport) to cut bandwidth.
- Heartbeats and pings to detect dead clients faster and reclaim resources.
- Backpressure strategy: drop-oldest or disconnect slow clients based on buffer depth.
- Unified envelope for WS messages with `type` and `payload` to simplify client parsing.
- Optional reliable delivery (ack/rewind) if product needs guaranteed real-time, otherwise stay best-effort.
