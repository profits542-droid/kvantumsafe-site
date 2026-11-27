package handlers

import (
	"crypto/sha256"
	"crypto/subtle"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"kvantum-safe-platform/internal/scanner"
)

const (
	publicScanPerMinute     = 10
	enterpriseScanPerMinute = 1000
	scanQuotaWindow         = time.Minute
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
	if key, ok := enterpriseScannerKey(r); ok {
		if !a.scanAPI.allow("key:"+scannerKeyID(key), enterpriseScanPerMinute) {
			writeScanRateLimit(w)
			return
		}
	} else if !a.scanAPI.allow("ip:"+clientIP(r), publicScanPerMinute) {
		writeScanRateLimit(w)
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

type scanLimitBody struct {
	OK    bool   `json:"ok"`
	Error string `json:"error"`
}

func writeScanRateLimit(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusTooManyRequests)
	_ = json.NewEncoder(w).Encode(scanLimitBody{
		OK:    false,
		Error: "Rate limit exceeded. Use X-API-Key for enterprise access.",
	})
}

func enterpriseScannerKey(r *http.Request) (string, bool) {
	candidates := []string{
		strings.TrimSpace(r.Header.Get("X-API-Key")),
		strings.TrimSpace(r.URL.Query().Get("api_key")),
	}
	keys := scannerAPIKeys()
	var matched string
	for _, candidate := range candidates {
		if candidate == "" {
			continue
		}
		for _, allowed := range keys {
			if subtle.ConstantTimeCompare([]byte(candidate), []byte(allowed)) == 1 {
				matched = allowed
			}
		}
	}
	if matched == "" {
		return "", false
	}
	return matched, true
}

func scannerAPIKeys() []string {
	raw := os.Getenv("SCANNER_API_KEYS")
	parts := strings.Split(raw, ",")
	keys := make([]string, 0, len(parts))
	for _, part := range parts {
		key := strings.TrimSpace(part)
		if key != "" {
			keys = append(keys, key)
		}
	}
	return keys
}

func scannerKeyID(key string) string {
	sum := sha256.Sum256([]byte(key))
	return hex.EncodeToString(sum[:])
}

type scanQuota struct {
	mu      sync.Mutex
	buckets map[string][]time.Time
}

func newScanQuota() *scanQuota {
	return &scanQuota{buckets: make(map[string][]time.Time)}
}

func (q *scanQuota) allow(id string, limit int) bool {
	now := time.Now()
	q.mu.Lock()
	defer q.mu.Unlock()
	for key, stamps := range q.buckets {
		fresh := freshScanStamps(stamps, now)
		if len(fresh) == 0 {
			delete(q.buckets, key)
			continue
		}
		if len(fresh) != len(stamps) {
			q.buckets[key] = fresh
		}
	}
	recent := q.buckets[id]
	if len(recent) >= limit {
		return false
	}
	q.buckets[id] = append(recent, now)
	return true
}

func freshScanStamps(stamps []time.Time, now time.Time) []time.Time {
	fresh := make([]time.Time, 0, len(stamps))
	for _, ts := range stamps {
		if now.Sub(ts) < scanQuotaWindow {
			fresh = append(fresh, ts)
		}
	}
	return fresh
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
