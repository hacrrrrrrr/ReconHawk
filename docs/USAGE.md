# ReconHawk CLI Usage

## Global help

    reconhawk --help
    reconhawk help

## Scan

    reconhawk scan example.com
    reconhawk scan example.com --depth 2 --workers 8 --timeout 10s --output report.json
    reconhawk scan example.com --jsonl observations.jsonl

## Passive subdomains

Uses certificate-transparency data through crt.sh.

    reconhawk subdomains example.com

## Historical URLs

Uses Common Crawl indexes.

    reconhawk historical example.com

## JavaScript endpoints

Fetches an authorized URL and extracts HTTP/HTTPS endpoints and same-origin paths.

    reconhawk js https://example.com

## Technology fingerprinting

    reconhawk fingerprint https://example.com

## Version

    reconhawk version

Every command accepts its own help flag:

    reconhawk scan --help
    reconhawk subdomains --help
    reconhawk historical --help
    reconhawk js --help
    reconhawk fingerprint --help

ReconHawk is reconnaissance-only. It does not exploit discovered endpoints.
