# Changelog

All notable changes to LiteLog are documented here.

Format follows [Keep a Changelog](https://keepachangelog.com/en/1.0.0/).  
This project adheres to [Semantic Versioning](https://semver.org/).

---

## [Unreleased]

### Planned
- Regex filters for `litelog tail`
- Custom TUI charts in the dashboard
- `json_extract` support in SQL queries
- Docker logging driver
- Distributed instances via Raft consensus

---

## [0.1.0] — 2026-03-07

### Added
- HTTP ingestion server on `localhost:8080`
- SQLite storage engine with WAL mode for concurrent writes
- `litelog start` — starts the server with configurable retention (`--retention`)
- `litelog ingest` — pipe stdin directly into LiteLog
- `litelog query` — run standard SQL against the log database
- `litelog tail` — stream live logs filtered by `--level` and `--service`
- `litelog dashboard` — full-screen BubbleTea terminal dashboard
- `litelog export` — export to JSON or CSV
- Async goroutine-based batch flush pipeline
- Automatic log retention enforcement
- GoReleaser-based cross-platform binary releases (Linux, macOS, Windows)
- GitHub Actions CI (build + vet on push/PR)
- Docusaurus documentation site with GitHub Pages deployment
