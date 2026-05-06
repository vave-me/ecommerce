# Setup and External Integrations

Use this checklist before local E2E testing, staging, or production deployment.

## 1. Configure Environment Variables

Create local runtime env from template and then set required integration values:

```bash
cp docker/.env.example docker/.env
```


- AI providers:
  - `OPENAI_API_KEY`
  - `ANTHROPIC_API_KEY`
  - `DEEPSEEK_API_KEY`
  - `GOOGLE_AI_API_KEY`
- Geocoding and maps:
  - `GEOCODING_API_KEY`
  - `NOMINATIM_PASSWORD`
  - `PBF_URL`
  - `REPLICATION_URL`
  - `IMPORT_WIKIPEDIA`
- Google OAuth:
  - `WEB_GOOGLE_OAUTH_CLIENT_ID`
  - `MOBILE_GOOGLE_OAUTH_CLIENT_ID`
  - `NEXT_PUBLIC_GOOGLE_CLIENT_ID`
  - `GOOGLE_CLIENT_SECRET`
  - `GOOGLE_PROJECT_ID`
  - `GOOGLE_AUTH_URI`
  - `GOOGLE_ISSUER`
- MinIO/S3:
  - `MINIO_ROOT_USER`
  - `MINIO_ROOT_PASSWORD`
  - `MINIO_ENDPOINT`
  - `MINIO_ACCESS_KEY_ID`
  - `MINIO_SECRET_ACCESS_KEY`
  - `MINIO_BUCKET`
  - `MINIO_USE_SSL`
- Payments/shipping/mail:
  - `STRIPE_KEY`, `STRIPE_SECRET`, `STRIPE_WEBHOOK_SECRET`
  - `DHL_API_KEY`
  - `SMTP_HOST`, `SMTP_PORT`, `SMTP_USERNAME`, `SMTP_PASSWORD`/`SMTP_PASS`, `DEFAULT_SMTP_SENDER`

## 2. Download Geocoding Map Data

Nominatim expects OSM `.pbf` data under `docker/nominatim/`.

```bash
mkdir -p docker/nominatim
curl -L -o docker/nominatim/germany-latest.osm.pbf \
  https://download.geofabrik.de/europe/germany-latest.osm.pbf
```

Use these env values:

```bash
PBF_URL=file:///data/germany-latest.osm.pbf
REPLICATION_URL=https://download.geofabrik.de/europe/germany-updates
IMPORT_WIKIPEDIA=false
```

## 3. Start Stack

```bash
docker compose --env-file ./docker/.env --profile microservices up -d
```

## 4. Initialize MinIO Bucket

MinIO endpoints:

- S3 API: `http://localhost:9096`
- Console: `http://localhost:9099`

Create the bucket named by `MINIO_BUCKET` in the console.

Optional CLI method:

```bash
docker run --rm --network host minio/mc \
  alias set local http://127.0.0.1:9096 "$MINIO_ROOT_USER" "$MINIO_ROOT_PASSWORD"
docker run --rm --network host minio/mc \
  mb --ignore-existing "local/$MINIO_BUCKET"
```

## 5. Apply Bucket Read Policy (Requested)

Policy file location:

- `docker/minio/public-read-policy.json`

Current policy grants public read for `arn:aws:s3:::classified/*`.

Apply with `mc`:

```bash
docker run --rm --network host -v "$PWD/docker/minio:/policies:ro" minio/mc \
  alias set local http://127.0.0.1:9096 "$MINIO_ROOT_USER" "$MINIO_ROOT_PASSWORD"
docker run --rm --network host -v "$PWD/docker/minio:/policies:ro" minio/mc \
  anonymous set-json /policies/public-read-policy.json local/classified
```

If your bucket is not `classified`, update `Resource` in the policy file before applying.

## 6. Verify Geocoding Import

First import is long-running. Monitor:

```bash
docker compose logs -f nominatim
```

Check service status:

```bash
docker compose ps geocoding nominatim
```

## 7. Verify Frontend Integration

Ensure frontend points at reachable backend and OAuth settings:

- `NEXT_PUBLIC_API_BASE_URL`
- `NEXT_PUBLIC_GOOGLE_CLIENT_ID`

If backend APIs are not reachable, auth/search/geocoding UI flows will fail.

## 8. Validate Before Publish

```bash
go generate ./...
go test ./...
```

For frontend:

```bash
npm run lint
npm run build
npm run test -- --runInBand
```
