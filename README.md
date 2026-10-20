# gofs template

The project template that `gofs init` copies. A full-stack Go app with two small, dependency-free runtime dependencies (chi's middleware and go-chi/cors) and nothing else:

- **Server**: Go standard library (`net/http`, `html/template`, `embed`) plus [chi](https://github.com/go-chi/chi)'s middleware (logger, recoverer, gzip, no-cache, strip slashes) and [go-chi/cors](https://github.com/go-chi/cors). Routes in one file, handlers next to their pages. CSRF and security headers are in `internal/server/middleware`.
- **Pages**: rendered on the server with `html/template`. Each folder under `internal/ui/app/pages` is a section with its templates and handlers side by side.
- **Interactivity**: small native web components in `internal/ui/app/components`, one folder per element. They ask the server for HTML fragments and swap them into the page. The one shared helper is `request.js`.
- **Data**: an in-memory example store in `internal/app/contact`. Postgres via plain SQL and [sqlc](https://sqlc.dev) when you need it (`sqlc.yaml`, `queries.sql`, `internal/db`).
- **Styling**: plain CSS in `internal/server/assets/css/app.css`.
- **Nothing else**: no Node.js, no npm, no bundler, no JavaScript libraries. `go.mod` has two requirements, neither of which has dependencies of its own.

## Run

```sh
make          # start the server with live reload at http://localhost:8080
make test     # go test
make lint     # go vet, staticcheck, govulncheck
make build    # build ./tmp/server
```

Development tools (sqlc, wgo, staticcheck, govulncheck) are pinned in the Makefile and run with `go run pkg@version`, so they never appear in `go.mod`.

## Examples

The template ships a small contacts app so that every pattern has a working example. Each is one page, one handler file and one or two components:

| Page            | Components                                              | Shows                                                                                    |
| --------------- | ------------------------------------------------------- | ---------------------------------------------------------------------------------------- |
| `/contacts`     | `live-search`, `load-more`, `edit-inline`, `row-delete` | search as you type, paging, edit a row in place, delete with a `<dialog>` confirmation   |
| `/contacts/new` | `validate-field`                                        | server-side validation of one field while typing; a plain form post with a flash message |

Every control is a working link or form before JavaScript loads. Handlers answer a plain form with a redirect and a component with a fragment; `ui.WantsFragment` tells the two apart.

Delete the `contact` package, the example pages and their components when you start your own app.

## Layout

```txt
cmd/server              main
internal
  app/contact           example backend package: domain rules, store, queries.sql
  auth                  user from request
  config                configuration from environment
  db                    database connection, transactions, migrations
  server
    assets/css|js       app.css and app.js (component loader, dialog fallback)
    handlers            generic handlers (hashed static files)
    middleware          CSRF and security headers (the rest comes from chi)
    routes.go           every route
  ui/app
    layout.tmpl         document shell
    pages/<section>     templates + handlers
    components/<tag>    element.js per custom element; request.js shared
sqlc.yaml               sqlc configuration
```

## Database

The example store is in memory so the template runs with nothing installed. To use Postgres: `go get github.com/jackc/pgx/v5`, add the blank import in `internal/db/db.go`, set `DSN` in `.env` (`make run` then starts Postgres with docker compose), and run `make generate` to turn `internal/app/contact/queries.sql` into a typed `query` package. The migrations in `internal/db/migrations` run at startup.

Objects are stored as JSON blobs with generated columns for the fields you search, sort or index. See the gofs docs for the reasoning.
