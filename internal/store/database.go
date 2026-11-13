package store

import (
	"context"
	"fmt"
	"strings"
)

type Block struct {
	ID    int64
	Key   string
	Lang  string
	Value string
}

type Product struct {
	ID     int64
	Key    string
	NameRU string
	NameEN string
	DescRU string
	DescEN string
	Price  string
	Status string
}

func (s *Store) migrateContent() error {
	_, err := s.db.Exec(`
		CREATE TABLE IF NOT EXISTS content (
			id INTEGER PRIMARY KEY,
			"key" TEXT NOT NULL,
			lang TEXT NOT NULL,
			value TEXT NOT NULL,
			UNIQUE("key", lang)
		);
		CREATE TABLE IF NOT EXISTS products (
			id INTEGER PRIMARY KEY,
			name_ru TEXT NOT NULL,
			name_en TEXT NOT NULL,
			desc_ru TEXT NOT NULL,
			desc_en TEXT NOT NULL,
			price TEXT NOT NULL,
			status TEXT NOT NULL
		);
	`)
	if err != nil {
		return err
	}
	if _, err = s.db.Exec(`ALTER TABLE products ADD COLUMN product_key TEXT`); err != nil && !strings.Contains(err.Error(), "duplicate column") {
		return err
	}
	_, err = s.db.Exec(`CREATE UNIQUE INDEX IF NOT EXISTS products_product_key ON products(product_key) WHERE product_key IS NOT NULL AND product_key <> ''`)
	return err
}

func (s *Store) seedIfEmpty(ctx context.Context) error {
	if err := s.seedContent(ctx); err != nil {
		return err
	}
	if err := s.refreshOfficialLegal(ctx); err != nil {
		return err
	}
	return s.seedProducts(ctx)
}

func (s *Store) refreshOfficialLegal(ctx context.Context) error {
	for _, row := range contentSeed() {
		switch row.key {
		case "legal_company", "legal_patent", "about_residency", "about_legal", "about_patent", "htp_certificate_label":
		default:
			continue
		}
		var current string
		err := s.db.QueryRowContext(ctx, `SELECT value FROM content WHERE "key" = ? AND lang = ?`, row.key, row.lang).Scan(&current)
		if err != nil {
			continue
		}
		if current == row.value || strings.Contains(current, "799-28-10-2026") {
			continue
		}
		if _, err := s.db.ExecContext(ctx, `UPDATE content SET value = ? WHERE "key" = ? AND lang = ?`, row.value, row.key, row.lang); err != nil {
			return err
		}
	}
	return nil
}

func (s *Store) seedContent(ctx context.Context) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	stmt, err := tx.PrepareContext(ctx, `INSERT OR IGNORE INTO content ("key", lang, value) VALUES (?, ?, ?)`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, row := range contentSeed() {
		if _, err := stmt.ExecContext(ctx, row.key, row.lang, row.value); err != nil {
			return fmt.Errorf("seed %s/%s: %w", row.key, row.lang, err)
		}
	}
	return tx.Commit()
}

func (s *Store) seedProducts(ctx context.Context) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	for _, item := range productSeed() {
		if _, err := tx.ExecContext(ctx, `
			UPDATE products SET product_key = ?
			WHERE (product_key IS NULL OR product_key = '') AND name_en = ?`,
			item.Key, item.NameEN,
		); err != nil {
			return err
		}
		var n int
		if err := tx.QueryRowContext(ctx, `SELECT COUNT(*) FROM products WHERE product_key = ?`, item.Key).Scan(&n); err != nil {
			return err
		}
		if n > 0 {
			continue
		}
		if _, err := tx.ExecContext(ctx, `
			INSERT INTO products (product_key, name_ru, name_en, desc_ru, desc_en, price, status)
			VALUES (?, ?, ?, ?, ?, ?, ?)`,
			item.Key, item.NameRU, item.NameEN, item.DescRU, item.DescEN, item.Price, item.Status,
		); err != nil {
			return err
		}
	}
	return tx.Commit()
}

func (s *Store) ContentByLang(ctx context.Context, lang string) (map[string]string, error) {
	rows, err := s.db.QueryContext(ctx, `SELECT "key", value FROM content WHERE lang = ?`, lang)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := make(map[string]string)
	for rows.Next() {
		var key, value string
		if err := rows.Scan(&key, &value); err != nil {
			return nil, err
		}
		out[key] = value
	}
	return out, rows.Err()
}

func (s *Store) ListContent(ctx context.Context) ([]Block, error) {
	rows, err := s.db.QueryContext(ctx, `SELECT id, "key", lang, value FROM content ORDER BY "key", lang`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []Block
	for rows.Next() {
		var item Block
		if err := rows.Scan(&item.ID, &item.Key, &item.Lang, &item.Value); err != nil {
			return nil, err
		}
		out = append(out, item)
	}
	return out, rows.Err()
}

func (s *Store) UpdateContent(ctx context.Context, key, lang, value string) error {
	res, err := s.db.ExecContext(ctx, `UPDATE content SET value = ? WHERE "key" = ? AND lang = ?`, value, key, lang)
	if err != nil {
		return err
	}
	n, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if n == 0 {
		return fmt.Errorf("блок %s/%s не найден", key, lang)
	}
	return nil
}

func (s *Store) ListProducts(ctx context.Context) ([]Product, error) {
	rows, err := s.db.QueryContext(ctx, `
		SELECT id, COALESCE(product_key, ''), name_ru, name_en, desc_ru, desc_en, price, status
		FROM products ORDER BY id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []Product
	for rows.Next() {
		var item Product
		if err := rows.Scan(&item.ID, &item.Key, &item.NameRU, &item.NameEN, &item.DescRU, &item.DescEN, &item.Price, &item.Status); err != nil {
			return nil, err
		}
		out = append(out, item)
	}
	return out, rows.Err()
}

func (s *Store) UpdateProduct(ctx context.Context, item Product) error {
	res, err := s.db.ExecContext(ctx, `
		UPDATE products
		SET name_ru = ?, name_en = ?, desc_ru = ?, desc_en = ?, price = ?, status = ?
		WHERE id = ?`,
		item.NameRU, item.NameEN, item.DescRU, item.DescEN, item.Price, item.Status, item.ID,
	)
	if err != nil {
		return err
	}
	n, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if n == 0 {
		return fmt.Errorf("продукт %d не найден", item.ID)
	}
	return nil
}

type seedRow struct {
	key   string
	lang  string
	value string
}

const (
	officialLegalRU = "Официальный резидент Парка высоких технологий Кыргызской Республики (Свидетельство № 799/10-С-26 от 28.10.2026, рег. № в реестре резидентов 799-28-10-2026). ОсОО «Квантум Сейф», ИНН 00209202610130, ОКПО 35220502. Государственное правовое депонирование ПО: патент Кыргызпатента № 1295 от 01.10.2026 г."
	officialLegalEN = "Official Resident of the High Technology Park of the Kyrgyz Republic (Certificate No. 799/10-С-26 dated Oct 28, 2026, Registry ID 799-28-10-2026). Quantum Safe LLC, TIN 00209202610130. State IP Deposit: Kyrgyzpatent Certificate No. 1295 dated Oct 01, 2026."
)

func contentSeed() []seedRow {
	return []seedRow{
		{key: "hero_kicker", lang: "ru", value: "Kvantum Safe Pro · .tech"},
		{key: "hero_kicker", lang: "en", value: "Kvantum Safe Pro · .tech"},
		{key: "hero_title", lang: "ru", value: "Постквантовая экосистема кибербезопасности нового поколения"},
		{key: "hero_title", lang: "en", value: "Next-Gen Post-Quantum Cybersecurity Ecosystem"},
		{key: "hero_lead", lang: "ru", value: "Kvantum Safe Pro — модульная платформа защиты финансового сектора и инфраструктуры Web3. Она закрывает сегодняшние угрозы периметру банка и заранее учитывает нагрузку квантовых вычислений на криптографию, которая уже стоит в промышленном контуре."},
		{key: "hero_lead", lang: "en", value: "Kvantum Safe Pro is a modular protection platform for banks and Web3 infrastructure. It addresses present-day threats to the financial perimeter and the coming pressure of quantum computing on cryptography now in production."},
		{key: "crypto_stack", lang: "ru", value: "Криптографический профиль поставки: гибрид NIST ML-KEM-768 и GmSSL SM4-GCM для сообщений SWIFT и XML."},
		{key: "crypto_stack", lang: "en", value: "Delivery cryptographic profile: a hybrid of NIST ML-KEM-768 and GmSSL SM4-GCM for SWIFT and XML messages."},
		{key: "intake_label", lang: "ru", value: "Домен организации"},
		{key: "intake_label", lang: "en", value: "Institution domain"},
		{key: "intake_button", lang: "ru", value: "Запустить сканирование периметра"},
		{key: "intake_button", lang: "en", value: "Run Perimeter Scan"},
		{key: "intake_hint", lang: "ru", value: "Открывает экспресс-отчёт по указанному домену."},
		{key: "intake_hint", lang: "en", value: "Opens the express report for the domain you enter."},
		{key: "products_title", lang: "ru", value: "Продуктовая триада"},
		{key: "products_title", lang: "en", value: "Core products"},
		{key: "deploy_title", lang: "ru", value: "Архитектура и внедрение"},
		{key: "deploy_title", lang: "en", value: "Architecture and deployment"},
		{key: "gateway_title", lang: "ru", value: "Автономный бронированный шлюз"},
		{key: "gateway_title", lang: "en", value: "Autonomous hardened gateway"},
		{key: "gateway_text", lang: "ru", value: "Reverse proxy в демилитаризованной зоне. Стоит перед системами банка и остаётся единственным входом из публичной сети."},
		{key: "gateway_text", lang: "en", value: "A reverse proxy in the demilitarized zone. It stands in front of core banking systems and is the only path in from the public network."},
		{key: "sdk_title", lang: "ru", value: "Компилируемое SDK"},
		{key: "sdk_title", lang: "en", value: "Compilable SDK"},
		{key: "sdk_text", lang: "ru", value: "Статическая библиотека .lib или .so, встраиваемая в код АБС банка. Экземпляр привязывается к оборудованию сервера."},
		{key: "sdk_text", lang: "en", value: "A static library, .lib or .so, built into the bank’s core-banking code. Each deployment is bound to the server hardware."},
		{key: "legal_title", lang: "ru", value: "Юридический мандат и комплаенс"},
		{key: "legal_title", lang: "en", value: "Legal mandate and compliance"},
		{key: "legal_company", lang: "ru", value: officialLegalRU},
		{key: "legal_company", lang: "en", value: officialLegalEN},
		{key: "legal_patent", lang: "ru", value: "Государственное правовое депонирование ПО: патент Кыргызпатента № 1295 от 01.10.2026 г."},
		{key: "legal_patent", lang: "en", value: "State IP Deposit: Kyrgyzpatent Certificate No. 1295 dated Oct 01, 2026."},
		{key: "htp_certificate_label", lang: "ru", value: "Открыть Свидетельство резидента ПВТ (PDF)"},
		{key: "htp_certificate_label", lang: "en", value: "View HTP Resident Certificate (PDF)"},
		{key: "about_kicker", lang: "ru", value: "ОсОО «Квантум Сейф»"},
		{key: "about_kicker", lang: "en", value: "Kvantum Safe LLC"},
		{key: "about_title", lang: "ru", value: "О компании"},
		{key: "about_title", lang: "en", value: "About the company"},
		{key: "about_lead", lang: "ru", value: "Kvantum Safe разрабатывает модульную защиту для банков и инфраструктуры Web3: периметр, комплаенс и постквантовый контур поставок."},
		{key: "about_lead", lang: "en", value: "Kvantum Safe builds modular protection for banks and Web3 infrastructure: the perimeter, compliance, and the post-quantum delivery profile."},
		{key: "about_residency_title", lang: "ru", value: "Резидентство ПВТ КР"},
		{key: "about_residency_title", lang: "en", value: "HTP residency"},
		{key: "about_residency", lang: "ru", value: officialLegalRU},
		{key: "about_residency", lang: "en", value: officialLegalEN},
		{key: "about_legal_title", lang: "ru", value: "Юридические реквизиты"},
		{key: "about_legal_title", lang: "en", value: "Legal details"},
		{key: "about_legal", lang: "ru", value: officialLegalRU},
		{key: "about_legal", lang: "en", value: officialLegalEN},
		{key: "about_patent", lang: "ru", value: "Государственное правовое депонирование ПО: патент Кыргызпатента № 1295 от 01.10.2026 г."},
		{key: "about_patent", lang: "en", value: "State IP Deposit: Kyrgyzpatent Certificate No. 1295 dated Oct 01, 2026."},
		{key: "about_team_title", lang: "ru", value: "Команда"},
		{key: "about_team_title", lang: "en", value: "Team"},
		{key: "about_role_1_title", lang: "ru", value: "Инженерная группа"},
		{key: "about_role_1_title", lang: "en", value: "Engineering"},
		{key: "about_role_1_text", lang: "ru", value: "Собирает сканер периметра, внутренний контроль файлов и криптографический мост в поставляемые модули."},
		{key: "about_role_1_text", lang: "en", value: "Builds the perimeter scanner, the internal file control, and the cryptographic bridge into deliverable modules."},
		{key: "about_role_2_title", lang: "ru", value: "Комплаенс"},
		{key: "about_role_2_text", lang: "ru", value: "Сверяет контур поставки с санкционными требованиями и правилами доступа к материалам банка."},
		{key: "about_role_2_title", lang: "en", value: "Compliance"},
		{key: "about_role_2_text", lang: "en", value: "Checks the delivery profile against sanctions requirements and the bank’s rules for access to sensitive material."},
		{key: "about_role_3_title", lang: "ru", value: "Внедрение"},
		{key: "about_role_3_text", lang: "ru", value: "Ставит шлюз в демилитаризованной зоне или встраивает SDK в автоматизированную банковскую систему."},
		{key: "about_role_3_title", lang: "en", value: "Deployment"},
		{key: "about_role_3_text", lang: "en", value: "Places the gateway in the demilitarized zone or builds the SDK into the core banking system."},
		{key: "marketplace_kicker", lang: "ru", value: "Kvantum Safe Pro"},
		{key: "marketplace_kicker", lang: "en", value: "Kvantum Safe Pro"},
		{key: "marketplace_title", lang: "ru", value: "Маркетплейс ИИ-решений"},
		{key: "marketplace_title", lang: "en", value: "AI solutions marketplace"},
		{key: "marketplace_lead", lang: "ru", value: "Карточки продуктов, их спецификации и цены берутся из базы. Статус и стоимость правятся в админке вместе с текстами."},
		{key: "marketplace_lead", lang: "en", value: "Product cards, specifications, and prices come from the database. Status and price are edited in the admin desk together with the copy."},
	}
}

func productSeed() []Product {
	return []Product{
		{
			Key:    "external_scanner",
			NameRU: "Внешний инфраструктурный сканер",
			NameEN: "External infrastructure scanner",
			DescRU: "Модуль scanner.go. Бесплатный аудит внешних TLS-шлюзов банка: сертификат, протокол на входе и то, что организация показывает в интернет.",
			DescEN: "Module scanner.go. A no-fee review of a bank’s public TLS gateways: the certificate, the protocol at the edge, and what the institution exposes to the internet.",
			Price:  "0",
			Status: "active",
		},
		{
			Key:    "internal_scanner",
			NameRU: "Внутренний сканер комплаенса",
			NameEN: "Internal compliance scanner",
			DescRU: "Модуль internal_fs.go. Ищет внутри контура банка материалы, которым не место в открытом доступе: выгрузки баз, SQL-дампы и приватные ключи Web3. Найденное изолируется правами 0600 — чтение и запись только у владельца.",
			DescEN: "Module internal_fs.go. Looks inside the bank’s own environment for material that should not sit in the open: database exports, SQL dumps, and Web3 private keys. What it finds is locked to mode 0600, readable and writable by the owner alone.",
			Price:  "included",
			Status: "active",
		},
		{
			Key:    "quantum_bridge",
			NameRU: "Постквантовый крипто-мост",
			NameEN: "Post-quantum crypto bridge",
			DescRU: "Модуль quantum_bridge.go. Сквозная защита сообщений SWIFT и XML. В профиле поставки — гибрид стандартов NIST ML-KEM-768 и GmSSL SM4-GCM.",
			DescEN: "Module quantum_bridge.go. End-to-end protection for SWIFT and XML payment messages. The delivery profile uses the hybrid of NIST ML-KEM-768 and GmSSL SM4-GCM.",
			Price:  "35000 USD / year",
			Status: "active",
		},
		{
			Key:    "compliance_orchestrator",
			NameRU: "ИИ-Оркестратор Комплаенса и Антифрода",
			NameEN: "AI Compliance & Anti-Fraud Orchestrator",
			DescRU: "Модуль комплексной защиты для банков и криптобирж. Включает ИИ-Судью P2P-чеков для предотвращения мошенничества и автоматический AML-скрининг по базам OFAC/Минюста КР до 4-го поколения родственников для защиты от вторичных санкций.",
			DescEN: "Comprehensive enterprise security module for commercial banks and crypto exchanges. Features an AI Receipt Judge to prevent P2P fraud and automated AML screening via OFAC/MoJ registries up to the 4th generation of kinship to eliminate secondary sanction risks.",
			Price:  "$25,000 / year",
			Status: "active",
		},
	}
}
