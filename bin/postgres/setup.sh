#!/bin/bash
set -e

psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" <<-EOSQL
    CREATE USER go_test_user WITH PASSWORD 'go_test_secret';
    CREATE DATABASE go_test;
    GRANT ALL PRIVILEGES ON DATABASE go_test TO go_test_user;
EOSQL

psql go_test --username "go_test_user" < /tmp/initial.sql