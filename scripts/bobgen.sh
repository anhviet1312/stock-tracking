#!/bin/sh
PSQL_DSN=postgres://postgres:postgres@127.0.0.1:5432/storyweb?sslmode=disable go run github.com/stephenafamo/bob/gen/bobgen-psql@v0.38.0 -c bobgen.yaml