package main

// GetRussianTemplate возвращает оригинальный синий сайт-сканер ПВТ КР с интеграцией комплаенс-блоков
func GetRussianTemplate() string {
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
			.input-field { padding: 12px; width: 60%; border-radius: 4px; border: none; margin-right: 10px; font-size: 16px; }
			.nav-link { color: white; text-decoration: none; font-weight: bold; background: #10b981; padding: 10px 20px; border-radius: 4px; }
			.nav-panel { background: #1b263b; padding: 15px; border-radius: 6px; margin-bottom: 20px; }
		</style>
	</head>
	<body>
		<div class="container">
			<div class="nav-panel">
				<span>New AI compliance modules are available for commercial banks of the KR: </span>
				<a href="/compliance" class="nav-link">Launch AI Compliance (up to 4th generation) →</a>
			</div>

			<div class="profile-box">
				📌 BUSINESS PROFILE OF HTP KR RESIDENT Kvantum Safe LLC: Comprehensive development of specialized software, design of intelligent AI assistants and autonomous AI agents, end-to-end regulatory and technological compliance control, as well as specialized legal assistance for financial institutions, distributed IT systems, and Web3 infrastructure.
			</div>
			
			<h1>🔍 Network Perimeter Express Audit (External AI Scanner)</h1>
			
			<form action="/report/mbank" method="GET" style="margin-top: 20px;">
				<input type="text" name="domain" class="input-field" placeholder="Enter bank domain or IP for scanning...">
				<button type="submit" class="btn">🚀 Run Scan</button>
			</form>
			
			<div style="margin-top: 30px;">
				<button type="button" class="btn" style="background: #f59e0b; color: #0b1329; margin-right: 10px;">🏆 Software Specification & Financial Benefit</button>
				<button type="button" class="btn" style="background: #38bdf8; color: #0b1329;">🏪 AI Solutions Marketplace</button>
			</div>
		</div>
	</body>
	</html>
	`
}

// GetComplianceTemplate возвращает новую финтех-страницу для банков КР
func GetComplianceTemplate() string {
	return `
	<!DOCTYPE html>
	<html lang="ru">
	<head>
		<meta charset="UTF-8">
		<title>Quantum Safe AI | Финтех-Комплаенс и Анти-Фрод</title>
		<style>
			body { font-family: 'Segoe UI', sans-serif; background: #f4f6f9; padding: 40px; color: #333; line-height: 1.6; }
			.container { max-width: 900px; margin: 0 auto; background: white; padding: 40px; border-radius: 12px; box-shadow: 0 4px 15px rgba(0,0,0,0.05); }
			.lang-switch { float: right; background: #1e3a8a; color: white; padding: 8px 15px; border-radius: 6px; text-decoration: none; font-weight: bold; }
			.btn-back { float: left; background: #475569; color: white; padding: 8px 15px; border-radius: 6px; text-decoration: none; font-weight: bold; }
			h1 { color: #1e3a8a; border-bottom: 2px solid #e5e7eb; padding-bottom: 10px; margin-top: 50px; }
			.demo-box { background: #f8fafc; border: 2px dashed #1e3a8a; padding: 25px; border-radius: 8px; margin-top: 30px; }
			.btn { display: inline-block; background: #10b981; color: white; padding: 12px 24px; border-radius: 6px; text-decoration: none; font-weight: bold; margin-right: 12px; margin-top: 12px; }
			.badge { background: #e0f2fe; color: #0369a1; padding: 4px 8px; border-radius: 4px; font-size: 14px; font-weight: bold; }
		</style>
	</head>
	<body>
		<div class="container">
			<a href="/" class="btn-back">← Вернуться к ИИ-Сканеру</a>
			<a href="/en" class="lang-switch">🌐 EN Version (Grant)</a>
			<h1>🛡️ Quantum Safe AI | Финтех-Комплаенс Контур</h1>
			<p><strong>Автономная ИИ-инфраструктура комплаенса для коммерческих банков и P2P-платформ Central Asia.</strong></p>
			<p><span class="badge">Резидент ПВТ</span> Официальный резидент Парка Высоких Технологий КР.</p>
			
			<h3>🚀 Модули Квантум-Контура:</h3>
			<ul>
				<li><strong>Quantum AI P2P Referee:</strong> Локальный движок компьютерного зрения (OCR), который за 2 секунды распознает чеки, извлекая ФИО плательщика капсом.</li>
				<li><strong>Банковский Комплаенс-Шлюз:</strong> Автоматический построитель графа скрытых взаимосвязей до 4-го колена по базам МинЮста КР и санкционным спискам.</li>
			</ul>

			<div class="demo-box">
				<h3>🤖 ИНТЕРАКТИВНЫЙ СТЕНД (LIVE DEMO В ТЕЛЕГРАМ):</h3>
				<a href="https://telegram.me" class="btn" target="_blank">⚡ Запустить ИИ-Судью Чеков</a>
				<a href="https://t.me" class="btn" style="background:#1e3a8a;" target="_blank">🏛️ Запустить Банковский Шлюз</a>
			</div>
		</div>
	</body>
	</html>`
}

// --- АНГЛИЙСКИЙ ШАБЛОН ГЛАВНОЙ СТРАНИЦЫ ---
func GetEnglishTemplate() string {
	return `
	<!DOCTYPE html>
	<html lang="en">
	<head>
		<meta charset="UTF-8">
		<title>Quantum Safe - AI Scanner</title>
		<style>
			body { font-family: sans-serif; background-color: #0d1b2a; color: #ffffff; padding: 40px; }
			.lang-panel { position: absolute; top: 20px; right: 20px; }
			.lang-panel a { color: #00b4d8; text-decoration: none; margin-left: 10px; font-weight: bold; }
			.container { max-width: 800px; margin: 0 auto; background: #1b263b; padding: 30px; border-radius: 8px; }
			.btn { background: #00b4d8; color: white; padding: 12px 24px; border: none; border-radius: 4px; cursor: pointer; font-size: 16px; }
		</style>
	</head>
	<body>
		<div class="lang-panel">
			<a href="/?lang=ru">🇷🇺 RU</a> | <a href="/?lang=en" style="border-bottom: 2px solid #00b4d8;">🇬🇧 EN</a>
		</div>
		<div class="container">
			<p><em>📌 BUSINESS PROFILE OF HTP KR RESIDENT Kvantum Safe LLC: Comprehensive development of specialized software, design of intelligent AI assistants and autonomous AI agents, end-to-end regulatory and technological compliance control, as well as specialized legal assistance for financial institutions, distributed IT systems, and Web3 infrastructure.</em></p>
			<h1>🔍 Network Perimeter Express Audit (External AI Scanner)</h1>
			<form action="/report/mbank" method="GET">
				<input type="hidden" name="lang" value="en">
				<input type="text" name="domain" placeholder="Enter bank domain (e.g., bakai.kg)" style="padding: 10px; width: 60%; border-radius: 4px; border: none;">
				<button type="submit" class="btn">🚀 Run Scan</button>
			</form>
		</div>
	</body>
	</html>
	`
}

// --- АНГЛИЙСКИЙ ШАБЛОН ВКЛАДКИ КОМПЛАЕНСА ---
func GetEnglishComplianceTemplate() string {
	return `
	<!DOCTYPE html>
	<html lang="en">
	<head>
		<meta charset="UTF-8">
		<title>Quantum Safe - Compliance</title>
		<style>
			body { font-family: sans-serif; background-color: #0d1b2a; color: #ffffff; padding: 40px; }
			.lang-panel { position: absolute; top: 20px; right: 20px; }
			.lang-panel a { color: #00b4d8; text-decoration: none; margin-left: 10px; font-weight: bold; }
			.container { max-width: 800px; margin: 0 auto; background: #1b263b; padding: 30px; border-radius: 8px; }
		</style>
	</head>
	<body>
		<div class="lang-panel">
			<a href="/compliance?lang=ru">🇷🇺 RU</a> | <a href="/compliance?lang=en" style="border-bottom: 2px solid #00b4d8;">🇬🇧 EN</a>
		</div>
		<div class="container">
			<h1>⚖️ AI Compliance Orchestrator</h1>
			<p>Interactive demonstration of the AI Receipt Judge and deep anti-money laundering screening up to the 4th generation of kinship (Ministry of Justice of the Kyrgyz Republic integration).</p>
			<p>Status: Ready for infrastructure deployment.</p>
		</div>
	</body>
	</html>
	`
}

// --- РУССКИЙ ШАБЛОН ВКЛАДКИ КОМПЛАЕНСА ---
func GetRussianComplianceTemplate() string {
	return `
	<!DOCTYPE html>
	<html lang="ru">
	<head>
		<meta charset="UTF-8">
		<title>Квантум Сейф - Комплаенс</title>
		<style>
			body { font-family: sans-serif; background-color: #0d1b2a; color: #ffffff; padding: 40px; }
			.lang-panel { position: absolute; top: 20px; right: 20px; }
			.lang-panel a { color: #00b4d8; text-decoration: none; margin-left: 10px; font-weight: bold; }
			.container { max-width: 800px; margin: 0 auto; background: #1b263b; padding: 30px; border-radius: 8px; }
		</style>
	</head>
	<body>
		<div class="lang-panel">
			<a href="/compliance?lang=ru" style="border-bottom: 2px solid #00b4d8;">🇷🇺 RU</a> | <a href="/compliance?lang=en">🇬🇧 EN</a>
		</div>
		<div class="container">
			<h1>⚖️ ИИ-Оркестратор Комплаенса</h1>
			<p>Интерактивная презентация ИИ-Судьи чеков и проверки до 4 колена Минюста КР.</p>
			<p>Статус: Готово к развертыванию инфраструктуры.</p>
		</div>
	</body>
	</html>
	`
}
