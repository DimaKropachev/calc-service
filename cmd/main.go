package main

import (
	"fmt"
	"os"

	"github.com/DimaKropachev/calculate-web-server/internal/application"
)

func main() {
	// Создаем новый экземпляр приложения
	app := application.New()
	// Запускаем сервер
	err := app.StartServer()
	if err != nil {
		// Выводим сообщение об ошибке в консоль
		fmt.Println("Ошибка при запуске сервера:", err)
		// Завершаем работу программы с кодом ошибки
		os.Exit(1)
	}
}
