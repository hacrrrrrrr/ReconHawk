# 🦅 ReconHawk

**Automated reconnaissance and attack-surface discovery for bug hunters.**

ReconHawk is a Go-based reconnaissance framework designed to turn an authorized target into structured, deduplicated attack-surface data. It focuses on discovery and observation rather than exploitation.

> **Authorized use only.** Run ReconHawk only against assets you own or targets explicitly permitted by the applicable bug-bounty program, penetration-test scope, or security-research authorization.

## ✨ Goals

- Fast, lightweight reconnaissance written in Go
- Concurrent HTTP/HTTPS probing with bounded workers
- Security-header, TLS, redirect and server fingerprint collection
- robots.txt, sitemap.xml and security.txt discovery
- Basic same-origin link extraction
- Deterministic JSON output for automation
- Clear separation between discovery and later analysis
- Designed to run comfortably on modest hardware

## 🧭 Pipeline

Target → HTTP probe → well-known resources → same-origin links → JSON report

## 🚀 Quick start

### Requirements

- Go 1.22+
- Network access to targets you are authorized to test

### Build

    git clone https://github.com/hacrrrrrrr/ReconHawk.git
    cd ReconHawk
    go build -o reconhawk ./cmd/reconhawk

### Run

    ./reconhawk scan https://example.com
    ./reconhawk scan https://example.com --output recon.json
    ./reconhawk scan https://example.com --workers 8 --timeout 8s

## 📦 Output

ReconHawk emits structured JSON containing observations, status codes, redirects, headers, TLS metadata and discovered same-origin URLs.

## 🗺️ Roadmap

- [x] Go CLI foundation
- [x] Concurrent HTTP probing
- [x] Redirect and TLS metadata
- [x] Security-header inventory
- [x] robots/sitemap/security.txt discovery
- [x] Same-origin URL extraction
- [ ] Passive subdomain provider interface
- [ ] DNS resolution module
- [ ] JavaScript endpoint extraction
- [ ] Historical URL providers
- [ ] Technology fingerprint database
- [ ] Scope-aware filtering
- [ ] JSONL / CSV / HTML reports
- [ ] Configurable rate limits and retry budgets
- [ ] Plugin/module API
- [ ] Continuous monitoring mode

## 🧱 Project layout

    ReconHawk/
    ├── cmd/reconhawk/       # CLI entrypoint
    ├── internal/recon/      # discovery engine
    ├── docs/                # design and usage docs
    ├── examples/            # example material
    ├── tests/               # integration-test material
    ├── go.mod
    ├── LICENSE
    └── README.md

## 🤝 Contributing

Issues and pull requests are welcome. Please keep contributions focused on authorized reconnaissance, reproducibility, reliability and clean output.

## 📫 Sponsorship & inquiries

**Kritik Bhattarai**  
Email: **hunterkritik@gmail.com**

## 📄 License

GPL-3.0. See LICENSE.
