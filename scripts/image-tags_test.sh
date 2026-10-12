#!/bin/sh
set -eu
tags=$(CDPATH= cd -- "$(dirname "$0")" && pwd)/image-tags.sh
out=$("$tags" v3.1.0-r1)
[ "$out" = "v3.1.0-r1
v3.1.0
latest" ]
out=$("$tags" v3.1.0)
[ "$out" = "v3.1.0
latest" ]
out=$("$tags" v3.4.0-rc.1)
[ "$out" = "v3.4.0-rc.1" ]
out=$("$tags" main)
[ "$out" = "latest" ]
ids=$(printf '%s\n' '[
  {"id":1,"created_at":"2020-01-01T00:00:00Z","metadata":{"container":{"tags":["v1.2.3-rc.1"]}}},
  {"id":2,"created_at":"2026-10-01T00:00:00Z","metadata":{"container":{"tags":["v1.2.3-rc.2"]}}},
  {"id":3,"created_at":"2020-01-01T00:00:00Z","metadata":{"container":{"tags":["latest"]}}}
]' | sh "$(dirname "$tags")/prune-rc-images.sh" --select 2026-09-12T00:00:00Z)
[ "$ids" = "1" ]
wf=$(CDPATH= cd -- "$(dirname "$0")/.." && pwd)/.github/workflows/release.yml
if grep -n '^[^ ]' "$wf" | grep -q 'ghcr.io/infinimesh/platform-operator'; then
  echo "release.yml image line is outside the run block" >&2
  exit 1
fi
