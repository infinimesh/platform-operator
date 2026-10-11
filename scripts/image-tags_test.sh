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
out=$("$tags" main)
[ "$out" = "latest" ]
