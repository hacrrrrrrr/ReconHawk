# 🦅 ReconHawk

**Automated reconnaissance and attack-surface discovery for bug hunters.**

ReconHawk is a Go-first reconnaissance framework for authorized security research. It collects and normalizes attack-surface observations without attempting exploitation.

> **Authorized use only.** Scan only assets you own or assets explicitly included in the applicable bug-bounty, pentest, or security-research scope.

## Features

- Concurrent HTTP/HTTPS probing
- DNS resolution
- Same-origin crawling
- robots.txt, sitemap.xml and security.txt discovery
- Passive subdomain discovery through certificate transparency
- Historical URL discovery through Common Crawl
- JavaScript/HTML endpoint extraction
- HTTP/TLS/security-header metadata
- Technology fingerprinting
- JSON and JSONL output
- Bounded concurrency and timeouts
- Modular provider architecture
- Go core with Python/Java extension path

## CLI

    reconhawk --help

Commands:

    scan
    subdomains
    historical
    js
    fingerprint
    version
    help

Examples:

    reconhawk scan example.com --depth 2 --workers 8 --output report.json
    reconhawk subdomains example.com
    reconhawk historical example.com
    reconhawk js https://example.com
    reconhawk fingerprint https://example.com

See [docs/INSTALL.md](docs/INSTALL.md) and [docs/USAGE.md](docs/USAGE.md).

## Architecture

    Target
      │
      ▼
    Go Core ── Scope ── Scheduler ── Dedup ── Output
      │
      ├── HTTP / DNS / crawler
      ├── Passive providers
      ├── Historical URLs
      ├── JS endpoint extraction
      └── Technology fingerprints
              │
              └── future Python / Java engines

Go remains the primary runtime. Other languages are optional specialized extensions, not required dependencies.

## Roadmap

- [x] Go CLI foundation
- [x] Concurrent HTTP probing
- [x] DNS resolution
- [x] Passive subdomain provider
- [x] JavaScript endpoint extraction
- [x] Historical URL provider
- [x] Technology fingerprinting
- [x] Scope-aware same-origin filtering
- [x] JSON / JSONL reports
- [x] Timeout and concurrency controls
- [x] Stable external plugin API
- [x] Continuous monitoring command
- [x] Additional passive providers
- [x] Persistent result cache

## Project layout

    ReconHawk/
    ├── cmd/reconhawk/
    ├── internal/recon/
    ├── docs/
    ├── examples/
    ├── tests/
    └── .github/

## Sponsorship & inquiries

**Kritik Bhattarai**  
**hunterkritik@gmail.com**

## License

GPL-3.0. See LICENSE.
