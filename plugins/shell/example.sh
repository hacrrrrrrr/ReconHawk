#!/usr/bin/env sh
target="${1:-}"
printf '{"plugin":"shell-example","target":"%s","type":"enrichment"}\n' "$target"
