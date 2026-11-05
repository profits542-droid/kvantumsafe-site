package main

import (
	"log"
	"net/http"
)

// Определение структуры для перевода страниц
type PageTranslation struct {
	Lang    string
	Content string
}

// Карта для хранения контента на разных языках
var translations = map[string]map[string]string{
	"ru": {
		"home":       GetRussianTemplate(),
		"compliance": GetRussianComplianceTemplate(), // Презентация ИИ-Судьи чеков и проверки Минюста КР
	},
	"en": {
		"home":       GetEnglishTemplate(),
		"compliance": GetEnglishComplianceTemplate(), // Английская версия комплаенс-модуля
	},
}

// Функция для определения языка пользователя (из URL ?lang=en или по умолчанию ru)
func getLanguage(r *http.Request) string {
	lang := r.URL.Query().Get("lang")
	if lang == "en" {
		return "en"
	}
	return "ru" // Язык по умолчанию
}

func main() {
		// 1. ГЛАВНАЯ СТРАНИЦА: Наш оригинальный синий ИИ-Сканер сетевого периметра ПВТ КР
		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			// Строгая проверка, чтобы корневой маршрут не перехватывал другие пути
			if r.URL.Path != "/" {
				http.NotFound(w, r)
				return
			}
	
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			
			// Прямая проверка параметра в ссылке: если ?lang=en, то отдаем английский шаблон
			if r.URL.Query().Get("lang") == "en" {
				w.Write([]byte(GetEnglishTemplate()))
			} else {
				// В противном случае отдаем стандартный русский шаблон
				w.Write([]byte(GetRussianTemplate()))
			}
		})
	
	// 2. ВКЛАДКА КОМПЛАЕНСА: Презентация ИИ-Судьи чеков и проверки до 4 колена Минюста КР
	http.HandleFunc("/compliance", func(w http.ResponseWriter, r *http.Request) {
		lang := getLanguage(r)
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Write([]byte(translations[lang]["compliance"]))
	})

	// Запуск веб-ядра Квантум Сейф
	log.Println("[SUCCESS] Квантум Сейф Веб-Ядро запущено на порту: 8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("[ERROR] Не удалось запустить сервер: %v", err)
	}
}

