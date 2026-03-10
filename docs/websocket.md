## WebSocket Architecture and Flow

### Components and Their Roles
- `ServerWS` (handler): Upgrades HTTP to WebSocket using Gorilla upgrader. Requires `userID` in context (set by auth middleware). Creates a `Client`, registers it with the hub, and starts pumps.
- `Client` (connection state): Holds the socket (`Conn`), outbound buffer (`Send`), user ID, and last known `Location`. Runs `ReadPump` and `WritePump` goroutines.
- `ReadPump` (inbound): Reads frames from the socket, decodes the envelope `{type, payload}`, and currently handles `USER_UPDATE_LOCATION` to update `Client.Location`. On errors, unregisters and closes the connection.
- `WritePump` (outbound): Listens on `Send` and writes messages to the socket; batches queued messages; sends periodic ping to keep the connection alive; closes on errors.
- `Hub` (orchestrator): Tracks connected clients (`Clients` map) and channels (`Register`, `UnRegister`, `BroadcastLocus`). Its `Run` loop adds/removes clients and fans out broadcasts to eligible clients.
- Message envelope: `Websocket{Type string, Payload interface{}}` serialized as JSON. Current outbound type: `LOCUS_NEW`; inbound type: `USER_UPDATE_LOCATION`.

### Data Flow (current behavior)
1) **Connect:** Client hits `/ws` (route to be registered). `ServerWS` upgrades to WS, builds `Client{Send: buffered chan}`, registers it via `Hub.Register`, starts `ReadPump` and `WritePump`.
2) **Location update:** Frontend sends `{type:"USER_UPDATE_LOCATION", payload:{lat,long}}`. `ReadPump` parses and stores `Client.Location` (mutex-protected). Clients without a location will be skipped for geo-filtered pushes.
3) **Broadcast new locus:** Service layer calls `hub.BroadcastNewLoci(locus)` → sends locus to `BroadcastLocus` channel.
4) **Hub fan-out:** `Hub.Run` wraps locus as `{type:"LOCUS_NEW", payload:locus}`, JSON-encodes, and for each client with a location, computes distance (`models.CalculateDistance`). If within 5km, it enqueues the message onto `client.Send`; if the channel is blocked, the client is dropped.
5) **Deliver:** `WritePump` drains `Send`, writes frames, batches any queued messages into the same write, and sends periodic pings to keep the connection alive.
6) **Disconnects:** Read/write errors trigger unregister → removal from `Clients` map and closing of the send channel.

### Why It’s Structured This Way
- **Separation of concerns:** Handler does upgrade + auth context; Hub manages membership and broadcast; Client pumps isolate IO from coordination.
- **Backpressure safety:** Send channel buffering + drop-on-block prevents a slow client from stalling the hub.
- **Geo-filtered relevance:** Distance check keeps broadcasts local (5km radius). Clients with no location are skipped to avoid noisy sends.
- **Keep-alive:** Ping timer avoids idle disconnects by intermediaries.
- **Typed envelopes:** Simple `{type, payload}` lets the frontend switch on message intent without multiple endpoints.

### Gaps / Next Steps
- Route wiring: add `router.HandleFunc("/ws", func(w,r){ ServerWS(hub,w,r) })` and ensure auth middleware sets context.
- Payload decoding for location could be simplified (currently gob+json); direct JSON is cleaner.
- Add pong/timeout handling and per-client close logging.
- Expose reply/view events over WS when backend supports them.