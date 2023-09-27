#!/bin/sh

go run ./cmd/examples/example.go -config config.yaml \
--verbose \
--timeout 10s \
--port 8080 \
--start "2020-01-01T16:08 EAT" \
--url="https://www.google.com" \
--uuid "123e4567-e89b-12d3-a456-426614174000" \
-ip="192.168.100.5" \
--mac "00:11:22:33:44:55" \
--email="email@example.com" \
--hostport="localhost:8080" \
--file "test.sh" \
--dir "/home" \
greet -name "John Doe" -greeting "Tsup,"