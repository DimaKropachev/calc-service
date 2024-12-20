package application

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/DimaKropachev/calculate-web-server/internal/transport"
)

// Структура описывающая приложение
type Application struct {
	Port string
}

// Функция-констурктор для создания нового экземпляра приложения
func New() *Application {
	return &Application{Port: "8080"}
}

// Метод для запуска приложения
func (app *Application) StartServer() error {
	// Создаем файл для логирования 
	logfile, err := os.OpenFile("logfile.txt", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer logfile.Close()

	log.SetOutput(logfile)

	// Создаем мультиплексор запросов
	mux := http.NewServeMux()
	mux.HandleFunc("/api/v1/calculate", transport.CalculateHandler)

	siteHandler := transport.LoggingMiddleware(mux)
	siteHandler = transport.PanicMiddleware(siteHandler)

	// Запускаем сервер
	log.Printf("Сервер запущен на порту: %s\n\n", app.Port)
	if err := http.ListenAndServe(fmt.Sprintf(":%s", app.Port), siteHandler); err != nil {
		// логируем и возвращаем ошибку если не удалось запустить сервер
		log.Fatalf("Не удалось запустить приложение, ошибка: %v\n", err)
		return err
	}
	return nil
}
