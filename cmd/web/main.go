package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"html/template"
	"io/fs"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	magzhanme "magzhan.me"
	"magzhan.me/internal/handlers"
	"magzhan.me/internal/i18n"
	"magzhan.me/internal/site"
)

func main() {
	addr := flag.String("addr", "127.0.0.1:8080", "address to listen on")
	flag.Parse()

	tz, err := time.LoadLocation(site.Me.Timezone)
	if err != nil {
		log.Fatalf("bad timezone %q: %v", site.Me.Timezone, err)
	}

	locales, err := i18n.Load(magzhanme.Locales, "locales", site.Me.DefaultLang)
	if err != nil {
		log.Fatalf("load locales: %v", err)
	}

	tmpls, err := parseTemplates(magzhanme.Templates)
	if err != nil {
		log.Fatalf("parse templates: %v", err)
	}

	srv := &handlers.Server{
		Tmpl:    tmpls,
		Locales: locales,
		Profile: site.Me,
		TZ:      tz,
	}

	mux := http.NewServeMux()

	staticFS, err := fs.Sub(magzhanme.Static, "static")
	if err != nil {
		log.Fatalf("static sub-fs: %v", err)
	}
	mux.Handle("GET /static/", cacheStatic(http.StripPrefix("/static/", http.FileServer(http.FS(staticFS)))))

	for _, lang := range site.Me.Languages {
		mux.Handle("GET "+lang.Path+"{$}", srv.Home(lang.Code))
	}

	mux.HandleFunc("GET cv.magzhan.me/{$}", srv.CVRedirect)
	mux.Handle("GET cv.magzhan.me/pdf", srv.CVPDF())
	mux.Handle("GET cv.magzhan.me/text", srv.CVText())

	mux.HandleFunc("GET /api/time", srv.Time)
	mux.HandleFunc("GET /", srv.NotFound)

	server := &http.Server{
		Addr:              *addr,
		Handler:           logging(mux),
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      15 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	go func() {
		log.Printf("listening on http://%s", *addr)
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("listen: %v", err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop
	log.Println("shutting down")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Printf("shutdown: %v", err)
	}
}

func logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		log.Printf("%s %s %s", r.Method, r.URL.Path, time.Since(start))
	})
}

func parseTemplates(fsys fs.FS) (map[string]*template.Template, error) {
	const layout = "templates/layout.html"
	out := map[string]*template.Template{}
	for _, page := range []string{"home", "404", "cv_pdf", "cv_text"} {
		t, err := template.ParseFS(fsys, layout, "templates/"+page+".html")
		if err != nil {
			return nil, fmt.Errorf("parse %s: %w", page, err)
		}
		out[page] = t
	}
	return out, nil
}

func cacheStatic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.URL.Path, "/static/fonts/") || strings.HasPrefix(r.URL.Path, "/static/vendor/") {
			w.Header().Set("Cache-Control", "public, max-age=31536000, immutable")
		} else {
			w.Header().Set("Cache-Control", "public, max-age=3600")
		}
		next.ServeHTTP(w, r)
	})
}
