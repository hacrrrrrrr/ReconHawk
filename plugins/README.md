# ReconHawk Polyglot Plugins

ReconHawk keeps the high-concurrency orchestration and CLI in Go, while optional external engines can be written in other languages.

## Supported plugin languages

- Go — native/core extensions
- Python — parsing, enrichment, custom security research logic
- Java — JVM-based analyzers and enterprise integrations
- JavaScript/TypeScript — browser/JS ecosystem analysis
- Rust — optional high-performance specialized modules
- Shell — lightweight glue and provider adapters

Plugins communicate through a simple JSON contract so the core does not need to embed another runtime.

A plugin receives the target as its first argument and emits a JSON object on stdout.

Do not execute untrusted plugins. Plugins execute with the operating-system permissions of the user running ReconHawk.

## Layout

    plugins/
    ├── go/
    ├── python/
    ├── java/
    ├── javascript/
    ├── rust/
    └── shell/

Language runtimes remain optional; the normal ReconHawk installation only requires Go.
