package handlers

import (
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"kvantum-safe-platform/internal/store"
)

var contentKeyPattern = regexp.MustCompile(`^[a-z0-9_]{1,64}$`)

func (a *App) Dashboard(w http.ResponseWriter, r *http.Request) {
	username, ok := a.currentUser(r)
	if !ok {
		http.Redirect(w, r, "/admin/login?lang="+requestLang(r), http.StatusSeeOther)
		return
	}
	ctx, cancel := requestContext(r)
	defer cancel()

	blocks, err := a.store.ListContent(ctx)
	if err != nil {
		log.Printf("список контента: %v", err)
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	products, err := a.store.ListProducts(ctx)
	if err != nil {
		log.Printf("список продуктов: %v", err)
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	data := a.baseView(r)
	data.Username = username
	data.Blocks = blocks
	data.Products = products
	data.Saved = r.URL.Query().Get("saved") == "1"
	data.CSRF = issueCSRF(w)
	a.render(w, "dashboard", data)
}

func (a *App) EditContent(w http.ResponseWriter, r *http.Request) {
	if _, ok := a.currentUser(r); !ok {
		http.Redirect(w, r, "/admin/login?lang="+requestLang(r), http.StatusSeeOther)
		return
	}
	r.Body = http.MaxBytesReader(w, r.Body, 64<<10)
	if err := r.ParseForm(); err != nil || !validCSRF(r) {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	key := r.PostFormValue("key")
	lang := r.PostFormValue("lang")
	value := strings.TrimSpace(r.PostFormValue("value"))
	if !contentKeyPattern.MatchString(key) || (lang != "ru" && lang != "en") || value == "" || len(value) > 8000 {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	ctx, cancel := requestContext(r)
	defer cancel()
	if err := a.store.UpdateContent(ctx, key, lang, value); err != nil {
		log.Printf("сохранение блока: %v", err)
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, dashboardPath(r.PostFormValue("ui_lang"), true), http.StatusSeeOther)
}

func (a *App) EditProduct(w http.ResponseWriter, r *http.Request) {
	if _, ok := a.currentUser(r); !ok {
		http.Redirect(w, r, "/admin/login?lang="+requestLang(r), http.StatusSeeOther)
		return
	}
	r.Body = http.MaxBytesReader(w, r.Body, 64<<10)
	if err := r.ParseForm(); err != nil || !validCSRF(r) {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	id, err := strconv.ParseInt(r.PostFormValue("id"), 10, 64)
	if err != nil || id <= 0 {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	item := store.Product{
		ID:      id,
		NameRU:  strings.TrimSpace(r.PostFormValue("name_ru")),
		NameEN:  strings.TrimSpace(r.PostFormValue("name_en")),
		DescRU:  strings.TrimSpace(r.PostFormValue("desc_ru")),
		DescEN:  strings.TrimSpace(r.PostFormValue("desc_en")),
		PriceRU: strings.TrimSpace(r.PostFormValue("price_ru")),
		PriceEN: strings.TrimSpace(r.PostFormValue("price_en")),
		Status:  strings.TrimSpace(r.PostFormValue("status")),
	}
	item.Price = item.PriceEN
	if !validProduct(item) {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	ctx, cancel := requestContext(r)
	defer cancel()
	if err := a.store.UpdateProduct(ctx, item); err != nil {
		log.Printf("сохранение продукта: %v", err)
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, dashboardPath(r.PostFormValue("ui_lang"), true), http.StatusSeeOther)
}

func validProduct(item store.Product) bool {
	fields := []string{item.NameRU, item.NameEN, item.DescRU, item.DescEN, item.PriceRU, item.PriceEN, item.Status}
	for _, field := range fields {
		if field == "" || len(field) > 4000 {
			return false
		}
	}
	return len(item.Status) <= 32
}

func dashboardPath(lang string, saved bool) string {
	path := "/admin/dashboard"
	if lang == "en" {
		path += "?lang=en"
	} else {
		path += "?lang=ru"
	}
	if saved {
		path += "&saved=1"
	}
	return path
}
