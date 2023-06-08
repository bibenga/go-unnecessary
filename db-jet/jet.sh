# go install github.com/go-jet/jet/v2/cmd/jet@latest

jet -dsn=postgresql://rds:sqlsql@host.docker.internal:5432/go?sslmode=disable \
    -path=dbjet \
    -schema=public \
    -ignore-tables schema_migrations \
    -path db-jet
