package handlers

import (
	"context"
	"crypto/rand"
	"crypto/subtle"
	"database/sql"
	"encoding/hex"
	"errors"
	"html/template"
	"log"
	"net"
	"net/http"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"golang.org/x/crypto/bcrypt"

	"kvantum-safe-platform/internal/store"
)

const (
	sessionCookie = "kvantum_session"
	csrfCookie    = "kvantum_csrf"
	sessionTTL    = 12 * time.Hour
)

type App struct {
	store     *store.Store
	pages     map[string]*template.Template
	hits      *rateGate
	dummyHash []byte
}

type viewData struct {
	Lang      string
	Year      int
	CSRF      string
	Domain    string
	Sent      bool
	Saved     bool
	ErrorCode string
	Username  string
	Path      string
	Checks    []checkRow
	Text      map[string]string
	Blocks    []store.Block
	Products  []store.Product
}

type checkRow struct {
	Domain string
	When   string
}

func New(st *store.Store, root string) (*App, error) {
	pages := map[string]string{
		"home":        "home.html",
		"about":       "about.html",
		"marketplace": "marketplace.html",
		"report":      "report.html",
		"login":       "login.html",
		"admin":       "admin.html",
	}
	parsed := make(map[string]*template.Template, len(pages))
	base := filepath.Join(root, "ui", "html", "base.html")
	for name, file := range pages {
		t, err := template.ParseFiles(base, filepath.Join(root, "ui", "html", "pages", file))
		if err != nil {
			return nil, err
		}
		parsed[name] = t
	}
	dashboard, err := template.ParseFiles(base, filepath.Join(root, "ui", "html", "admin", "dashboard.html"))
	if err != nil {
		return nil, err
	}
	parsed["dashboard"] = dashboard
	dummy, err := bcrypt.GenerateFromPassword([]byte("not-a-real-password"), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	return &App{
		store:     st,
		pages:     parsed,
		hits:      newRateGate(8, time.Minute),
		dummyHash: dummy,
	}, nil
}

func (a *App) Home(w http.ResponseWriter, r *http.Request) {
	a.renderPublic(w, r, "home")
}

func (a *App) About(w http.ResponseWriter, r *http.Request) {
	a.renderPublic(w, r, "about")
}

func (a *App) Marketplace(w http.ResponseWriter, r *http.Request) {
	a.renderPublic(w, r, "marketplace")
}

func (a *App) renderPublic(w http.ResponseWriter, r *http.Request, page string) {
	data, err := a.publicView(r, requestLang(r), "", "")
	if err != nil {
		log.Printf("%s: %v", page, err)
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	a.render(w, page, data)
}

func (a *App) Report(w http.ResponseWriter, r *http.Request) {
	if !a.hits.allow(clientIP(r)) {
		http.Error(w, "too many requests", http.StatusTooManyRequests)
		return
	}
	lang := requestLang(r)
	raw := r.URL.Query().Get("domain")
	domain, ok := normalizeDomain(raw)
	if !ok {
		data, err := a.publicView(r, lang, strings.TrimSpace(raw), "invalid_domain")
		if err != nil {
			log.Printf("главная: %v", err)
			http.Error(w, "internal error", http.StatusInternalServerError)
			return
		}
		a.render(w, "home", data)
		return
	}

	ctx, cancel := requestContext(r)
	defer cancel()
	if err := a.store.AddDomainCheck(ctx, domain); err != nil {
		log.Printf("сохранение домена: %v", err)
	}

	data := a.baseView(r)
	data.Lang = lang
	data.Domain = domain
	a.render(w, "report", data)
}

func (a *App) LoginForm(w http.ResponseWriter, r *http.Request) {
	if _, ok := a.currentUser(r); ok {
		http.Redirect(w, r, adminPath(requestLang(r)), http.StatusSeeOther)
		return
	}
	data := a.baseView(r)
	data.CSRF = issueCSRF(w)
	a.render(w, "login", data)
}

func (a *App) Login(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, 8<<10)
	if err := r.ParseForm(); err != nil || !validCSRF(r) {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	lang := requestLang(r)
	username := strings.TrimSpace(r.PostFormValue("username"))
	password := r.PostFormValue("password")

	ctx, cancel := requestContext(r)
	defer cancel()
	userID, err := a.store.Authenticate(ctx, username, password, a.dummyHash)
	if err != nil {
		if !a.store.InvalidLogin(err) {
			log.Printf("вход: %v", err)
			http.Error(w, "internal error", http.StatusInternalServerError)
			return
		}
		data := a.baseView(r)
		data.Lang = lang
		data.ErrorCode = "invalid_login"
		data.CSRF = issueCSRF(w)
		a.render(w, "login", data)
		return
	}

	token, err := newToken()
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	expires := time.Now().Add(sessionTTL)
	if err := a.store.CreateSession(ctx, userID, token, expires); err != nil {
		log.Printf("сессия: %v", err)
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	http.SetCookie(w, &http.Cookie{
		Name:     sessionCookie,
		Value:    token,
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		Expires:  expires,
		MaxAge:   int(sessionTTL.Seconds()),
	})
	http.Redirect(w, r, adminPath(lang), http.StatusSeeOther)
}

func (a *App) Logout(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil || !validCSRF(r) {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	if cookie, err := r.Cookie(sessionCookie); err == nil && cookie.Value != "" {
		ctx, cancel := requestContext(r)
		defer cancel()
		_ = a.store.DeleteSession(ctx, cookie.Value)
	}
	http.SetCookie(w, &http.Cookie{
		Name:     sessionCookie,
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   -1,
	})
	http.Redirect(w, r, "/admin/login?lang="+requestLang(r), http.StatusSeeOther)
}

func (a *App) Admin(w http.ResponseWriter, r *http.Request) {
	username, ok := a.currentUser(r)
	if !ok {
		http.Redirect(w, r, "/admin/login?lang="+requestLang(r), http.StatusSeeOther)
		return
	}
	ctx, cancel := requestContext(r)
	defer cancel()
	items, err := a.store.RecentDomainChecks(ctx, 50)
	if err != nil {
		log.Printf("список доменов: %v", err)
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	rows := make([]checkRow, 0, len(items))
	for _, item := range items {
		rows = append(rows, checkRow{
			Domain: item.Domain,
			When:   item.CreatedAt.UTC().Format("02.01.2006 15:04 UTC"),
		})
	}
	data := a.baseView(r)
	data.Username = username
	data.Checks = rows
	data.CSRF = issueCSRF(w)
	a.render(w, "admin", data)
}

func (a *App) currentUser(r *http.Request) (string, bool) {
	cookie, err := r.Cookie(sessionCookie)
	if err != nil || cookie.Value == "" {
		return "", false
	}
	ctx, cancel := requestContext(r)
	defer cancel()
	username, err := a.store.UsernameBySession(ctx, cookie.Value)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			log.Printf("проверка сессии: %v", err)
		}
		return "", false
	}
	return username, true
}

func (a *App) baseView(r *http.Request) viewData {
	return viewData{Lang: requestLang(r), Year: time.Now().Year(), Path: r.URL.Path}
}

func (a *App) publicView(r *http.Request, lang, domain, errorCode string) (viewData, error) {
	data := a.baseView(r)
	data.Lang = lang
	data.Domain = domain
	data.ErrorCode = errorCode
	ctx, cancel := requestContext(r)
	defer cancel()
	text, err := a.store.ContentByLang(ctx, lang)
	if err != nil {
		return data, err
	}
	products, err := a.store.ListProducts(ctx)
	if err != nil {
		return data, err
	}
	data.Text = text
	data.Products = products
	return data, nil
}

func (a *App) render(w http.ResponseWriter, page string, data viewData) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Cache-Control", "no-store")
	if err := a.pages[page].ExecuteTemplate(w, "base", data); err != nil {
		log.Printf("шаблон %s: %v", page, err)
	}
}

func requestLang(r *http.Request) string {
	if r.URL.Query().Get("lang") == "en" {
		return "en"
	}
	return "ru"
}

func adminPath(lang string) string {
	if lang == "en" {
		return "/admin/dashboard?lang=en"
	}
	return "/admin/dashboard"
}

func normalizeDomain(raw string) (string, bool) {
	raw = strings.TrimSpace(strings.ToLower(raw))
	raw = strings.TrimPrefix(raw, "https://")
	raw = strings.TrimPrefix(raw, "http://")
	if i := strings.IndexAny(raw, "/?#"); i >= 0 {
		raw = raw[:i]
	}
	raw = strings.TrimSuffix(raw, ".")
	if len(raw) < 4 || len(raw) > 253 || strings.Contains(raw, "..") {
		return "", false
	}
	labels := strings.Split(raw, ".")
	if len(labels) < 2 {
		return "", false
	}
	for _, label := range labels {
		if len(label) == 0 || len(label) > 63 {
			return "", false
		}
		if label[0] == '-' || label[len(label)-1] == '-' {
			return "", false
		}
		for _, ch := range label {
			if (ch < 'a' || ch > 'z') && (ch < '0' || ch > '9') && ch != '-' {
				return "", false
			}
		}
	}
	return raw, true
}

func issueCSRF(w http.ResponseWriter) string {
	token, err := newToken()
	if err != nil {
		return ""
	}
	http.SetCookie(w, &http.Cookie{
		Name:     csrfCookie,
		Value:    token,
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   int((2 * time.Hour).Seconds()),
	})
	return token
}

func validCSRF(r *http.Request) bool {
	cookie, err := r.Cookie(csrfCookie)
	form := r.PostFormValue("csrf")
	if err != nil || cookie.Value == "" || form == "" || len(cookie.Value) != len(form) {
		return false
	}
	return subtle.ConstantTimeCompare([]byte(cookie.Value), []byte(form)) == 1
}

func newToken() (string, error) {
	buf := make([]byte, 32)
	if _, err := rand.Read(buf); err != nil {
		return "", err
	}
	return hex.EncodeToString(buf), nil
}

func clientIP(r *http.Request) string {
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return host
}

func requestContext(r *http.Request) (context.Context, context.CancelFunc) {
	return context.WithTimeout(r.Context(), 5*time.Second)
}

type rateGate struct {
	mu    sync.Mutex
	limit int
	span  time.Duration
	hits  map[string][]time.Time
}

func newRateGate(limit int, span time.Duration) *rateGate {
	return &rateGate{limit: limit, span: span, hits: make(map[string][]time.Time)}
}

func (g *rateGate) allow(key string) bool {
	now := time.Now()
	g.mu.Lock()
	defer g.mu.Unlock()
	prev := g.hits[key]
	recent := make([]time.Time, 0, len(prev)+1)
	for _, ts := range prev {
		if now.Sub(ts) < g.span {
			recent = append(recent, ts)
		}
	}
	if len(recent) >= g.limit {
		g.hits[key] = recent
		return false
	}
	g.hits[key] = append(recent, now)
	return true
}
