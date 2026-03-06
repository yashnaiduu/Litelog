# Roadmap

This document outlines the planned development trajectory for LiteLog.

> This roadmap reflects current priorities and may evolve based on community feedback and contributions.

---

## Phase 1 — Foundation ✅ Complete

- HTTP ingestion server
- SQLite storage engine (WAL mode, batched writes)
- `litelog start`, `ingest`, `query`, `tail`, `dashboard`, `export`
- Async goroutine-based flush pipeline
- Log retention policies
- GoReleaser cross-platform binary releases
- CI/CD via GitHub Actions
- Documentation site

---

## Phase 2 — Query & Streaming Enhancements

- **Regex filters** in `litelog tail` (e.g., `--match "panic.*db"`)
- **`json_extract` support** in SQL queries for nested log fields
- **Named queries** — save and re-run common query patterns
- **Colorized tail output** with configurable themes
- **Custom TUI charts** in the dashboard (error rate graphs, throughput sparklines)

---

## Phase 3 — Integrations

- **Docker logging driver** — route container logs directly to LiteLog
- **SD-notify support** — systemd service integration
- **HTTP push target** — forward logs to an external LiteLog instance
- **OpenTelemetry ingest endpoint** — accept OTLP-formatted log payloads

---

## Phase 4 — Scale-Out

- **Distributed mode** — multiple LiteLog instances with Raft consensus
- **Read replicas** — follower nodes for query-only workloads
- **Remote WAL sync** — stream WAL frames to remote storage (S3, GCS)

---

## Community

See [GitHub Issues](https://github.com/yashnaiduu/Litelog/issues) to vote on or propose new features.  
Check [CONTRIBUTING.md](CONTRIBUTING.md) to get involved.
