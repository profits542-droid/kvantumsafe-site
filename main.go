package main

import (
	"log"
	"net/http"
)

func main() {
	// 1. Маршрут для русскоязычной версии (Банки КР и СНГ)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Строгая проверка, чтобы роут не перехватывал другие несуществующие пути
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Write([]byte(GetRussianTemplate()))
	})

	// 2. Маршрут для англоязычной версии (Специально для гранта .TECH)
	http.HandleFunc("/en", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Write([]byte(GetEnglishTemplate()))
	})

	// 3. Запуск веб-сервера сайта на стандартном порту 80
	// Если порт 80 занят в системе, можно временно переключить на ":8081"
	port := ":80"
	log.Printf("🌍 Официальный сайт Quantum Safe успешно запущен на http://localhost%s", port)
	
	err := http.ListenAndServe(port, nil)
	if err != nil {
		log.Printf("⚠️ Порт 80 занят, переключаемся на резервный порт 8081...")
		port = ":8081"
		log.Printf("🌍 Сайт Quantum Safe запущен на http://localhost%s", port)
		if errRetry := http.ListenAndServe(port, nil); errRetry != nil {
			log.Fatalf("🚨 Критическая ошибка запуска веб-сервера: %v", errRetry)
		}
	}
}
