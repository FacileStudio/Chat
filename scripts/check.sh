#!/usr/bin/env sh
#
# Repository quality gate for Chat.
#
#   sh scripts/check.sh           gofmt + vet + test + filet
#   sh scripts/check.sh --format  rewrite Go sources in place

set -eu

mode="all"
case "${1:-}" in
--format) mode="format" ;;
"") ;;
*)
  echo "usage: $0 [--format]" >&2
  exit 2
  ;;
esac

cd "$(git rev-parse --show-toplevel)"

if [ -z "${GO:-}" ]; then
  if [ -n "${GOROOT:-}" ] && [ -x "$GOROOT/bin/go" ]; then GO="$GOROOT/bin/go"; else GO=go; fi
fi
if [ -z "${GOFMT:-}" ]; then
  if [ -n "${GOROOT:-}" ] && [ -x "$GOROOT/bin/gofmt" ]; then GOFMT="$GOROOT/bin/gofmt"; else GOFMT=gofmt; fi
fi

if ! command -v "$GO" >/dev/null 2>&1 && [ ! -x "$GO" ]; then
  echo "check: no usable go ('$GO')" >&2
  exit 1
fi

if [ "$mode" = "format" ]; then
  "$GO" fmt ./...
  exit 0
fi

status=0

unformatted="$("$GOFMT" -l . | grep -v '^vendor/' || true)"
if [ -n "$unformatted" ]; then
  echo "gofmt: unformatted files (run 'sh scripts/check.sh --format'):"
  echo "$unformatted"
  status=1
fi

"$GO" vet ./... || status=1
"$GO" test ./... || status=1

if command -v filet >/dev/null 2>&1; then
  filet check . || status=1
fi

if [ "$status" -ne 0 ]; then
  echo "check failed"
fi
exit "$status"
