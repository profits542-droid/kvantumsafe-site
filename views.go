package main

// GetRussianTemplate возвращает русскоязычную версию сайта для банков КР и СНГ
func GetRussianTemplate() string {
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
			.lang-switch:hover { background: #1d4ed8; }
			h1 { color: #1e3a8a; border-bottom: 2px solid #e5e7eb; padding-bottom: 10px; margin-top: 0; }
			h3 { color: #0f172a; margin-top: 25px; }
			ul { padding-left: 20px; }
			li { margin-bottom: 8px; }
			.demo-box { background: #f8fafc; border: 2px dashed #1e3a8a; padding: 25px; border-radius: 8px; margin-top: 30px; }
			.btn { display: inline-block; background: #10b981; color: white; padding: 12px 24px; border-radius: 6px; text-decoration: none; font-weight: bold; margin-right: 12px; margin-top: 12px; transition: background 0.2s; }
			.btn:hover { background: #059669; }
			.badge { background: #e0f2fe; color: #0369a1; padding: 4px 8px; border-radius: 4px; font-size: 14px; font-weight: bold; }
		</style>
	</head>
	<body>
		<div class="container">
			<a href="/en" class="lang-switch">🌐 EN Version (Grant)</a>
			<h1>🛡️ Quantum Safe AI</h1>
			<p><strong>Автономная ИИ-инфраструктура комплаенса для коммерческих банков и P2P-платформ Central Asia.</strong></p>
			<p><span class="badge">Резидент ПВТ</span> Официальный резидент Парка Высоких Технологий Кыргызской Республики.</p>
			
			<h3>🎯 Проблема, которую мы решаем:</h3>
			<p>Финтех-сектор и банковская система Центральной Азии несут убытки из-за риска <strong>вторичных санкций (OFAC, SDN, ЕС)</strong> и мошеннических схем типа «Треугольник». Традиционный ручной аудит занимает много времени, а облачные ИИ-сервисы нарушают банковскую тайну.</p>

			<h3>🚀 Модули Квантум-Контура:</h3>
			<ul>
				<li><strong>Quantum AI P2P Referee:</strong> Локальный движок компьютерного зрения (OCR Tesseract v5), который за 2 секунды распознает чеки, извлекая точные суммы, банки и ФИО плательщика капсом без отправки данных в интернет.</li>
				<li><strong>Банковский Комплаенс-Шлюз:</strong> Автоматический построитель графа скрытых взаимосвязей <strong>до 4-го колена</strong> (соучредители, связанные компании, родственники бенефициаров) по базам МинЮста КР и санкционным спискам.</li>
				<li><strong>Высокопроизводительное Ядро Go:</strong> Бэкенд на Golang обрабатывает тысячи запросов в секунду через легковесные горутины.</li>
			</ul>

			<div class="demo-box">
				<h3>🤖 ИНТЕРАКТИВНЫЙ СТЕНД (LIVE DEMO В ТЕЛЕГРАМ):</h3>
				<p>Протестируйте работу модулей нашего ИИ-ядра в реальном времени через официальные демонстрационные шлюзы:</p>
				<a href="https://telegram.me/kvantumsafe_p2p_bot" class="btn" target="_blank">⚡ Запустить ИИ-Судью Чеков</a>
				<a href="https://t.me/KvantumSafe_Lead_Bot" class="btn" style="background:#1e3a8a;" target="_blank">🏛️ Запустить Банковский Шлюз</a>
			</div>
		</div>
	</body>
	</html>`
}

// GetEnglishTemplate возвращает англоязычную версию сайта для фонда .TECH
func GetEnglishTemplate() string {
	return `
	<!DOCTYPE html>
	<html lang="en">
	<head>
		<meta charset="UTF-8">
		<title>Quantum Safe AI | Enterprise Bank Compliance & P2P Anti-Fraud</title>
		<style>
			body { font-family: 'Segoe UI', sans-serif; background: #f4f6f9; padding: 40px; color: #333; line-height: 1.6; }
			.container { max-width: 900px; margin: 0 auto; background: white; padding: 40px; border-radius: 12px; box-shadow: 0 4px 15px rgba(0,0,0,0.05); }
			.lang-switch { float: right; background: #1e3a8a; color: white; padding: 8px 15px; border-radius: 6px; text-decoration: none; font-weight: bold; }
			.lang-switch:hover { background: #1d4ed8; }
			h1 { color: #1e3a8a; border-bottom: 2px solid #e5e7eb; padding-bottom: 10px; margin-top: 0; }
			h3 { color: #0f172a; margin-top: 25px; }
			ul { padding-left: 20px; }
			li { margin-bottom: 8px; }
			.demo-box { background: #f8fafc; border: 2px dashed #1e3a8a; padding: 25px; border-radius: 8px; margin-top: 30px; }
			.btn { display: inline-block; background: #10b981; color: white; padding: 12px 24px; border-radius: 6px; text-decoration: none; font-weight: bold; margin-right: 12px; margin-top: 12px; transition: background 0.2s; }
			.btn:hover { background: #059669; }
			.badge { background: #e0f2fe; color: #0369a1; padding: 4px 8px; border-radius: 4px; font-size: 14px; font-weight: bold; }
		</style>
	</head>
	<body>
		<div class="container">
			<a href="/" class="lang-switch">🌐 RU Version</a>
			<h1>🛡️ Quantum Safe AI</h1>
			<p><strong>Autonomous On-Premise Compliance Infrastructure for Central Asian Banks & Crypto Platforms.</strong></p>
			<p><span class="badge">HTP Resident</span> Official resident of the High-Tech Park of the Kyrgyz Republic.</p>
			
			<h3>🎯 The Problem We Solve:</h3>
			<p>The fintech and banking sector in Central Asia faces severe risks from <strong>secondary sanctions (OFAC, SDN, EU)</strong> and P2P "Triangle Fraud". Traditional manual auditing is highly inefficient, and third-party cloud APIs leak confidential customer transactions.</p>

			<h3>🚀 Core Technological Modules:</h3>
			<ul>
				<li><strong>Quantum AI P2P Referee:</strong> Local on-premise Tesseract OCR engine that parses payment receipts in under 2 seconds, extracting exact amounts, banks, and full names without external cloud leaks.</li>
				<li><strong>Bank Compliance Gate:</strong> Automated B2B gateway scanning transactions up to the 4th level of affiliation (cross-checking Kyrgyz Ministry of Justice registries and global sanction lists).</li>
				<li><strong>High-Performance Go Backend:</strong> Built with pure Golang, utilizing lightweight concurrent goroutines to process thousands of transactions per second.</li>
			</ul>

			<div class="demo-box">
				<h3>🤖 INTERACTIVE SANDBOX (TELEGRAM DEMO BUILDS):</h3>
				<p>You can audit our core AI processing loops in real-time using our live sandboxed interface:</p>
				<a href="https://t.me/kvantumsafe_p2p_bot" class="btn" target="_blank">⚡ Launch P2P Referee Demo</a>
				<a href="https://t.me/KvantumSafe_Lead_Bot" class="btn" style="background:#1e3a8a;" target="_blank">🏛️ Launch Bank Compliance Gate</a>
			</div>
		</div>
	</body>
	</html>`
}
