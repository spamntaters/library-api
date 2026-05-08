# library-api

Go GraphQL API for a personal library. PostgreSQL backend with Open Library integration.

## Workflow

- **Always work in a feature branch.** Never commit directly to `main`. Use the naming convention `feat/<short-description>`.

## Commands

```
make up          Start PostgreSQL via docker compose
make down        Stop PostgreSQL
make migrate-up  Run goose migrations
make migrate-down Rollback goose migrations
make generate    Run gqlgen then sqlc (order matters)
make dev         generate + run server
make run         Run server only
make test        go test ./... -v
make test-coverage  Run with coverage report
```

## Architecture

- **Entrypoint**: `cmd/server/main.go`
- **GraphQL**: gqlgen, schema at `gql/schema.graphql`, playground at `/query`, endpoint at `/graphql`
- **Database**: sqlc generates Go from `internal/db/queries/*.sql` → `internal/db/*.go`
- **Migrations**: goose in `migrations/`
- **External API**: `internal/openlibrary/` client for openlibrary.org
- **Config**: loaded from `.env` via `internal/config/`

## Key Gotchas

### Code generation
- `make generate` runs **gqlgen first, then sqlc**. gqlgen needs sqlc types to compile resolvers.
- gqlgen validates by compiling the package — it will fail if the code doesn't build. Run `go mod tidy` if you get missing dependency errors during generation.
- Helper functions (`dbToGraphQLBook`, `ptrToInt`, etc.) live at the bottom of `schema.resolvers.go`. gqlgen preserves code below the generated section.

### GraphQL ID type
- All GraphQL `ID!` fields map to Go `int` (not `string`), configured via `gqlgen.yml`:
  ```yaml
  models:
    ID:
      model: github.com/99designs/gqlgen/graphql.IntID
  ```

### Testing
- Mocks live in `internal/graphql/resolver_test.go` (`MockStore`).
- `MockStore` must implement `db.Querier` — when sqlc changes function signatures (e.g. adding Params structs), update both the mock and any `store.On(...)` call sites in `resolvers_test.go`.
- Use `testResolver(store, olClient)` from `resolver_test.go` to construct a resolver for unit tests.

### Env vars
- `.env` (gitignored) — copy `.env.example` for defaults:
  ```
  DATABASE_URL=postgres://library:library@localhost:5432/library?sslmode=disable
  PORT=8080
  OPEN_LIBRARY_BASE_URL=https://openlibrary.org
  ```
