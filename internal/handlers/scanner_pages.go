package handlers

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"kvantum-safe-platform/internal/scanner"
)

func (a *App) ScannerHome(w http.ResponseWriter, r *http.Request) {
	data := a.baseView(r)
	if r.URL.Query().Get("domain") != "" {
		data.Domain = strings.TrimSpace(r.URL.Query().Get("domain"))
	}
	if data.Domain == "" {
		data.Domain = "bakai.kg"
	}
	a.render(w, "scanner", data)
}

func (a *App) ScannerReport(w http.ResponseWriter, r *http.Request) {
	data, ok := a.scanOrHome(w, r)
	if !ok {
		return
	}
	a.render(w, "scanner_report", data)
}

func (a *App) ScannerDownload(w http.ResponseWriter, r *http.Request) {
	a.render(w, "scanner_download", a.baseView(r))
}

func (a *App) ScannerBilling(w http.ResponseWriter, r *http.Request) {
	a.render(w, "scanner_billing", a.baseView(r))
}

func (a *App) ScannerAPIDocs(w http.ResponseWriter, r *http.Request) {
	a.render(w, "scanner_api_docs", a.baseView(r))
}

func (a *App) ScannerAPI(w http.ResponseWriter, r *http.Request) {
	if !a.hits.allow(clientIP(r)) {
		http.Error(w, "too many requests", http.StatusTooManyRequests)
		return
	}
	lang := requestLang(r)
	raw := r.URL.Query().Get("domain")
	domain, ok := normalizeDomain(raw)
	if !ok {
		http.Error(w, "invalid domain", http.StatusBadRequest)
		return
	}
	result := scanner.ScanDomainTLS(domain, lang)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	_ = json.NewEncoder(w).Encode(map[string]any{
		"domain":       result.Domain,
		"ip":           result.IP,
		"tls_version":  result.TLSVersion,
		"cipher_suite": result.CipherSuite,
		"subject":      result.Subject,
		"issuer":       result.Issuer,
		"not_before":   formatScanTime(result.NotBefore),
		"not_after":    formatScanTime(result.NotAfter),
		"ok":           result.OK(),
		"error":        result.Error,
		"logs":         result.Logs,
	})
}

func (a *App) scanOrHome(w http.ResponseWriter, r *http.Request) (viewData, bool) {
	if !a.hits.allow(clientIP(r)) {
		http.Error(w, "too many requests", http.StatusTooManyRequests)
		return viewData{}, false
	}
	lang := requestLang(r)
	raw := strings.TrimSpace(r.URL.Query().Get("domain"))
	if raw == "" {
		raw = "bakai.kg"
	}
	domain, ok := normalizeDomain(raw)
	data := a.baseView(r)
	data.Lang = lang
	data.Domain = raw
	if !ok {
		data.ErrorCode = "invalid_domain"
		a.render(w, "scanner", data)
		return viewData{}, false
	}
	result := scanner.ScanDomainTLS(domain, lang)
	data.Domain = domain
	data.Scan = scanViewFrom(result, lang)
	return data, true
}

func scanViewFrom(result scanner.ScanResult, lang string) scanView {
	view := scanView{
		IP:        dashScan(result.IP),
		Version:   dashScan(result.TLSVersion),
		Cipher:    dashScan(result.CipherSuite),
		Subject:   dashScan(result.Subject),
		Issuer:    dashScan(result.Issuer),
		NotAfter:  "—",
		Validity:  "—",
		Error:     result.Error,
		Logs:      strings.Join(result.Logs, "\n"),
		Connected: result.OK(),
	}
	if view.Error == "" {
		if lang == "en" {
			view.Error = "none"
		} else {
			view.Error = "нет"
		}
	}
	if !result.NotAfter.IsZero() {
		view.NotAfter = result.NotAfter.UTC().Format("2006-01-02 15:04:05 MST")
		left := time.Until(result.NotAfter)
		if left > 0 {
			days := int(left.Hours() / 24)
			if lang == "en" {
				view.Validity = "valid for about " + itoa(days) + " days"
			} else {
				view.Validity = "действителен ещё ~" + itoa(days) + " дн."
			}
		} else if lang == "en" {
			view.Validity = "expired"
		} else {
			view.Validity = "срок истёк"
		}
	}
	return view
}

func formatScanTime(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.UTC().Format(time.RFC3339)
}

func dashScan(s string) string {
	if s == "" {
		return "—"
	}
	return s
}

func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	var buf [12]byte
	i := len(buf)
	for n > 0 {
		i--
		buf[i] = byte('0' + n%10)
		n /= 10
	}
	return string(buf[i:])
}
