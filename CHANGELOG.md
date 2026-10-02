# Changelog

## 0.3.0 — 2026-10-02

### Added
- Working Go-first CLI for scan, subdomains, historical, js, fingerprint, check, monitor, cache and version.
- Bounded concurrent same-origin crawling.
- Context-aware HTTP requests and redirect limits.
- Persistent JSON cache with atomic writes and cache TTL support.
- Certificate Transparency passive subdomain provider.
- Common Crawl historical URL discovery.
- JavaScript/HTML endpoint extraction.
- HTTP, TLS and DNS observations.
- Technology fingerprinting.
- Non-destructive security checks with evidence and remediation.
- Continuous monitoring with configurable intervals and finite iterations.
- JSON and JSONL output.
- CI checks for formatting, tests, vet and build.

### Fixed
- Initial documentation-only repository state now has an executable implementation.
- Cache writes are protected from partially written files by temporary-file replacement.
- Crawl concurrency is bounded to prevent unbounded request fan-out.
- Redirect handling is limited.
- Response bodies are size-limited to avoid excessive memory use.
- Passive discovery results are normalized and deduplicated.

### Safety
ReconHawk is intended for authorized security research only. It does not attempt exploitation or destructive actions.

## 0.1.0
- Initial project documentation.
