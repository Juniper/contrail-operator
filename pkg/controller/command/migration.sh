#!/usr/bin/env bash

export PGPASSWORD=contrail123

psql -w -h localhost -U root -d postgres <<END_OF_SQL
DROP DATABASE contrail_test_migrated;
DROP DATABASE contrail_test_backup;
END_OF_SQL

set -e
createdb -w -h localhost -U root contrail_test_migrated
psql -w -h localhost -U root -d contrail_test_migrated -f /usr/share/contrail/gen_init_psql.sql
psql -w -h localhost -U root -d contrail_test_migrated -f /usr/share/contrail/init_psql.sql

commandutil convert --intype yaml --in /backups/db_migrated.yml --outtype rdbms -c /etc/contrail/migration.yml

psql -w -h localhost -U root -d postgres <<END_OF_SQL
BEGIN;
ALTER DATABASE contrail_test RENAME TO contrail_test_backup;
ALTER DATABASE contrail_test_migrated RENAME TO contrail_test;
COMMIT;
END_OF_SQL