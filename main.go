package main

import (
	"log"
	"net/http"
)

func main() {
	// 1. ГЛАВНАЯ СТРАНИЦА: Наш оригинальный ИИ-Сканер сетевого периметра ПВТ КР
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		
		// Если в URL передан параметр ?lang=en — отдаем чистый английский HTML напрямую
		if r.URL.Query().Get("lang") == "en" {
			w.Write([]byte(GetHardcodedEnglishTemplate()))
		} else {
			// Во всех остальных случаях отдаем оригинальный русский шаблон из views.go
			w.Write([]byte(GetRussianTemplate()))
		}
	})

	// 2. ВКЛАДКА КОМПЛАЕНСА
	http.HandleFunc("/compliance", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		
		if r.URL.Query().Get("lang") == "en" {
			w.Write([]byte(GetHardcodedEnglishComplianceTemplate()))
		} else {
			w.Write([]byte(GetRussianComplianceTemplate()))
		}
	})

	log.Println("[SUCCESS] Квантум Сейф Веб-Ядро запущено на порту: 8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("[ERROR] Не удалось запустить сервер: %v", err)
	}
}

// --- ЖЕСТКО ПРОПИСАННЫЙ АНГЛИЙСКИЙ ШАБЛОН ГЛАВНОЙ СТРАНИЦЫ ---
func GetHardcodedEnglishTemplate() string {
	return `
	<!DOCTYPE html>
	<html lang="en">
	<head>
		<meta charset="UTF-8">
		<title>Quantum Safe - AI Scanner</title>
		<style>
			body { font-family: sans-serif; background-color: #0d1b2a; color: #ffffff; padding: 40px; }
			.container { max-width: 800px; margin: 0 auto; background: #1b263b; padding: 30px; border-radius: 8px; }
			.btn { background: #00b4d8; color: white; padding: 12px 24px; border: none; border-radius: 4px; cursor: pointer; font-size: 16px; font-weight: bold; }
			.profile-box { background: #14213d; padding: 20px; border-radius: 6px; border-left: 4px solid #00b4d8; margin-bottom: 30px; font-style: italic; }
			.input-field { padding: 12px; width: 60%; border-radius: 4px; border: none; margin-right: 10px; font-size: 16px; background: #0d1b2a; color: white; border: 1px solid #00b4d8; }
			.nav-link { color: white; text-decoration: none; font-weight: bold; background: #10b981; padding: 10px 20px; border-radius: 4px; }
			.nav-panel { background: #14213d; padding: 15px; border-radius: 6px; margin-bottom: 20px; display: flex; justify-content: space-between; align-items: center; }
		</style>
	</head>
	<body>
		<div class="container">
			<div class="nav-panel">
				<span>New AI compliance modules are available for commercial banks of the KR:</span>
				<a href="/compliance?lang=en" class="nav-link">Launch AI Compliance (up to 4th generation) →</a>
			</div>

			<div class="profile-box">
				📌 BUSINESS PROFILE OF HTP KR RESIDENT Kvantum Safe LLC: Comprehensive development of specialized software, design of intelligent AI assistants and autonomous AI agents, end-to-end regulatory and technological compliance control, as well as specialized legal assistance for financial institutions, distributed IT systems, and Web3 infrastructure.
			</div>
			
			<h1>🔍 Network Perimeter Express Audit (External AI Scanner)</h1>
			
			<form action="/report/mbank" method="GET" style="margin-top: 20px;">
				<input type="hidden" name="lang" value="en">
				<input type="text" name="domain" class="input-field" placeholder="Enter bank domain or IP for scanning (e.g., bakai.kg)...">
				<button type="submit" class="btn">🚀 Run Scan</button>
			</form>
		</div>
	</body>
	</html>
	`
}

// --- ЖЕСТКО ПРОПИСАННЫЙ АНГЛИЙСКИЙ ШАБЛОН КОМПЛАЕНСА ---
func GetHardcodedEnglishComplianceTemplate() string {
	return `
	<!DOCTYPE html>
	<html lang="en">
	<head>
		<meta charset="UTF-8">
		<title>Quantum Safe - Compliance</title>
		<style>
			body { font-family: sans-serif; background-color: #0d1b2a; color: #ffffff; padding: 40px; }
			.container { max-width: 800px; margin: 0 auto; background: #1b263b; padding: 30px; border-radius: 8px; }
			.nav-link { color: white; text-decoration: none; font-weight: bold; background: #00b4d8; padding: 10px 20px; border-radius: 4px; }
		</style>
	</head>
	<body>
		<div class="container">
			<p><a href="/?lang=en" class="nav-link">← Back to AI Scanner</a></p>
			<h1 style="margin-top: 30px;">⚖️ AI Compliance Orchestrator</h1>
			<p style="font-size: 18px; line-height: 1.6;">Interactive demonstration of the AI Receipt Judge and deep anti-money laundering screening up to the 4th generation of kinship (Ministry of Justice of the Kyrgyz Republic integration).</p>
			<p style="color: #10b981; font-weight: bold;">Status: Ready for infrastructure deployment.</p>
		</div>
	</body>
	</html>
	`
}
