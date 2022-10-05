#!/bin/bash -eu

go test $(go list ./... | grep -v /vendor/) -coverprofile=cover.out -v ; go tool cover -html=cover.out -o cover.html
