package handlers

import (
	"log"
	"net/http"
	"strings"
	"time"
	"unicode/utf8"
)

func (a *App) OrderAgent(w http.ResponseWriter, r *http.Request) {
	if !a.hits.allow(clientIP(r)) {
		http.Error(w, "too many requests", http.StatusTooManyRequests)
		return
	}
	r.Body = http.MaxBytesReader(w, r.Body, 16<<10)
	if err := r.ParseForm(); err != nil || !validCSRF(r) {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	lang := requestLang(r)
	company := strings.TrimSpace(r.PostFormValue("company"))
	contact := strings.TrimSpace(r.PostFormValue("contact"))
	task := strings.TrimSpace(r.PostFormValue("task"))
	if utf8.RuneCountInString(company) < 2 || utf8.RuneCountInString(company) > 120 ||
		utf8.RuneCountInString(contact) < 3 || utf8.RuneCountInString(contact) > 160 ||
		utf8.RuneCountInString(task) < 10 || utf8.RuneCountInString(task) > 2000 {
		data, err := a.publicView(r, lang, "", "invalid_order")
		if err != nil {
			log.Printf("заявка: %v", err)
			http.Error(w, "internal error", http.StatusInternalServerError)
			return
		}
		data.CSRF = issueCSRF(w)
		a.render(w, "home", data)
		return
	}
	ctx, cancel := requestContext(r)
	defer cancel()
	if err := a.store.SaveAgentRequest(ctx, company, contact, task, time.Now().UTC().Format(time.RFC3339)); err != nil {
		log.Printf("заявка: %v", err)
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/?lang="+lang+"&ordered=1#order", http.StatusSeeOther)
}
