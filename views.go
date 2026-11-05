package main

// GetRussianTemplate возвращает оригинальный синий сайт-сканер ПВТ КР с интеграцией комплаенс-блоков
func GetRussianTemplate() string {
	return `
	<!DOCTYPE html>
	<html lang="ru">
	<head>
		<meta charset="UTF-8">
		<title>Экспресс-Аудит Сетевого Пеминетра | ОсОО Квантум Сейф</title>
		<style>
			body { font-family: 'Segoe UI', sans-serif; background: #0b1329; padding: 40px; color: #f4f6f9; line-height: 1.6; }
			.container { max-width: 1000px; margin: 0 auto; }
			.profile-box { background: #1c2541; border: 1px solid #3a506b; padding: 25px; border-radius: 12px; box-shadow: 0 4px 15px rgba(0,0,0,0.2); margin-bottom: 30px; font-size: 15px; }
			.scanner-box { background: #1c2541; padding: 35px; border-radius: 12px; text-align: center; border: 1px solid #3a506b; box-shadow: 0 4px 15px rgba(0,0,0,0.2); }
			.nav-panel { background: #1e3a8a; padding: 15px 25px; border-radius: 8px; margin-bottom: 25px; display: flex; justify-content: space-between; align-items: center; box-shadow: 0 4px 10px rgba(0,0,0,0.1); }
			.nav-link { color: white; text-decoration: none; font-weight: bold; background: #10b981; padding: 10px 20px; border-radius: 6px; transition: background 0.2s; }
			.nav-link:hover { background: #059669; }
			h2 { color: white; font-size: 24px; margin-bottom: 25px; }
			.input-field { width: 55%; padding: 12px 20px; background: #0b1329; border: 1px solid #3a506b; color: white; border-radius: 6px; font-size: 16px; margin-right: 15px; }
			.btn-scan { background: #38bdf8; color: #0b1329; padding: 12px 28px; border-radius: 6px; font-weight: bold; border: none; cursor: pointer; font-size: 16px; transition: background 0.2s; }
			.btn-scan:hover { background: #0ea5e9; }
			.btn-menu { display: inline-block; background: #f59e0b; color: #0b1329; padding: 10px 20px; border-radius: 6px; text-decoration: none; font-weight: bold; margin-right: 15px; margin-top: 25px; }
		</style>
	</head>
	<body>
		<div class="container">
			<!-- Навигационная плашка перехода к банковскому контуру -->
			<div class="nav-panel">
				<span>🏛️ Доступны новые ИИ-модули комплаенса для коммерческих банков КР:</span>
				<a href="/compliance" class="nav-link">Запустить ИИ-Комплаенс (до 4 колена) →</a>
			</div>

			<!-- ПРОФИЛЬ ДЕЯТЕЛЬНОСТИ РЕЗИДЕНТА ПВТ -->
			<div class="profile-box">
				📌 <strong>ПРОФИЛЬ ДЕЯТЕЛЬНОСТИ РЕЗИДЕНТА ПВТ КР ОсОО Квантум Сейф:</strong> Комплексная разработка специализированного программного обеспечения, проектирование интеллектуальных ИИ-ассистентов и автономных ИИ-агентов, сквозной нормативно-правовой и технологический комплаенс-контроль, а также оказание профильной юридической помощи для финансовых организаций, распределенных ИТ-систем и Web3-инфраструктуры.
			</div>

			<!-- БЛОК СКАНИРОВАНИЯ СЕТЕВОГО ПЕРИМЕТРА -->
			<div class="scanner-box">
				<h2>🔍 Экспресс-Аудит Сетевого Периметра (Внешний ИИ-Сканер)</h2>
				<input type="text" class="input-field" placeholder="Введите домен или IP для сканирования...">
				<button class="btn-scan" onclick="alert('Локальное ИИ-сканирование сетевого контура запущено...')">Запустить сканирование</button>
				
				<br>
				<a href="#" class="btn-menu">🏆 Спецификация ПО и Финансовая Выгода</a>
				<a href="#" class="btn-menu" style="background:#38bdf8;">🏪 Маркетплейс ИИ-решений</a>
			</div>
		</div>
		<script>
    const urlParams = new URLSearchParams(window.location.search);
    	<script>
		const urlParams = new URLSearchParams(window.location.search);
		if (urlParams.get('lang') === 'en') {
			document.body.innerHTML = document.body.innerHTML
				.replace(/Запустить ИИ-Комплаенс \(до 4 колена\) →/g, "Launch AI Compliance (up to 4th generation) →")
				.replace(/ПРОФИЛЬ ДЕЯТЕЛЬНОСТИ РЕЗИДЕНТА ПВТ КР ОсОО Квантум Сейф:/g, "BUSINESS PROFILE OF HTP KR RESIDENT Kvantum Safe LLC:")
				.replace(/Комплексная разработка специализированного программного обеспечения, проектирование интеллектуальных ИИ-ассистентов и автономных ИИ-агентов, сквозной нормативно-правовой и технологический комплаенс-контроль, а также оказание профильной юридической помощи для финансовых организаций, распределенных ИТ-систем и Web3-инфраструктуры./g, "Comprehensive development of specialized software, design of intelligent AI assistants and autonomous AI agents, end-to-end regulatory and technological compliance control, as well as specialized legal assistance for financial institutions, distributed IT systems, and Web3 infrastructure.")
				.replace(/Экспресс-Аудит Сетевого Периметра \(Внешний ИИ-Сканер\)/g, "Network Perimeter Express Audit (External AI Scanner)")
				.replace(/Введите домен или IP для сканирования.../g, "Enter bank domain or IP for scanning...")
				.replace(/Запустить сканирование/g, "Run Scan")
				.replace(/Спецификация ПО и Финансовая Выгода/g, "Software Specification & Financial Benefit")
				.replace(/Маркетплейс ИИ-решений/g, "AI Solutions Marketplace");
		}
	</script>
	</body>
	</html>`

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
