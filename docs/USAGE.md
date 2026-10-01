# ReconHawk CLI Usage

## Help

    reconhawk --help
    reconhawk scan --help
    reconhawk subdomains --help
    reconhawk historical --help
    reconhawk js --help
    reconhawk fingerprint --help
    reconhawk monitor --help

## Scan

    reconhawk scan example.com --depth 2 --workers 8 --timeout 10s --output report.json

## Passive subdomains

    reconhawk subdomains example.com

## Historical URLs

    reconhawk historical example.com

## JavaScript endpoints

    reconhawk js https://example.com

## Technology fingerprint

    reconhawk fingerprint https://example.com

## Continuous monitoring

    reconhawk monitor example.com --interval 10m
    reconhawk monitor example.com --interval 1m --iterations 10

Monitoring emits JSON snapshots suitable for logs or downstream automation.

## Plugin API

A plugin executable receives the target as its first argument and writes one JSON object to stdout.

Example:

    {"finding":"custom observation","value":"example"}

The Go API exposes RunPlugin so applications can integrate language-specific tooling without coupling it to the core.

Only use ReconHawk against authorized assets.
