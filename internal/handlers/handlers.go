package handlers

import (
	"html/template"
	"log"
	"net/http"

	"magzhan.me/internal/i18n"
	"magzhan.me/internal/site"
)

type Server struct {
	Tmpl    map[string]*template.Template
	Locales *i18n.Registry
	Profile site.Profile
}

type Page struct {
	Locale  *i18n.Locale
	Profile site.Profile
	Links   NavSection
}

type NavSection struct {
	Label string
	Items []site.Link
}

func (s *Server) page(lang string) Page {
	loc := s.Locales.Get(lang)
	return Page{
		Locale:  loc,
		Profile: s.Profile,
		Links:   NavSection{Label: loc.NavLinks, Items: s.Profile.Links},
	}
}

func (s *Server) render(w http.ResponseWriter, name string, status int, data any) {
	tmpl, ok := s.Tmpl[name]
	if !ok {
		log.Printf("render: no template %q", name)
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(status)
	if err := tmpl.ExecuteTemplate(w, "layout", data); err != nil {
		log.Printf("render %s: %v", name, err)
	}
}

func (s *Server) Home(lang string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		s.render(w, "home", http.StatusOK, s.page(lang))
	}
}

func (s *Server) NotFound(w http.ResponseWriter, r *http.Request) {
	loc := s.Locales.Get(s.Profile.DefaultLang)
	s.render(w, "404", http.StatusNotFound, Page{
		Locale:  loc,
		Profile: s.Profile,
	})
}

func (s *Server) cvPage() Page {
	return Page{
		Locale:  s.Locales.Get(s.Profile.DefaultLang),
		Profile: s.Profile,
	}
}

func (s *Server) CVRedirect(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/pdf", http.StatusFound)
}

func (s *Server) CVPDF() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		s.render(w, "cv_pdf", http.StatusOK, s.cvPage())
	}
}

func (s *Server) CVText() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		s.render(w, "cv_text", http.StatusOK, s.cvPage())
	}
}
