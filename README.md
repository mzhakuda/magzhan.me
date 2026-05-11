# magzhan.me

Code for https://magzhan.me/.

## Technology

- Go (stdlib `net/http`, `html/template`, `embed`)
- [HTMX](https://htmx.org) for the live-updating time in the bio
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
internal/site/         profile struct: timezone, links, socials
internal/i18n/         locale loader; bio placeholder substitution
internal/handlers/     home, 404, /api/time
locales/en.json        all user-facing copy
templates/             layout.html + per-page content blocks
static/                CSS, woff2 fonts, vendored htmx, og.png, fav.svg
```

## Customization

Two files hold everything personal:

- `internal/site/config.go` — timezone, `Links`, `Socials`, language list
- `locales/en.json` — title, tagline, bio paragraphs, city, company

Bio paragraphs support inline HTML and these placeholders:

- `{time}` — wrapped in an HTMX-polling span (updates every 60s)
- `{todays_date}` — current date in the profile timezone
- `{city}` — `Locale.City`
- `{company_name}` / `{company_url}` — from the locale

## Inspiration

Aesthetic and structure cribbed from [muan/site](https://github.com/muan/site).
