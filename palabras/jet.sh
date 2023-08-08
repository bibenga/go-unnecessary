# go install github.com/go-jet/jet/v2/cmd/jet@latest

jet -source sqlite \
    -dsn palabras.db \
    -path ./palabras/gen \
    -schema pl \
    -ignore-tables schema_migrations 
