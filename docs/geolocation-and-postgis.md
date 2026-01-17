# Geolocation and PostGIS

## Coordinate system
- **WGS84 (SRID 4326):** Lat/long stored as `geography` in PostGIS for accurate distance on the spheroid.
- **Insertion:** `ST_MakePoint(long, lat)::geography` when creating a locus.
- **Reading:** `ST_Y(location::geometry)` -> lat, `ST_X(location::geometry)` -> long for API payloads.

## Queries
- **Bounding box fetch:**
  - SQL: `ST_Within(location::geometry, ST_MakeEnvelope($1,$2,$3,$4,4326))`
  - Params: `(SouthWestLong, SouthWestLat, NorthEastLong, NorthEastLat)`
  - Service maps rows to `models.Locus` with lat/long and counts.
- **Locus location lookup:** `GetLocusLocation` uses `ST_Y/ST_X` to return precise coords for view broadcasts.

## Distance and envelope reasoning
- **Envelope vs radius:** Envelope keeps queries simple for map viewport fetches. For circular radius, swap to `ST_DWithin(location, ST_MakePoint(lng, lat)::geography, radius_meters)`.
- **Precision:** Geography type accounts for Earth curvature; geometry is only used for `ST_Within` on the envelope for performance.

## Geolocation on the client
- `DropLocus` tries `navigator.geolocation` to seed coords; falls back to provided seed or defaults.
- Frontend stores coords as `lat/lng`; API expects `{lat, long}`. JSON tags on `GeoPoint` ensure backend emits `lat/long` so the frontend mapping can read either casing.

## Data flow for a new locus
1. User submits message + coords.
2. Handler parses; middleware injects `userID` from JWT.
3. Service builds `CreateLociParams` with `ST_MakePoint(long, lat)`.
4. PostGIS stores geography; sqlc returns lat/long via `ST_Y/ST_X`.
5. Service maps to `models.Locus` and broadcasts on the hub (with counts).
6. Frontend updates state and re-renders map/feed.

## Gotchas and mitigations
- **0/0 coordinates:** If malformed data is encountered, frontend centers on a safe default instead of (0,0) and mapping tolerates multiple key shapes.
- **SRID consistency:** All envelopes use 4326; mixing SRIDs will break `ST_Within`.
- **Performance:** Envelope queries are index-friendly if a `GIST` index on geography exists (`CREATE INDEX loci_location_gix ON loci USING GIST(location);`).
