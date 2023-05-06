#!/bin/sh
go clean -testcache
go test -v -run TestAcceptance ./internal/tests/