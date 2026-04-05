# radar-service

Radar service: gRPC data layer for web scraping, browsing, and crawling. Uses Crawl4ai and browser-use for agentic data extraction.

## Architecture (Phase 1)

- **core/** – Spider registry, proxy pool, scheduler, agent orchestrator
- **spiders/** – Per-source adapters (Indeed, LinkedIn, Bloomberg stubs)
- **extractors/** – Per-source parsers (Indeed extractor)
- **discovery/** – DiscoveryService gRPC (CrawlSource RPC)

## From repo root

```bash
make radar-init           # First-time: install deps, Playwright
make radar-generate-proto # Regenerate protos in hs-protos (run from users root)
make radar-run            # Start gRPC server on :50051
make up-founder           # web + postgres + redis + radar (compose, with hot reload)
```

Radar gRPC: `localhost:50051` (or `radar:50051` from other containers).

**Services:**
- `JobSearchService` – FindJobs, AnalyzeJob, ApplyToJob
- `DiscoveryService` – CrawlSource (spider orchestration)

## Docker (Bridgr Compose)

From **`bridgr-api`** repo root:

```bash
make up   # or: docker compose up --build -d
```

This builds and starts the **`radar`** service (`docker/radar/Dockerfile`) on **`0.0.0.0:50051`** (host port **50051**). The image generates Python gRPC stubs from **`proto/radar`** (same tree as the Go module) so it does **not** clone the private **`hs-protos`** Git dependency used by local `uv sync`. The Go API and worker read **`RADAR_ADDR`** via `config/docker.development.env` (default **`radar:50051`** on the Compose network); `internal/config.Config.RadarAddr` exposes this for wiring.

**Dependencies:** the Docker image omits **`browser-use`** (only referenced in TODO stubs); Discovery / JobSearch / PDF paths use **crawl4ai**, **playwright**, **python-jobspy**, etc. For LLM-backed extraction set **`OPENAI_API_KEY`**, **`GEMINI_API_KEY`**, or **`LLM_URL`** in `docker.development.env` or compose `environment` overrides.

## Docker entrypoint (like Go services)

- `run` – from /app (standalone, no volume)
- `run-local` – from mounted path (no reload)
- `watch` – from mounted path with hot reload (default in compose)
