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
root=$(CDPATH= cd -- "$(dirname "$0")/.." && pwd)
ruby -ryaml -e 'ARGV.each {|p| YAML.load_file(p)}' "$root"/.github/workflows/*.yml
bad=$(mktemp)
trap 'rm -f "$bad"' EXIT
cat > "$bad" <<'EOF'
name: bad
on: push
jobs:
  image:
    runs-on: ubuntu-latest
    steps:
      - run: |
          echo ok
`not-yaml`.
EOF
if ruby -ryaml -e 'YAML.load_file(ARGV[0])' "$bad" >/dev/null 2>&1; then
  echo "a workflow with an unindented line in a run block was accepted" >&2
  exit 1
fi
