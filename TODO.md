# TODO

## High Priority

### Reading status
- Books have `owned` boolean but no way to track "read" status
- Need to add `read` boolean or `date_read` field to `books` table
- Consider: should this be a status enum (OWNED, READ, WISHLIST) or separate fields?
- Requires: migration, sqlc query updates, GraphQL schema changes, resolver updates

### Pagination metadata
- List queries (`books`, `authors`, `series`, `tags`) return raw arrays with `limit`/`offset` but no `totalCount`
- Clients cannot determine if more pages exist
- Consider `Connection` type with `nodes`, `totalCount`, `hasNextPage`

### Database unique constraints
- `authors.name` — no unique constraint, allows duplicate authors
- `books.isbn` / `books.isbn13` — no unique constraint, allows duplicate books
- `ScanAndAddBook` relies on `GetAuthorByName` which is fragile without uniqueness

### N+1 query problems
- `Series.Books` resolver fires one `GetBook` query per `SeriesBook` row
- `Book.tags` resolver fires one `GetBookTags` query per book when listing
- Consider batch queries or DataLoader pattern

### No graceful shutdown
- `main.go` uses `http.ListenAndServe` without signal handling
- Server does not drain connections on SIGTERM/SIGINT
- Add `os.Signal` handling and `http.Server.Shutdown()`

### No HTTP timeout on Open Library client
- `http.Client` uses default settings (no timeout)
- Slow Open Library responses could hang the server
- Add configurable timeout (e.g. 10s default)

### ScanAndAddBook always sets owned=true
- Scanned books are always marked as owned
- User may want to scan to wishlist without marking as owned
- Add an `owned` parameter to `scanAndAddBook`

### No CORS
- Server has no CORS headers
- Browser-based clients will be blocked
- Add CORS middleware with configurable allowed origins

## Medium Priority

### Tag CRUD gaps
- No `deleteTag` mutation
- No `updateTag` mutation (cannot rename tags)
- No `tag(id: ID!)` query (only list all tags)

### No local text search
- `books` query filters by `owned`, `authorID`, `tagID` but not by title/description
- Add `searchBooks(query: String!)` query
- Would need full-text search index (GIN/tsvector)

### SimilarBooks returns source book
- `SimilarBooks` and `SimilarSeries` search Open Library by title
- Results include the exact book/series being queried
- Filter out source from results

### No structured GraphQL errors
- Custom errors (`ValidationError`, `NotFoundError`, `DuplicateError`) returned as raw strings
- Clients cannot easily distinguish error types
- Implement `gqlerror.Extension` for structured error codes

### No Dockerfile
- No multi-stage Dockerfile for building/deploying
- Add Dockerfile for containerized deployment

### No CI pipeline
- No GitHub Actions or CI workflow
- Add CI for linting, testing, and building

### coverage.out tracked in git
- `.gitignore` only excludes `.env` and `bin/`
- Add `coverage.out` and other common Go ignores

## Lower Priority

### Reading status
- Books have `owned` boolean but no way to track "read" status
- Need to add `read` boolean or `date_read` field to `books` table
- Consider: should this be a status enum (OWNED, READ, WISHLIST) or separate fields?
- Requires: migration, sqlc query updates, GraphQL schema changes, resolver updates

### No ratings or reviews
- No ability to rate or review books
- Add `Rating` type with mutations/queries

### No bulk operations
- Cannot add/remove multiple tags from a book at once
- Cannot add multiple books to a series
- Add bulk mutations like `addTagsToBook(bookID: ID!, tagIDs: [ID!]!)`

### No request logging middleware
- No middleware for logging requests, measuring latency, or panic recovery
- Add logging middleware

### No input sanitization
- User input not sanitized beyond parameterized queries
- Add sanitization for string fields

### Magic numbers for default limits
- Default limits hardcoded: `10` for similar books, `20` for Open Library search
- Define constants or make configurable

### No API versioning
- GraphQL endpoint has no versioning strategy
- Consider `/graphql/v1` or schema-level versioning

### Repetitive dbToGraphQL conversion
- Each resolver manually converts DB types to GraphQL types
- Consider centralized conversion layer or code generation

### Utility functions scattered
- `ptrToStr`, `ptrToInt`, `ptrToBool` live at bottom of `schema.resolvers.go`
- Move to shared `internal/util` package

### No Open Library caching
- Repeated lookups hit the API every time
- Add in-memory LRU cache for responses

### No Open Library retry logic
- No retry for transient failures (rate limiting, network errors)
- Add retry with exponential backoff

### No config validation
- `config.Load()` silently uses defaults if env vars missing
- Validate required values (e.g. `DATABASE_URL` not empty)

### No index on series_books.position
- Queries ordering by position would benefit from index
- Add `CREATE INDEX idx_series_books_position ON series_books(series_id, position)`

### No index on book_tags.tag_id
- No index for reverse lookups (all books with a given tag)
- Add `CREATE INDEX idx_book_tags_tag_id ON book_tags(tag_id)`

### No updated_at trigger
- `updated_at` manually set in SQL queries
- Add PostgreSQL trigger to auto-update

### Only one migration file
- All schema in single `001_initial.sql`
- Split into multiple migrations as features are added

### No integration tests
- All tests use mocks, no real PostgreSQL tests
- Add integration tests with testcontainers or test database

### No Open Library client unit tests
- `client_test.go` only has real API tests gated by build tag
- Add unit tests with `httptest.NewServer`

### GetBookByKey unused
- Defined in client interface but never called from any resolver
- Either use it or remove it
