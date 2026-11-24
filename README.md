# Loci

A Go API for location-based messaging: users pin short messages (loci) to map coordinates, reply in threads, and track views.

## Why

Most conversations belong to a room or a feed, not a place. Loci ties every thread to a coordinate, so discussion stays local to where it happened, map-viewport queries are first-class, and connected clients can receive new activity nearby in real time.

## Quickstart

You need Go 1.24.5 or newer (`go version`), Docker with Compose v2.1 or newer (`docker compose version`), and ports 3000, 5432, and 6379 free.

```bash
git clone https://github.com/wycliff-ochieng/Loci.git
cd Loci
docker compose up -d --wait db redis
go run ./cmd/main.go
```

The server applies the migrations in `internal/store/migrations/` at startup, prints an `OK` line naming each applied migration, and listens on `http://localhost:3000`. Leave it running, then in a second terminal:

```bash
curl -s 'http://localhost:3000/api/get/loci/?SouthWestLat=-1.5&SouthWestLong=36.6&NorthEastLat=-1.0&NorthEastLong=37.0'
```

An empty database returns:

```json
[]
```

Run `go run` from the repository root, because migrations are resolved relative to the working directory. The defaults in `internal/config` match `docker-compose.yaml`, so no `.env` file is needed for local development.

## Usage

This walkthrough registers a user, logs in, pins a locus near Nairobi, reads it back, replies, and records a view. It needs the server from the Quickstart.

```bash
BASE=http://localhost:3000

# 1. Create a user.
curl -s -X POST "$BASE/register" \
  -H 'Content-Type: application/json' \
  -d '{"username":"alice","firstname":"Alice","lastname":"Wanjiku","email":"alice@example.com","password":"correcthorse"}'

# 2. Log in and keep the access token.
TOKEN=$(curl -s -X POST "$BASE/login" \
  -H 'Content-Type: application/json' \
  -d '{"email":"alice@example.com","username":"alice","password":"correcthorse"}' \
  | sed -E 's/.*"AccessToken":"([^"]+)".*/\1/')

# 3. Pin a locus and keep its id.
LOCUS_ID=$(curl -s -X POST "$BASE/api/post/loci" \
  -H "Authorization: Bearer $TOKEN" \
  -H 'Content-Type: application/json' \
  -d '{"message":"First locus","location":{"lat":-1.2864,"long":36.8172}}' \
  | sed -E 's/.*"ID":"([^"]+)".*/\1/')

# 4. Read the loci inside a bounding box around it.
curl -s "$BASE/api/get/loci/?SouthWestLat=-1.5&SouthWestLong=36.6&NorthEastLat=-1.0&NorthEastLong=37.0"

# 5. Reply to the locus and record a view.
curl -s -X POST "$BASE/api/loci/$LOCUS_ID/reply" \
  -H "Authorization: Bearer $TOKEN" \
  -H 'Content-Type: application/json' \
  -d '{"content":"First reply"}'

curl -s -X POST "$BASE/api/loci/$LOCUS_ID/view" \
  -H "Authorization: Bearer $TOKEN"
```

Step 2 prints the user, `AccessToken`, and `RefreshToken`. Step 4 returns the locus with its location and counters. After step 5, the same bounding-box query returns `view_count` 1 and `replies_count` 1:

```json
[{"id":"e560a268-d3da-4d00-a516-1e6ed95e1633","user_id":"085a4e24-903d-4698-a454-64d1722a1da6","message":"First locus","location":{"lat":-1.2864,"long":36.8172},"created_at":"2026-11-24T08:42:20.381741Z","view_count":1,"replies_count":1}]
```

Reads of loci and replies are public. Posting a locus, replying, recording a view, and opening the WebSocket require the bearer token.

### Real-time events

Connect to `/ws` with the same bearer token. Browser WebSocket clients cannot set an `Authorization` header, so when you proxy through the bundled nginx configuration, pass the token in the query string (`/ws?token=$TOKEN`); the proxy turns it into the header the API expects. For a non-browser client, send the header directly (this uses `$BASE` and `$TOKEN` from the walkthrough above):

```bash
curl -s -i -N --max-time 3 -H "Connection: Upgrade" -H "Upgrade: websocket" \
  -H "Sec-WebSocket-Key: dGhlIHNhbXBsZSBub25jZQ==" -H "Sec-WebSocket-Version: 13" \
  -H "Authorization: Bearer $TOKEN" "$BASE/ws"
```

A `101 Switching Protocols` response means the connection is open. Send your position so the hub can filter events for you:

```json
{"type":"USER_UPDATE_LOCATION","payload":{"lat":-1.2864,"long":36.8172}}
```

The hub pushes JSON envelopes with a `type` field and a `payload` field. The types are `LOCUS_NEW`, `REPLY_EVENT`, and `VIEW_EVENT`. A client that never sends a location receives no events.

### Endpoints

| Method | Path | Auth | Purpose |
| --- | --- | --- | --- |
| POST | `/register` | none | Create a user from `username`, `firstname`, `lastname`, `email`, `password`. |
| POST | `/login` | none | Log in with `email` (or `username`) and `password`; returns `User`, `AccessToken`, `RefreshToken`. |
| GET | `/api/get/loci/` | none | List loci inside a bounding box, the rectangle between two opposite corners: `SouthWestLat`, `SouthWestLong`, `NorthEastLat`, `NorthEastLong`. |
| POST | `/api/post/loci` | JWT | Create a locus from `message` and `location: {lat, long}`; the author comes from the token. |
| GET | `/api/loci/{id}/replies` | none | List a locus's replies, newest first. |
| POST | `/api/loci/{id}/reply` | JWT | Add a reply from `content`. |
| POST | `/api/loci/{id}/view` | JWT | Record at most one view per user per locus; increments `view_count` once. A repeat call returns an error without changing the count. |
| GET | `/ws` | JWT | WebSocket connection for real-time locus, reply, and view events. |

## Configuration

`internal/config` reads environment variables and loads a `.env` file from the directory you run the server in. The defaults point at the services in `docker-compose.yaml`:

| Variable | Default | What it changes |
| --- | --- | --- |
| `DB_HOST` | `localhost` | Postgres host. |
| `DB_PORT` | `5432` | Postgres port. |
| `DB_NAME` | `loci_db` | Database name. |
| `DB_USER` | `admin` | Database user. |
| `DB_PASSWORD` | `secretpassword` | Database password. |
| `DB_SSLMODE` | `disable` | Postgres SSL mode. |
| `REDIS_ADDR` | `localhost:6379` | Redis address used by the per-user posting limiter. |
| `REDIS_PASSWORD` | `secretpassword` | Redis password. |
| `JWTSecret` | `mydogiscalledrufus` | Secret the auth middleware validates tokens with. Leave it at the default: login tokens are currently signed with this value from `pkg/auth`, so a different `JWTSecret` makes every issued token fail validation. |
| `CORS_ALLOWED_ORIGINS` | `http://localhost:3001` | Comma-separated browser origins allowed to call the API. Set it to the web client's origin. |

For example, to run the API against infrastructure on non-default ports:

```bash
DB_PORT=15432 REDIS_ADDR=localhost:16379 go run ./cmd/main.go
```

## Contributing

Run `make lint`, `make test`, and `make build` before opening a pull request. The repository has no tests yet, so a patch that adds behavior should add the first test for the package it touches. If you change SQL in `internal/store/queries/`, regenerate the typed code with `sqlc generate` and commit `sqlc/`; if you change the schema, add a goose migration to `internal/store/migrations/`. Ask questions and propose changes in an issue.

## Deployment

The container image, Kubernetes manifests, nginx proxy configs, and CI pipeline are covered in [docs/deployment.md](docs/deployment.md). Deeper design notes live in [docs/architecture.md](docs/architecture.md), [docs/geolocation-and-postgis.md](docs/geolocation-and-postgis.md), [docs/replies-and-views.md](docs/replies-and-views.md), and [docs/websocket.md](docs/websocket.md).

## License

No license has been specified for this repository yet; add one before distributing or reusing the code.
