#!/bin/sh
# GitHub refuses auto-merge on a PR that is already mergeable ("is in clean
# status"), so merge that one directly.
set -eu
b=$1
if out=$(gh pr merge "$b" --squash --auto --delete-branch 2>&1); then
  printf '%s\n' "$out"
  exit 0
fi
printf '%s\n' "$out" >&2
case $out in
  *"is in clean status"*) gh pr merge "$b" --squash --delete-branch ;;
  *) exit 1 ;;
esac
