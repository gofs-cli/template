# AGENTS.md

Guidance for AI agents working in this repository: the project template copied by `gofs init`. A full-stack Go app with server-rendered pages, declarative shadow DOM and almost no JavaScript. The shipped contacts app is an example to be deleted when building your own app.

## Hard constraints

- Runtime dependencies are only the Go standard library, `github.com/go-chi/chi/v5` (middleware only, never its router) and `github.com/go-chi/cors`; neither has transitive dependencies. A Postgres driver (`github.com/jackc/pgx/v5`) is the one addition an app may make. Do not add others.
- No Node.js, npm, bundlers or JavaScript libraries. Never add `package.json`. The only client JavaScript is a small inline `<script type="module">` in the `contact-list` component template; everything else is server-rendered.
- Templates are `html/template` files with the `.gohtml` suffix (not `.tmpl`), embedded via `//go:embed app` in `internal/ui/ui.go`; `make run` watches them and rebuilds.
- Handlers hold no business logic; that belongs in `internal/app/<package>`.
- Pages are rendered whole. There is no fragment/fetch layer and no client-side rendering: forms post, handlers redirect or re-render, links navigate.
- Every control works as a link or form before JavaScript loads; the component's inline module script only enhances server-rendered markup (the contact list hides its edit rows until Edit is pressed). Without it the forms stay visible and submit normally.
- Browser features: Baseline widely available freely; newly available only with a fallback in the component that needs it.

## Commands

```sh
make            # cp .env (from example), start Postgres if DSN is set, run with live reload (wgo)
make test       # go test ./cmd/... ./internal/...
make lint       # go vet, staticcheck, govulncheck
make build      # build ./tmp/server
make fmt        # gofmt -w cmd internal
make generate   # sqlc generate (see sqlc.yaml)
make gendata    # docker compose down/up, then re-migrates the DB via a gendata-tagged test
make dbup       # start Postgres with docker compose (only if DSN is set in .env)
make dbdown
```

- Development tools (sqlc, wgo, staticcheck, govulncheck) are pinned in the Makefile and run with `go run pkg@version`, so they never appear in `go.mod`. `make run` (wgo) watches `.go`, `.gohtml`, `.css`, `.js` and `.sql` files for live reload.
- CI (`.github/workflows/`) runs `make test`, `make lint`, `make build`, a gofmt check (`gofmt -l cmd internal`), and gitleaks.
- There are no unit tests shipped. The only `_test.go` is `internal/data/data_test.go`, which has the `gendata` build tag and is *not* run by `make test`; `make gendata` runs it against Postgres.

## Layout and flow

Requests flow: `cmd/server/main.go` → importing `internal/ui` parses all templates (panics on error) → `server.New` (opens DB, seeds stores) → `Routes()` → page handlers → `ui.Render` → the template set in the embed FS.

- `internal/server/routes.go` — every route, with Go 1.22 `ServeMux` patterns (`GET /contacts`, `POST /contacts/{id}`, read with `r.PathValue("id")`). Two muxes: assets (revalidated via ETag, CORS GET) vs. everything else (no-cache pages, auth, recover, CORS).
- `internal/server/middleware.go` — chi middleware only. Note lists are `slices.Reverse`d before wrapping, so *first in the literal list = outermost*.
- `internal/ui/app/pages/<section>/` — one folder per section: `handlers.go` (handler constructors), `model.go` (view data), `*.gohtml` templates. Handlers are functions taking their dependencies (stores) and returning `http.Handler`, e.g. `home.Index()`, `contacts.Index(s.contacts)`.
- `internal/ui/app/components/<tag>/` — declarative shadow DOM components: `<tag>.gohtml` holds the markup, styles and any behavior. Behavior is an inline `<script type="module">` next to the host in the same file; there are no separate component JS assets.
- `internal/app/<package>` — domain rules and stores. The example is `internal/app/contact` (a thread-safe in-memory store with `ValidationError`).
- `internal/server/assets/` — `app.css`, embedded and served with crc32 ETags via `handlers.NewHashedAssets` (`internal/server/handlers/assets.go`). There is no app script: behavior lives in the component template.

### Templates (`internal/ui`)

- `mustParse()` walks the embedded FS once, parses each `.gohtml` file into a single set under its path, and panics with the file involved on error. Every template is a `{{ define "path/style/name" }}` block; template names are path-style (`home/index`, `components/contact-row`, `layouts/head`, `components/site-header`).
- A duplicate `{{ define }}` name silently keeps the last one seen (html/template behaviour); there is no duplicate check.
- `components/contact-row` is an inline `{{ block }}` inside the rows `range`. `block` is shorthand for `define` + `template`, so the row is defined where it is used and remains available by name (`{{ template "components/contact-row" . }}`) for reuse elsewhere.
- `ui.Render(w, "contacts/index", data)` executes a named template. It buffers the output and answers 500 with the template error if execution fails.
- Every page starts with `{{ template "layouts/head" dict "Title" "..." }}` and ends with `{{ template "layouts/foot" . }}`; `layouts.gohtml` holds the document shell and the site header include. `head` expects a dict (or nil), so pass `dict` rather than a struct without `.Title`.
- The only custom template func is `dict`, for passing parameters with `{{ template "components/site-header" dict ... }}` (`funcs.go`). Missing keys read as nil, so use `{{ with .Key }}` for optional parameters.

### Frontend components (declarative shadow DOM)

Use a component when a region owns its markup, styles or behavior and can be reasoned about on its own: the site header, the contacts table with its search, the inline edit rows. Keep one-off page content (home, not found, the new-contact form) as page templates; do not wrap every piece of markup.

- One component is one file: `internal/ui/app/components/<tag>/<tag>.gohtml`, defining `components/<tag>`. The tag must contain a hyphen.
- The file holds the host and a `<template shadowrootmode="open">` with the markup and a `<style>`. The browser attaches the shadow root while parsing, so the markup renders with no JavaScript.
- Behavior, when needed, is an inline `<script type="module">` after the host in the same file. Scripts inside the shadow template are inert (the parser builds the template contents without a browsing context); the module script runs with the page and defines the custom element. Guard with `customElements.get(...)` in case the component is included twice.
- Include it with `{{ template "components/<tag>" . }}` (pass `dict` for parameters, like `layouts/head`).
- Styles are scoped to the shadow root and inherit custom properties (`--space-*`, `--color-*`). Host-level layout (e.g. `contact-list { display: block }`) goes in `app.css`; everything else about the component belongs in its `<style>`, not the global stylesheet.
- The pattern in one sentence: server-rendered markup and CSS in a declarative shadow root, optional behavior in a colocated module script, one file per component.

### Inline edit rows

- The `components/contact-list` component renders the search form, the count, the table and the "Load more" link. Each contact is two `<tr>`s: a display row with an Edit button and a delete form, and a `.row-edit` row with the edit form. The form's `action` (`/contacts/{id}`) and the hidden `q`/`page` fields are set server-side at render time.
- The inline module in the component hides `.row-edit` rows on load unless the server marked them `data-open`, and toggles them with `[data-edit]` / `[data-cancel]` buttons inside its shadow root. Validation errors re-render the list with the row marked `data-open` and the submitted values.
- Delete is a plain `POST /contacts/{id}/delete` form; search is a plain GET form; paging is a link. Nothing needs JavaScript.

## Configuration, auth, data

- Config comes from the environment (`internal/config/config.go`): `HOST`, `PORT`, `ALLOWED_ORIGINS`, `ENV` (`local`/`dev`/`prod`), `DSN`. `make` copies `.env.example` to `.env` if missing.
- `ENV=local` makes `auth.Middleware` inject `auth.LocalUser`; token parsing is an unimplemented TODO otherwise (`internal/auth/context.go:32`).
- The default data path is the in-memory `contact.Store` seeded with `SampleContacts()` — the app runs with zero infrastructure. Postgres is opt-in: add the pgx blank import in `internal/db/db.go`, set `DSN`, start it with docker compose (`make dbup`), and `make generate` to turn `queries.sql` into a typed package. Migrations in `internal/db/migrations` run **at server startup** and must be idempotent. Objects are stored as JSONB blobs with generated columns for searched/sorted fields (`migrations.sql`).
- `sqlc.yaml` has one entry per package with a `queries.sql`; generated code is not committed in the template (no `internal/app/contact/query` dir yet).

## Gotchas

- A component's host element (`<site-header>`) is unknown to the browser without declarative shadow DOM support; if you need older browsers, use a `<div>` wrapper instead or add a small polyfill.
- Component behavior is an inline module script, so a strict Content-Security-Policy would need a hash or nonce for it; no CSP is configured by default.
- Asset URLs are unversioned, so the assets middleware sends `Cache-Control: no-cache` and relies on ETags; do not add a long `max-age` without versioned URLs.
- `.env` is local-only (gitignored); never commit it. Do not commit your app's secrets into this template.