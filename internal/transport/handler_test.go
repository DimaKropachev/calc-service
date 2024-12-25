package transport

import (
	"bytes"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

// Структура запроса
type Request struct {
	Method string
	Data   string
}

// Структура ответа
type Response struct {
	StatusCode int
	Result     string
}

// Структура тестового кейса
type TestCase struct {
	Req  *Request
	Resp *Response
}

// Функция для теста обработчика CalculateHandler
func TestCalculateHandler(t *testing.T) {
	// Отключаем логирование на время тестов
	oldLog := log.Writer()
	defer log.SetOutput(oldLog)
	log.SetOutput(io.Discard)

	// Создаем список тестов
	testCases := []TestCase{
		{
			Req: &Request{
				Method: "POST",
				Data:   "{\"expression\": \"2+2\"}",
			},
			Resp: &Response{
				StatusCode: http.StatusOK,
				Result:     "{\"result\":4}",
			},
		},
		{
			Req: &Request{
				Method: "GET",
			},
			Resp: &Response{
				StatusCode: http.StatusMethodNotAllowed,
				Result:     "{\"error\":\"Method Not Allowed\"}",
			},
		},
		{
			Req: &Request{
				Method: "POST",
				Data:   "{\"expression\": \"-2+2-(*3\"}",
			},
			Resp: &Response{
				StatusCode: http.StatusUnprocessableEntity,
				Result:     "{\"error\":\"Expression is not valid\"}",
			},
		},
		{
			Req: &Request{
				Method: "POST",
				Data:   "{\"expression\": 17}",
			},
			Resp: &Response{
				StatusCode: http.StatusBadRequest,
				Result:     "{\"error\":\"Invalid JSON format\"}",
			},
		},
	}

	for numCase, testCase := range testCases {
		url := "http://localhost/api/v1/calculate"
		// Создаем поле data для POST запроса
		data := bytes.NewBuffer([]byte(testCase.Req.Data))

		// Создаем новый HTTP запрос
		req := httptest.NewRequest(testCase.Req.Method, url, data)
		// Создаем новый HTTP рекодер для записи ответа
		w := httptest.NewRecorder()

		// Вызываем функцию-обработчик
		CalculateHandler(w, req)

		// Проверка статуса ответа
		if w.Code != testCase.Resp.StatusCode {
			t.Errorf("[%d] wrong StatusCode: got %d, expected %d", numCase, w.Code, testCase.Resp.StatusCode)
		}

		// Получаем тело ответа
		resp := w.Result()
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			t.Fatalf("failed to read body: %v", err)
		}

		// Сравниваем тело ответа с ожидаемым результатом
		bodyStr := string(body)
		if bodyStr != testCase.Resp.Result {
			t.Errorf("[%d] wrong Response: got %+v, expected %+v", numCase, bodyStr, testCase.Resp.Result)
		}
	}
}
