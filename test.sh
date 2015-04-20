#!/bin/bash

echo "Create test database"
psql -U postgres -c 'DROP DATABASE IF EXISTS test_convos;'
psql -U postgres -c 'CREATE DATABASE test_convos;'

echo "Run test migrations"
migrate -url "postgres://localhost/test_convos?sslmode=disable" -path ./migrations reset

echo "Run tests"
go test ./...
