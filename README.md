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

### Docker Compose Deployment (Frontend + Backend + Nginx)
This repository now includes production deployment orchestration in `docker-compose.yaml` for:
- `frontend` (builds from `../geodrop-notes`)
- `backend` (Go API)
- `nginx` reverse proxy (TLS termination and routing)
- `db` (PostGIS)
- `redis`

#### 1) Prepare environment
Use one of the concrete env files in this folder:
- `.env.dev` for local/IP testing (HTTP)
- `.env.prod` for production (HTTPS)

Update at minimum:
- `JWTSecret`
- `DB_PASSWORD`
- `REDIS_PASSWORD`
- `CORS_ALLOWED_ORIGINS`
- `VITE_WS_PROXY_URL`

#### 2) Add TLS certificates
Place certificates in `nginx/certs/`:
- `nginx/certs/fullchain.pem`
- `nginx/certs/privkey.pem`

#### 3) Build and start
From this `Loci/` folder:
- Dev: `docker compose --env-file .env.dev up -d --build`
- Prod: `docker compose --env-file .env.prod up -d --build`

#### 4) Verify
- `docker compose ps`
- `docker compose logs -f nginx backend frontend`
- Open `https://<your-domain>`

#### 5) Deploy on a server
Copy these paths to your server (keeping the same relative layout):
- `Loci/`
- `geodrop-notes/`

Then run on server in `Loci/`:
- edit `.env.prod` with your real values
- add certs to `nginx/certs/`
- `docker compose --env-file .env.prod up -d --build`

docker compose --env-file .env.dev up -d --build

Dev/local (no domain): docker compose --env-file .env.dev up -d --build
Prod/domain: docker compose --env-file .env.prod up -d --build
Stop: docker compose --env-file .env.dev down (or .env.prod)

## CI/CD (Lint -> Test -> Build -> Push)

This repository includes a GitHub Actions workflow that runs:
- Lint: `gofmt` and `go vet`
- Test: `go test -v ./...`
- Build: Go binary build (`cmd/main.go`)
- Docker image build and push to:
	- Docker Hub
	- GitHub Container Registry (GHCR)

### Required GitHub Secrets
Add these in GitHub repository settings -> Secrets and variables -> Actions:
- `DOCKERHUB_USERNAME`
- `DOCKERHUB_TOKEN` (Docker Hub access token, not your password)

The workflow automatically uses `GITHUB_TOKEN` for GHCR push.

### Image Names Produced
- Docker Hub: `<DOCKERHUB_USERNAME>/loci-backend:<tag>`
- GHCR: `ghcr.io/<github-owner>/loci-backend:<tag>`

`latest` is pushed from the default branch, while branch/tag/SHA tags are also generated.

## VPS Deployment Without Domain (IP Only)

You can deploy with plain server IP using HTTP (no TLS certificate needed).

### 1) Prepare VPS
Install Docker and Compose plugin:

```bash
sudo apt update
sudo apt install -y ca-certificates curl gnupg
sudo install -m 0755 -d /etc/apt/keyrings
curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo gpg --dearmor -o /etc/apt/keyrings/docker.gpg
sudo chmod a+r /etc/apt/keyrings/docker.gpg
echo \
	"deb [arch=$(dpkg --print-architecture) signed-by=/etc/apt/keyrings/docker.gpg] https://download.docker.com/linux/ubuntu \
	$(. /etc/os-release && echo $VERSION_CODENAME) stable" | \
	sudo tee /etc/apt/sources.list.d/docker.list > /dev/null
sudo apt update
sudo apt install -y docker-ce docker-ce-cli containerd.io docker-buildx-plugin docker-compose-plugin
sudo usermod -aG docker $USER
```

Logout/login once after adding your user to the docker group.

### 2) Pull and configure project

```bash
git clone <your-repo-url>
cd Loci
cp .env.dev .env
```

Edit `.env`:
- Set strong values for `JWTSecret`, `DB_PASSWORD`, `REDIS_PASSWORD`
- Set `CORS_ALLOWED_ORIGINS=http://<YOUR_VPS_IP>`
- Keep `NGINX_CONF=nginx.dev.conf` for non-TLS HTTP deployment

### 3) Build and run stack

```bash
docker compose --env-file .env up -d --build
docker compose ps
docker compose logs -f nginx backend frontend
```

Open in browser:
- `http://<YOUR_VPS_IP>`

### 4) Use prebuilt backend image from CI (optional)
If you want to use the image pushed by GitHub Actions instead of building locally:
- In `docker-compose.yaml`, replace backend build block with image:

```yaml
backend:
	image: <DOCKERHUB_USERNAME>/loci-backend:latest
```

Then run:

```bash
docker compose --env-file .env pull backend
docker compose --env-file .env up -d
```

### 5) Update process on VPS

```bash
cd Loci
git pull
docker compose --env-file .env pull
docker compose --env-file .env up -d
docker image prune -f
```

## Push Latest Changes to GitHub

From your local machine:

```bash
cd Loci
git add .
git commit -m "ci: add loci lint-test-build-push workflow and VPS guide"
git push origin <your-branch>
```

After merging to default branch, the workflow pushes `latest` automatically to Docker Hub and GHCR.