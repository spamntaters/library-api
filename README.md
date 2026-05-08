# library-api

Go GraphQL API for managing a personal book library. Tracks books, authors, series, and tags with Open Library integration for ISBN scanning and book discovery.

## Features

- **Book management** — CRUD operations with ownership and read status tracking
- **Author management** — Link books to authors with optional bios
- **Series tracking** — Organize books into series with position ordering; find missing (unowned) books in a series
- **Tagging** — Categorize books with custom tags
- **Open Library integration** — Scan ISBNs to auto-populate book data, search for similar books/series, lookup by ISBN
- **GraphQL API** — Full CRUD via GraphQL with pagination support on list queries
- **GraphQL Playground** — Interactive query explorer at `/query`

## Prerequisites

- Go 1.25+
- Docker & Docker Compose
- [gqlgen](https://gqlgen.com/) — `go install github.com/99designs/gqlgen@latest`
- [sqlc](https://sqlc.dev/) — `go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest`
- [goose](https://github.com/pressly/goose) — `go install github.com/pressly/goose/v3/cmd/goose@latest`

## Quick Start

```bash
# 1. Clone and enter the project
cd library-api

# 2. Start PostgreSQL
make up

# 3. Run migrations
make migrate-up

# 4. Start the server
make dev
```

Server runs at `http://localhost:8080`. GraphQL playground at `http://localhost:8080/query`.

## Configuration

Copy `.env.example` to `.env` and adjust as needed:

```env
DATABASE_URL=postgres://library:library@localhost:5432/library?sslmode=disable
PORT=8080
OPEN_LIBRARY_BASE_URL=https://openlibrary.org
```

| Variable | Default | Description |
|---|---|---|
| `DATABASE_URL` | `postgres://library:library@localhost:5432/library?sslmode=disable` | PostgreSQL connection string |
| `PORT` | `8080` | HTTP server port |
| `OPEN_LIBRARY_BASE_URL` | `https://openlibrary.org` | Open Library API base URL |

## Commands

| Command | Description |
|---|---|
| `make up` | Start PostgreSQL via Docker Compose |
| `make down` | Stop PostgreSQL |
| `make migrate-up` | Run goose migrations |
| `make migrate-down` | Rollback goose migrations |
| `make generate` | Run gqlgen then sqlc (order matters) |
| `make dev` | Generate code and start server |
| `make run` | Start server only |
| `make test` | Run tests with verbose output |
| `make test-coverage` | Run tests with coverage report |

## API

### Endpoints

| Path | Description |
|---|---|
| `/graphql` | GraphQL API endpoint |
| `/query` | GraphQL Playground (interactive UI) |
| `/health` | Health check (returns 200 if database is reachable, 503 otherwise) |

### GraphQL Schema

#### Queries

| Query | Description |
|---|---|
| `books(owned, authorID, tagID, limit, offset)` | List books with optional filters |
| `book(id)` | Get a single book by ID |
| `authors(limit, offset)` | List authors |
| `author(id)` | Get a single author by ID |
| `series(limit, offset)` | List series |
| `seriesByID(id)` | Get a single series by ID |
| `seriesMissingBooks(seriesID, limit, offset)` | Get unowned books in a series |
| `similarBooks(bookID, limit)` | Find similar books via Open Library |
| `similarSeries(seriesID, limit)` | Find similar series via Open Library |
| `searchOpenLibrary(query, limit)` | Search Open Library |
| `lookupByISBN(isbn)` | Look up a book by ISBN on Open Library |
| `tags(limit, offset)` | List tags |

#### Mutations

| Mutation | Description |
|---|---|
| `createBook(input)` | Create a new book |
| `updateBook(id, input)` | Update a book |
| `deleteBook(id)` | Delete a book |
| `toggleOwned(id)` | Toggle book ownership |
| `toggleRead(id)` | Toggle book read status |
| `createAuthor(input)` | Create a new author |
| `updateAuthor(id, input)` | Update an author |
| `deleteAuthor(id)` | Delete an author |
| `createSeries(input)` | Create a new series |
| `updateSeries(id, input)` | Update a series |
| `deleteSeries(id)` | Delete a series |
| `addBookToSeries(bookID, seriesID, position)` | Add a book to a series |
| `removeBookFromSeries(bookID, seriesID)` | Remove a book from a series |
| `createTag(name)` | Create a new tag |
| `addTagToBook(bookID, tagID)` | Add a tag to a book |
| `removeTagFromBook(bookID, tagID)` | Remove a tag from a book |
| `scanAndAddBook(isbn)` | Scan ISBN via Open Library and add to library |

### Example Queries

```graphql
# List owned books with pagination
query {
  books(owned: true, limit: 10, offset: 0) {
    id
    title
    author { name }
    owned
    read
    tags { name }
  }
}

# Find missing books in a series
query {
  seriesMissingBooks(seriesID: 1, limit: 20) {
    id
    title
    author { name }
  }
}

# Scan an ISBN and add to library
mutation {
  scanAndAddBook(isbn: "9780451524935") {
    book { id title }
    external { title authorName }
    alreadyExists
  }
}
```

## Project Structure

```
├── cmd/server/main.go          # Application entrypoint
├── gql/
│   └── schema.graphql          # GraphQL schema definition
├── internal/
│   ├── config/                 # Environment configuration
│   ├── db/
│   │   ├── queries/*.sql       # SQL queries (sqlc source)
│   │   ├── *.sql.go            # Generated Go code from sqlc
│   │   └── store.go            # Database store with transaction support
│   ├── graphql/
│   │   ├── schema.resolvers.go # Resolver implementations
│   │   ├── helpers.go          # Conversion helpers (dbToGraphQLBook, etc.)
│   │   ├── validation.go       # Input validation
│   │   ├── generated.go        # Generated gqlgen code
│   │   └── *_test.go           # Tests
│   └── openlibrary/            # Open Library API client
├── migrations/                 # goose migration files
├── gqlgen.yml                  # gqlgen configuration
├── sqlc.yaml                   # sqlc configuration
└── Makefile                    # Developer commands
```

## Architecture

- **GraphQL**: [gqlgen](https://gqlgen.com/) generates Go code from `gql/schema.graphql`. Resolvers live in `internal/graphql/schema.resolvers.go`.
- **Database**: [sqlc](https://sqlc.dev/) generates type-safe Go code from SQL queries in `internal/db/queries/*.sql`.
- **Migrations**: [goose](https://github.com/pressly/goose) manages schema migrations in `migrations/`.
- **External API**: `internal/openlibrary/` provides a client for [openlibrary.org](https://openlibrary.org).
- **Configuration**: Loaded from `.env` via `internal/config/`.

## Development Workflow

### Adding a new feature

1. Create a feature branch: `git checkout -b feat/<short-description>`
2. If adding database changes:
   - Write a new migration in `migrations/`
   - Add or update SQL queries in `internal/db/queries/`
3. If adding GraphQL fields:
   - Update `gql/schema.graphql`
4. Run `make generate` (runs gqlgen first, then sqlc)
5. Implement resolvers in `internal/graphql/schema.resolvers.go`
6. Add helper functions to `internal/graphql/helpers.go` (not `schema.resolvers.go` — gqlgen will delete them on regeneration)
7. Write tests in `internal/graphql/resolvers_test.go`
8. Update the mock in `internal/graphql/resolver_test.go` if sqlc changed function signatures
9. Run `make test` to verify

### Code generation order

`make generate` runs **gqlgen first, then sqlc** because gqlgen needs sqlc-generated types to compile resolvers. If gqlgen fails due to missing types, run `sqlc generate` manually first, then `gqlgen generate`.

If you get missing dependency errors during generation, run `go mod tidy` first.

### Working with the database

```bash
# Start PostgreSQL
make up

# Run migrations
make migrate-up

# Rollback last migration
make migrate-down

# Stop PostgreSQL
make down
```

## Testing

```bash
# Run all tests
make test

# Run with coverage
make test-coverage

# Run a specific test
go test ./internal/graphql/... -run TestQuery_Books -v
```

Tests use mocks (`MockStore` in `internal/graphql/resolver_test.go`). The mock must implement `db.Querier` — when sqlc changes function signatures, update both the mock and any `store.On(...)` call sites in `resolvers_test.go`.

Open Library client tests are integration-only, gated by the `TEST_OL_API` build tag:

```bash
TEST_OL_API=1 go test ./internal/openlibrary/... -v
```

## Troubleshooting

### gqlgen fails with "undefined" errors

Run `go mod tidy`, then try `make generate` again. If helper functions are missing, run `sqlc generate` first to update db types, then `gqlgen generate`.

### Database connection refused

Ensure PostgreSQL is running: `make up`. Check `DATABASE_URL` in your `.env`.

### Port already in use

Change `PORT` in your `.env` or stop the process using port 8080.

### Migration fails

Check that the database exists and is accessible. Run `make down && make up` to restart PostgreSQL, then `make migrate-up`.
