# gofs template

The project template that `gofs init` copies. A full-stack Go app with two small, dependency-free runtime dependencies (chi's middleware and go-chi/cors) and nothing else:

- **Server**: Go standard library (`net/http`, `html/template`, `embed`) plus [chi](https://github.com/go-chi/chi)'s middleware (logger, recoverer, gzip, no-cache, strip slashes) and [go-chi/cors](https://github.com/go-chi/cors). Routes in one file, handlers next to their pages. Middleware stays in `internal/server/middleware`.
- **Pages**: rendered whole on the server with `html/template`. Each folder under `internal/ui/app/pages` is a section with its templates and handlers side by side. Pages include `layouts/head` and `layouts/foot` for the document shell.
- **Interactivity**: declarative shadow DOM components in `internal/ui/app/components`, one `.gohtml` per element. The browser attaches the shadow roots as it parses the page, with no JavaScript. Behavior, when needed, is a small inline `<script type="module">` in the same file, outside the shadow template; only the contact list needs it, to hide and show the inline edit rows.
- **Data**: an in-memory example store in `internal/app/contact`. Postgres via plain SQL and [sqlc](https://sqlc.dev) when you need it (`sqlc.yaml`, `queries.sql`, `internal/db`).
- **Styling**: plain CSS in `internal/server/assets/css/app.css`, plus scoped styles inside each shadow root.
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

The template ships a small contacts app so that every pattern has a working example. Each is one page and one handler file:

| Page            | Shows                                                                                  |
| --------------- | -------------------------------------------------------------------------------------- |
| `/contacts`     | a GET search form, a "Load more" link, inline edit rows (toggled by a small script in the component) and delete as a plain form post |
| `/contacts/new` | a plain form post with server-rendered validation errors                               |

Every control is a working link or form before JavaScript loads: the edit forms are rendered into the page and stay visible without JS, and every state change is a plain form post followed by a redirect.

Delete the `contact` package, the example pages and their components when you start your own app.

## Components

Use a component for a region that owns its markup, styles or behavior, such as the site header or the contacts list. One component is one `.gohtml` file under `internal/ui/app/components/<tag>/`:

- the host element and a `<template shadowrootmode="open">` holding the markup and a `<style>`; the browser attaches it while parsing, so it renders with no JavaScript;
- optional behavior in an inline `<script type="module">` after the host in the same file (scripts inside the shadow template are inert), defining the custom element;
- include it with `{{ template "components/<tag>" . }}`.

Styles inside the shadow root are scoped and inherit the app's custom properties. Global CSS is only for page layout and light-DOM primitives. Keep one-off page content (home, not found, the new-contact form) as page templates.

## Layout

```txt
cmd/server              main
internal
  app/contact           example backend package: domain rules, store, queries.sql
  auth                  user from request
  config                configuration from environment
  db                    database connection, transactions, migrations
  server
    assets/css          app.css
    handlers            generic handlers (hashed static files)
    middleware          CORS and the chi middleware stack
    routes.go           every route
  ui/app
    layouts.gohtml      document head and foot
    pages/<section>     templates + handlers
    components/<tag>    <tag>.gohtml (optional inline module script)
sqlc.yaml               sqlc configuration
```

## Database

The example store is in memory so the template runs with nothing installed. To use Postgres: `go get github.com/jackc/pgx/v5`, add the blank import in `internal/db/db.go`, set `DSN` in `.env` (`make run` then starts Postgres with docker compose), and run `make generate` to turn `internal/app/contact/queries.sql` into a typed `query` package. The migrations in `internal/db/migrations` run at startup.

Objects are stored as JSON blobs with generated columns for the fields you search, sort or index. See the gofs docs for the reasoning.
