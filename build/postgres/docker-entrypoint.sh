#!/bin/bash
set -e

psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" <<-EOSQL
	CREATE USER bacsgateway; ALTER user bacsgateway with password 'bacsgatewaypwd';
	CREATE DATABASE bacsgateway OWNER bacsgateway;
EOSQL
