#!/usr/bin/env sh

echo "sh: Building..." && go build -o server && echo "sh: Running..." && ./server