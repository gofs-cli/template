# AGENTS.md

Guidance for AI agents working in this repository: the project template copied by `gofs init`.

## Hard constraints

- Runtime dependencies are the Go standard library, `github.com/go-chi/chi/v5` (middleware only, not its router) and `github.com/go-chi/cors`; both have no transitive dependencies. A Postgres driver is the one addition an app may make. Do not add others.
- No Node.js, npm, bundlers or JavaScript libraries. Never add `package.json`. Components are plain `HTMLElement` subclasses in `internal/ui/components/<tag>/element.js`.
- Pages are `html/template` files with the `.tmpl` suffix. Handlers hold no business logic; that belongs in `internal/app/<package>`. Components hold no business logic either: they fetch, swap and react.
- The server returns HTML fragments to components
- Every control works as a link or form before JavaScript loads; components enhance server-rendered HTML.
- Never build HTML strings from data in JavaScript. `innerHTML` only for same-origin fragments the server rendered; `textContent` for everything else.
- Browser features: Baseline widely available freely; newly available with a fallback (see `app.js`).

## Conventions

- All routes in `internal/server/routes.go`. A state-changing route answers a plain form with a 303 redirect and a component request (`ui.WantsFragment`) with a fragment.
- One component per folder, defined with `customElements.define` using the folder name as the tag. Shared code goes in `request.js` only when three components need the same lines.
- Development tools are pinned in the Makefile with `go run pkg@version`.

## Commands

```sh
make        # run with live reload
make test   # tests, including template parsing
make lint   # vet, staticcheck, govulncheck
make build
```
