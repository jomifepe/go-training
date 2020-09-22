#!/usr/bin/env sh

chmod +x ./build-run.sh
reflex -d none -sr '(.*\.(go|sh)|go\.mod)$' -- ./build-run.sh