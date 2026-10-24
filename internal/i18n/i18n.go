package i18n

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"path"
	"strings"
)

type Locale struct {
	Lang     string    `json:"lang"`
	LangTag  string    `json:"lang_tag"`
	Title    string    `json:"title"`
	Tagline  string    `json:"tagline"`
	NavLinks string    `json:"nav_links"`
	CV       CVStrings `json:"cv"`
}

type CVStrings struct {
	PageTitle       string `json:"page_title"`
	TabText         string `json:"tab_text"`
	TabPDF          string `json:"tab_pdf"`
	Download        string `json:"download"`
	FallbackOpen    string `json:"fallback_open"`
	TextPlaceholder string `json:"text_placeholder"`
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
