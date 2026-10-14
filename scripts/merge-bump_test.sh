#!/bin/sh
set -eu
here=$(CDPATH= cd -- "$(dirname "$0")" && pwd)
dir=$(mktemp -d)
trap 'rm -rf "$dir"' EXIT
cat > "$dir/gh" <<'GH'
#!/bin/sh
echo "$*" >> "$GH_LOG"
case "$*" in
  *--auto*) [ -n "$GH_AUTO_ERR" ] || exit 0; echo "$GH_AUTO_ERR" >&2; exit 1 ;;
esac
GH
chmod +x "$dir/gh"
export PATH="$dir:$PATH" GH_LOG="$dir/log"

: > "$GH_LOG"
GH_AUTO_ERR= sh "$here/merge-bump.sh" release/core-v1 >/dev/null
[ "$(cat "$GH_LOG")" = "pr merge release/core-v1 --squash --auto --delete-branch" ]

: > "$GH_LOG"
GH_AUTO_ERR='GraphQL: Pull request Pull request is in clean status (enablePullRequestAutoMerge)' sh "$here/merge-bump.sh" release/core-v1 >/dev/null 2>&1
[ "$(cat "$GH_LOG")" = "pr merge release/core-v1 --squash --auto --delete-branch
pr merge release/core-v1 --squash --delete-branch" ]

: > "$GH_LOG"
if GH_AUTO_ERR='HTTP 403' sh "$here/merge-bump.sh" release/core-v1 >/dev/null 2>&1; then
  echo "merge-bump.sh hid an auto-merge error" >&2
  exit 1
fi
[ "$(cat "$GH_LOG")" = "pr merge release/core-v1 --squash --auto --delete-branch" ]
echo "merge-bump tests passed"
