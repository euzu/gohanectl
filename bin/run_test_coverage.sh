#!/bin/bash
cd hanectl && go test ./... -coverprofile=coverage.out -coverpkg=./... && go tool cover -html=coverage.out
