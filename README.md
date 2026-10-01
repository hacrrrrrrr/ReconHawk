# 🦅 ReconHawk

**Automated reconnaissance and attack-surface discovery for bug hunters.**

ReconHawk is a Go-based framework for organizing authorized reconnaissance into structured, deduplicated data. It focuses on discovery and observation rather than exploitation.

> **Authorized use only.** Use ReconHawk only against assets you own or targets explicitly permitted by the applicable bug-bounty program, penetration-test scope, or security-research authorization.

## Goals

- Lightweight Go-based reconnaissance
- Bounded concurrency and configurable timeouts
- HTTP/HTTPS observation and metadata collection
- Security-header and TLS inventory
- robots.txt, sitemap.xml and security.txt discovery
- Same-origin URL discovery
- JSON output for automation
- Modular architecture for future discovery providers

## Planned pipeline

```text
Target → Scope → HTTP discovery → DNS/passive sources → URL/JS discovery → deduplication → report
```

## Roadmap

- [x] Go CLI foundation
- [x] Concurrent HTTP probing
- [x] DNS resolution
- [ ] Passive subdomain providers
- [ ] JavaScript endpoint extraction
- [ ] Historical URL providers
- [ ] Technology fingerprinting
- [x] Same-origin scope filtering
- [x] JSON and JSONL reports
- [x] Bounded concurrency and request timeout
- [ ] Plugin API
- [ ] Continuous monitoring

## Project layout

```text
ReconHawk/
├── cmd/reconhawk/       # CLI
├── internal/recon/      # discovery engine
├── docs/                # documentation
├── examples/            # examples
├── tests/               # tests
├── go.mod
├── LICENSE
└── README.md
```

## Contributing

Issues and pull requests are welcome. Keep contributions focused on authorized reconnaissance, reliability, reproducibility and clean output.

## Sponsorship & inquiries

**Kritik Bhattarai**  
Email: **hunterkritik@gmail.com**

## License

GPL-3.0. See LICENSE.
