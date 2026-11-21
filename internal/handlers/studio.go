package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"time"

	"kvantum-safe-platform/internal/studio"
)

var exportNamePattern = regexp.MustCompile(`^kvantum-dataset-\d{8}-\d{6}\.jsonl$`)

func (a *App) Studio(w http.ResponseWriter, r *http.Request) {
	data, err := a.studioView(r, "", "", "ollama", "grok-3", "")
	if err != nil {
		log.Printf("студия: %v", err)
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	data.CSRF = issueCSRF(w)
	a.render(w, "studio", data)
}

func (a *App) ExportDataset(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, 8<<10)
	if err := r.ParseForm(); err != nil || !validCSRF(r) {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	ctx, cancel := requestContext(r)
	defer cancel()
	rows, err := a.store.ListProtocols(ctx)
	if err != nil {
		log.Printf("экспорт: %v", err)
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	if err := os.MkdirAll(a.exportDir, 0o755); err != nil {
		log.Printf("экспорт: %v", err)
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	name := "kvantum-dataset-" + time.Now().UTC().Format("20060102-150405") + ".jsonl"
	path := filepath.Join(a.exportDir, name)
	file, err := os.OpenFile(path, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0o600)
	if err != nil {
		log.Printf("экспорт: %v", err)
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	enc := json.NewEncoder(file)
	for _, row := range rows {
		record := map[string]string{
			"prompt":             studio.Mask(row.Prompt),
			"attack_scenario":    studio.Mask(row.AttackScenario),
			"response":           studio.Mask(row.Response),
			"system_reaction":    studio.Mask(row.SystemReaction),
			"risk_level":         row.RiskLevel,
			"remediation_advice": studio.Mask(row.RemediationAdvice),
			"category":           row.Category,
		}
		if err := enc.Encode(record); err != nil {
			file.Close()
			log.Printf("экспорт: %v", err)
			http.Error(w, "internal error", http.StatusInternalServerError)
			return
		}
	}
	if err := file.Close(); err != nil {
		log.Printf("экспорт: %v", err)
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	data, err := a.studioView(r, name, "", "ollama", "grok-3", "")
	if err != nil {
		log.Printf("студия: %v", err)
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	data.CSRF = issueCSRF(w)
	a.render(w, "studio", data)
}

func (a *App) DownloadDataset(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("file")
	if !exportNamePattern.MatchString(name) {
		http.NotFound(w, r)
		return
	}
	path := filepath.Join(a.exportDir, name)
	w.Header().Set("Content-Type", "application/jsonl")
	w.Header().Set("Content-Disposition", "attachment; filename=\""+name+"\"")
	http.ServeFile(w, r, path)
}

func (a *App) GenerateConnector(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, 8<<10)
	if err := r.ParseForm(); err != nil || !validCSRF(r) {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	mode := r.PostFormValue("engine")
	model := r.PostFormValue("model")
	if mode != "ollama" && mode != "xai" {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	code, err := studio.Generate(mode, model)
	if err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	notice := ""
	if r.PostFormValue("api_key") != "" {
		notice = "key"
	}
	data, err := a.studioView(r, "", code, mode, model, notice)
	if err != nil {
		log.Printf("студия: %v", err)
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	data.CSRF = issueCSRF(w)
	a.render(w, "studio", data)
}

func (a *App) studioView(r *http.Request, exportName, generated, engine, model, notice string) (viewData, error) {
	data := a.baseView(r)
	ctx, cancel := requestContext(r)
	defer cancel()
	rows, err := a.store.ListProtocols(ctx)
	if err != nil {
		return data, err
	}
	data.Protocols = rows
	data.ExportName = exportName
	data.Generated = generated
	if engine == "" {
		engine = "ollama"
	}
	if model == "" {
		model = "grok-3"
	}
	data.Engine = engine
	data.Model = model
	data.StudioNotice = notice
	return data, nil
}
