# magzhan.me

Code for https://magzhan.me/.

## Technology

- Go (stdlib `net/http`, `html/template`, `embed`)
- Self-hosted [IBM Plex Sans](https://www.ibm.com/plex/) (woff2)

## Development

Requires Go 1.22+.

```
$ go run ./cmd/web
```

Then open <http://127.0.0.1:8080>. The server listens on `127.0.0.1:8080` by
default; override with `-addr`.

## Layout

```
cmd/web/main.go        entry: routing, embed FS, graceful shutdown
embed.go               //go:embed for templates, locales, static
internal/site/         profile struct: links, languages
internal/i18n/         locale loader
internal/handlers/     home, 404, cv
locales/en.json        all user-facing copy
templates/             layout.html + per-page content blocks
static/                CSS, woff2 fonts, CV PDF, og.png, fav.svg
```

## Customization

Two files hold everything personal:

- `internal/site/config.go` — `Links` (shown on the home page), language list
- `locales/en.json` — name (title), tagline, CV page labels

## Inspiration

Aesthetic and structure cribbed from [muan/site](https://github.com/muan/site).
