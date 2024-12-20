package transport

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"

	"github.com/DimaKropachev/calculate-web-server/pkg/calculate"
)

var ErrInvalidExpression = errors.New("Expression is not valid")

// Структура для парсинга выражения из JSON
type Expression struct {
	Expr string `json:"expression"`
}

// Структура для представления ошибки в JSON
type ResponseError struct {
	Error string `json:"error"`
}

// Структура для представления результата в JSON
type ResponseResult struct {
	Result float64 `json:"result"`
}

/*
Функция CalculateHandler обрабатывает HTTP-запросы.
Она принимает JSON с выражением, выполняет вычисление и возвращает результат или ошибку
*/
func CalculateHandler(w http.ResponseWriter, r *http.Request) {
	// Устанавливаем заголовок для HTTP-ответа
	w.Header().Add("Content-Type", "application/json")

	// Проверяем метод запроса
	if r.Method != http.MethodPost {
		log.Printf("Responce:\n[ERROR] метод %s не разрешен\n", r.Method)
		w.WriteHeader(http.StatusMethodNotAllowed)
		SendErrorJSON(w, http.StatusText(http.StatusMethodNotAllowed))
		return
	}

	// Читаем тело зпроса
	data, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("Response:\n[ERROR] не удалось прочитать тело запроса: %v\n", err)
		w.WriteHeader(http.StatusBadRequest)
		SendErrorJSON(w, "Не удалось прочитать тело запроса")
		return
	}
	// В конце обработки запроса закрываем тело запроса
	defer r.Body.Close()

	// Создаем новый экземпляр структуры expression
	expression := &Expression{}

	// Распарсиваем JSON данные с выражением в структуру Expression
	err = json.Unmarshal(data, expression)
	if err != nil {
		log.Printf("Response:\n[ERROR] ошибка при распарсивании JSON: %v\n", err)
		w.WriteHeader(http.StatusBadRequest)
		SendErrorJSON(w, "Invalid JSON format")
		return
	}

	// Вызываем функцию для вычисления выражения
	result, err := calculate.Calc(expression.Expr)
	if err != nil {
		// Если произошла ошибка, логируем сообщение об ошибке и отправляем статус 422
		log.Printf("Response:\nExpression is not valid:\n%v", err)
		w.WriteHeader(http.StatusUnprocessableEntity)
		SendErrorJSON(w, ErrInvalidExpression.Error())
	} else {
		// Если вычисление прошло успешно, логируем результат и отправляем статус 200
		log.Printf("Response:\nResult = %f\n", result)
		w.WriteHeader(http.StatusOK)
		SendResultJSON(w, result)
	}

}

// Функция для отправки сообщения об ошибке в виде JSON
func SendErrorJSON(w http.ResponseWriter, err string) {
	resp := &ResponseError{}
	resp.Error = err

	data, jsonErr := json.Marshal(resp)
	if jsonErr != nil {
		log.Printf("[ERROR] ошибка при формировании JSON ответа: %v\n", jsonErr)
		http.Error(w, "Ошибка при формировании ответа", http.StatusInternalServerError)
		return
	}

	w.Write(data)
}

// Функция для отправки результата в виде JSON
func SendResultJSON(w http.ResponseWriter, result float64) {
	resp := &ResponseResult{}
	resp.Result = result

	data, jsonErr := json.Marshal(resp)
	if jsonErr != nil {
		log.Printf("[ERROR] ошибка при формировании JSON ответа: %v\n", jsonErr)
		http.Error(w, "Ошибка при формировании ответа", http.StatusInternalServerError)
		return
	}

	w.Write(data)
}
