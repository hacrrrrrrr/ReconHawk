# Architecture

ReconHawk is designed as a modular reconnaissance pipeline.

## Principles

1. Authorization first.
2. Bounded concurrency.
3. Observation over exploitation.
4. Structured output.
5. Scope-aware discovery.

## Planned flow

Target → scope validation → discovery providers → normalization → deduplication → report
