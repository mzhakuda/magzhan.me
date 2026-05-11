package i18n

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/fs"
	"path"
	"strings"
)

type Locale struct {
	Lang         string   `json:"lang"`
	LangTag      string   `json:"lang_tag"`
	Title        string   `json:"title"`
	Tagline      string   `json:"tagline"`
	Bio          []string `json:"bio"`
	City         string   `json:"city"`
	Company      string   `json:"company"`
	CompanyURL   string   `json:"company_url"`
	NavLinks     string   `json:"nav_links"`
	NavSocial    string   `json:"nav_social"`
	NavAbout     string   `json:"nav_about"`
	FooterSource string   `json:"footer_source"`

	bioPrepared []string
}

type Registry struct {
	byCode      map[string]*Locale
	defaultLang string
}

func Load(fsys fs.FS, dir, defaultLang string) (*Registry, error) {
	r := &Registry{byCode: map[string]*Locale{}, defaultLang: defaultLang}
	entries, err := fs.ReadDir(fsys, dir)
	if err != nil {
		return nil, fmt.Errorf("read locales dir: %w", err)
	}
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".json") {
			continue
		}
		b, err := fs.ReadFile(fsys, path.Join(dir, e.Name()))
		if err != nil {
			return nil, fmt.Errorf("read %s: %w", e.Name(), err)
		}
		var loc Locale
		if err := json.Unmarshal(b, &loc); err != nil {
			return nil, fmt.Errorf("parse %s: %w", e.Name(), err)
		}
		if loc.Lang == "" {
			return nil, fmt.Errorf("%s: missing lang field", e.Name())
		}
		loc.prepareBio()
		r.byCode[loc.Lang] = &loc
	}
	if _, ok := r.byCode[defaultLang]; !ok {
		return nil, fmt.Errorf("default locale %q not loaded", defaultLang)
	}
	return r, nil
}

func (r *Registry) Get(code string) *Locale {
	if loc, ok := r.byCode[code]; ok {
		return loc
	}
	return r.byCode[r.defaultLang]
}

// prepareBio substitutes the placeholders that don't change between requests
// ({city}, {company_name}, {company_url}). What's left for RenderBio is the
// per-request {time} and {todays_date}.
func (l *Locale) prepareBio() {
	r := strings.NewReplacer(
		"{city}", l.City,
		"{company_name}", l.Company,
		"{company_url}", l.CompanyURL,
	)
	l.bioPrepared = make([]string, len(l.Bio))
	for i, p := range l.Bio {
		l.bioPrepared[i] = r.Replace(p)
	}
}

// RenderBio returns the locale's bio with the per-request placeholders filled.
// timeHTML is inserted as-is (caller is responsible for safety).
func (l *Locale) RenderBio(timeHTML template.HTML, todaysDate string) []template.HTML {
	r := strings.NewReplacer(
		"{time}", string(timeHTML),
		"{todays_date}", todaysDate,
	)
	out := make([]template.HTML, len(l.bioPrepared))
	for i, p := range l.bioPrepared {
		out[i] = template.HTML(r.Replace(p))
	}
	return out
}
