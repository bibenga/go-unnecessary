# go install github.com/go-jet/jet/v2/cmd/jet@latest

jet -dsn=postgresql://postgres:postgres@db:5432/postgres?sslmode=disable \
    -path=dbjet \
    -schema=public \
    -ignore-tables schema_migrations \
    -path db-jet
