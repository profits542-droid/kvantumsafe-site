package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"kvantum-safe-platform/internal/handlers"
	"kvantum-safe-platform/internal/store"
)

func main() {
	log.SetFlags(log.Ldate | log.Ltime | log.LUTC)
	log.SetPrefix("kvantum ")

	root, err := findRoot()
	if err != nil {
		log.Fatal(err)
	}

	dbPath := filepath.Join(root, "data", "kvantum.db")
	st, err := store.Open(dbPath)
	if err != nil {
		log.Fatalf("база данных: %v", err)
	}
	defer st.Close()

	if err := ensureAdmin(st); err != nil {
		log.Fatalf("администратор: %v", err)
	}

	app, err := handlers.New(st, root)
	if err != nil {
		log.Fatal(err)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /{$}", app.Home)
	mux.HandleFunc("POST /order", app.OrderAgent)
	mux.HandleFunc("GET /about", app.About)
	mux.HandleFunc("GET /marketplace", app.Marketplace)
	mux.HandleFunc("GET /studio", app.Studio)
	mux.HandleFunc("POST /studio/export", app.ExportDataset)
	mux.HandleFunc("POST /studio/codegen", app.GenerateConnector)
	mux.HandleFunc("GET /studio/download", app.DownloadDataset)
	mux.HandleFunc("GET /report", app.Report)
	mux.HandleFunc("GET /scanner", app.ScannerHome)
	mux.HandleFunc("GET /scanner/report", app.ScannerReport)
	mux.HandleFunc("GET /scanner/download", app.ScannerDownload)
	mux.HandleFunc("GET /scanner/billing", app.ScannerBilling)
	mux.HandleFunc("GET /scanner/api-docs", app.ScannerAPIDocs)
	mux.HandleFunc("GET /scanner/api/scan", app.ScannerAPI)
	mux.HandleFunc("GET /admin/login", app.LoginForm)
	mux.HandleFunc("POST /admin/login", app.Login)
	mux.HandleFunc("POST /admin/logout", app.Logout)
	mux.HandleFunc("GET /admin/dashboard", app.Dashboard)
	mux.HandleFunc("POST /admin/edit-content", app.EditContent)
	mux.HandleFunc("POST /admin/edit-product", app.EditProduct)
	mux.HandleFunc("GET /admin", app.Admin)
	mux.Handle("GET /static/", http.StripPrefix("/static/", http.FileServer(http.Dir(filepath.Join(root, "ui", "static")))))
	mux.Handle("GET /docs/", docsHandler(filepath.Join(root, "ui", "static", "docs")))

	server := &http.Server{
		Addr:              ":8080",
		Handler:           withAccessLog(withSecurityHeaders(mux)),
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      20 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	log.Printf("сайт открыт: http://localhost:8080")
	log.Printf("английская версия: http://localhost:8080/?lang=en")
	if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatal(err)
	}
}

func findRoot() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir, nil
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return "", fmt.Errorf("не найден go.mod, запустите команду из папки проекта")
		}
		dir = parent
	}
}

const (
	adminUsername = "admin"
	adminPassword = "KvantumSafe2026!"
)

func ensureAdmin(st *store.Store) error {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	if err := st.SetPassword(ctx, adminUsername, adminPassword); err != nil {
		return err
	}
	log.Printf("админ-панель: http://localhost:8080/admin/login")
	log.Printf("логин: %s", adminUsername)
	return nil
}

func withSecurityHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "DENY")
		w.Header().Set("Referrer-Policy", "same-origin")
		w.Header().Set("Content-Security-Policy", "default-src 'self'; style-src 'self' 'unsafe-inline'; img-src 'self'; form-action 'self'; frame-ancestors 'none'; base-uri 'self'")
		next.ServeHTTP(w, r)
	})
}

type statusRecorder struct {
	http.ResponseWriter
	status int
}

func (s *statusRecorder) WriteHeader(code int) {
	s.status = code
	s.ResponseWriter.WriteHeader(code)
}

func docsHandler(dir string) http.Handler {
	allowed := map[string]bool{
		"htp_certificate.pdf":    true,
		"htp_certificate_en.pdf": true,
		"patent_1259.pdf":        true,
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		name := strings.TrimPrefix(r.URL.Path, "/docs/")
		if !allowed[name] || strings.Contains(name, "/") {
			http.NotFound(w, r)
			return
		}
		if r.URL.Query().Get("download") == "1" {
			w.Header().Set("Content-Disposition", "attachment; filename=\""+name+"\"")
		}
		http.ServeFile(w, r, filepath.Join(dir, name))
	})
}

func withAccessLog(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rec := &statusRecorder{ResponseWriter: w, status: http.StatusOK}
		start := time.Now()
		next.ServeHTTP(rec, r)
		path := r.URL.RequestURI()
		if strings.HasPrefix(path, "/static/") {
			return
		}
		log.Printf("%s %s -> %d (%s)", r.Method, path, rec.status, time.Since(start).Round(time.Millisecond))
	})
}
