package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestScanQuotaLimitsAndDropsStaleBuckets(t *testing.T) {
	q := newScanQuota()
	for i := 0; i < publicScanPerMinute; i++ {
		if !q.allow("ip:a", publicScanPerMinute) {
			t.Fatalf("public request %d was rejected", i+1)
		}
	}
	if q.allow("ip:a", publicScanPerMinute) {
		t.Fatal("public request above the limit was allowed")
	}
	if !q.allow("ip:b", publicScanPerMinute) {
		t.Fatal("a different address was rejected")
	}

	q.mu.Lock()
	q.buckets["old"] = []time.Time{time.Now().Add(-2 * time.Minute)}
	q.mu.Unlock()
	if !q.allow("ip:c", publicScanPerMinute) {
		t.Fatal("request after cleanup was rejected")
	}
	q.mu.Lock()
	_, stale := q.buckets["old"]
	q.mu.Unlock()
	if stale {
		t.Fatal("stale bucket was kept")
	}
}

func TestEnterpriseScannerKey(t *testing.T) {
	t.Setenv("SCANNER_API_KEYS", " alpha , beta ")
	header := httptest.NewRequest(http.MethodGet, "/scanner/api/scan", nil)
	header.Header.Set("X-API-Key", "alpha")
	if key, ok := enterpriseScannerKey(header); !ok || key != "alpha" {
		t.Fatalf("header key = %q ok=%v", key, ok)
	}

	query := httptest.NewRequest(http.MethodGet, "/scanner/api/scan?api_key=beta", nil)
	if key, ok := enterpriseScannerKey(query); !ok || key != "beta" {
		t.Fatalf("query key = %q ok=%v", key, ok)
	}

	missing := httptest.NewRequest(http.MethodGet, "/scanner/api/scan", nil)
	if _, ok := enterpriseScannerKey(missing); ok {
		t.Fatal("empty request was treated as an enterprise key")
	}
	wrong := httptest.NewRequest(http.MethodGet, "/scanner/api/scan?api_key=nope", nil)
	if _, ok := enterpriseScannerKey(wrong); ok {
		t.Fatal("unknown key was accepted")
	}
}
