package calculate

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/DimaKropachev/calculate-web-server/pkg/errors"
)

// Функция проверяет правильность написания выражения
func CheckExpression(tokens []string) error {
	// Создаем map с ошибками
	Errors := map[error][]string{
		errors.ErrInvalidToken:        {},
		errors.ErrBeforeFirstBracket:  {},
		errors.ErrAfterFirstBracket:   {},
		errors.ErrBeforeSecondBracket: {},
		errors.ErrAfterSecondBracket:  {},
		errors.ErrOperationsContract:  {},
		errors.ErrIntegersContract:    {},
	}

	var (
		resultError   string
		countBrackets int
		trueBrackets  bool
	)

	// Если выражение пустое сразу возвращаем ошибку
	if len(tokens) == 0 {
		return errors.ErrEmptyExpression
	}

	// Считаем скобки в выражении
	for _, token := range tokens {
		if token == "(" {
			countBrackets++
		} else if token == ")" {
			countBrackets--
		}
		// если нарушен порядок скобок то добавляем в resultError ошибку об этом
		if countBrackets < 0 {
			resultError += fmt.Sprintf("[ERROR] %v\n", errors.ErrBracketsInExpression)
			break
		}
	}

	// Если скобки расставлены правильно, то проверяем содержимое внутри скобок
	if countBrackets == 0 {
		trueBrackets = true
		err := CheckBrackets(tokens)
		if err != nil {
			resultError += fmt.Sprintf("[ERROR] %v", err)
		}
	} else {
		trueBrackets = false
		resultError += fmt.Sprintf("[ERROR] %v\n", errors.ErrBracketsInExpression)
	}

	// Создаем счетчики подряд идущих математических операций и чисел
	var (
		countNum, countOper int
	)

	// Перебираем токены выражения
	for i, token := range tokens {
		// Если token стоит на 1 месте и является знаком, то добавляем в resultError ошибку об этом
		if i == 0 {
			if IsOperation(token) {
				resultError += fmt.Sprintf("[ERROR] %v, expression[%d]\n", errors.ErrFirstOperation, i)
			}
		}
		// Если token стоит на последнем месте и является знаком, то добавляем в resultError ошибку об этом
		if i == len(tokens)-1 {
			if IsOperation(token) {
				resultError += fmt.Sprintf("[ERROR] %v, expression[%d]\n", errors.ErrLastOperation, i)
			}
		}
		if token == "(" {
			// Если token - открывающая скобка
			// Обнуляем счетчики подряд идущих чисел и математических операций
			countNum = 0
			countOper = 0

			// Если скобки в выражении расставлены правильно, то проверяем стоящие перед и после скобки token-ы
			// Если token не соответствует правилам, то добавляем в map с ошибками индекс этого токена
			if trueBrackets {
				if tokens[i+1] == ")" {
					continue
				}
				if !IsInteger(tokens[i+1]) {
					Errors[errors.ErrAfterFirstBracket] = append(Errors[errors.ErrAfterFirstBracket], fmt.Sprintf("%d", i))
				}
				if i > 0 {
					if !IsOperation(tokens[i-1]) {
						Errors[errors.ErrBeforeFirstBracket] = append(Errors[errors.ErrBeforeFirstBracket], fmt.Sprintf("%d", i))
					}
				}
			}
		} else if token == ")" {
			// Если token - закрывающая скобка
			// Обнуляем счетчики подряд идущих чисел и математических операций
			countNum = 0
			countOper = 0

			// Если скобки в выражении расставлены правильно, то проверяем стоящие перед и после скобки token-ы
			// Если token не соответствует правилам, то добавляем в map с ошибками индекс этого токена
			if trueBrackets {
				if tokens[i-1] == "(" {
					continue
				}
				if !IsInteger(tokens[i-1]) {
					Errors[errors.ErrBeforeSecondBracket] = append(Errors[errors.ErrBeforeSecondBracket], fmt.Sprintf("%d", i))
				}
				if i < len(tokens)-2 {
					if !IsOperation(tokens[i+1]) {
						Errors[errors.ErrAfterSecondBracket] = append(Errors[errors.ErrAfterSecondBracket], fmt.Sprintf("%d", i))
					}
				}
			}
		} else {
			// Если token - не скобка
			if IsOperation(token) {
				// Если token является математической операцией, то счетчик подряд идущих математический операций увеличиваем, а счетчик подряд идущих чисел обнуляем
				countOper++
				countNum = 0
			} else if IsInteger(token) {
				// Если token является числом, то счетчик подряд идущих чисел увеличиваем, а счетчик подряд идущих математических операций обнуляем
				countNum++
				countOper = 0
			} else {
				// Если token не является ни числом, ни математической операцией, то добавляем в map с ошибками ошибку о неизвестном token-е и индекс данного token-а
				// Также обнуляем счетчики подряд идущих чисел и математических операций
				countNum = 0
				countOper = 0
				Errors[errors.ErrInvalidToken] = append(Errors[errors.ErrInvalidToken], fmt.Sprintf("%d", i))
			}

			// Если счетчик подряд идущих чисел равен 2, то добавлем в map с ошибками ошибку об этом и индекс первого подряд идущего token-а
			if countNum == 2 {
				Errors[errors.ErrIntegersContract] = append(Errors[errors.ErrIntegersContract], fmt.Sprintf("%d", i-1))
			}
			// Если счетчик подряд идущих математических операций равен 2, то добавлем в map с ошибками ошибку об этом и индекс первого подряд идущего token-а
			if countOper == 2 {
				Errors[errors.ErrOperationsContract] = append(Errors[errors.ErrOperationsContract], fmt.Sprintf("%d", i-1))
			}
		}
	}

	// Перебираем map с ошибками и добавляем информацию о них в resultError
	for err, indexes := range Errors {
		if len(indexes) > 0 {
			resultError += fmt.Sprintf("[ERROR] %v, expression[%v]\n", err, strings.Join(indexes, ","))
		}
	}

	// Если resultError не пустая, то добавляем в начало список токенов для удобного чтения ошибок, иначе вернем nil
	if resultError != "" {
		return fmt.Errorf("expression = %v\n%v", tokens, resultError)
	}
	return nil
}

// Функция проверяет наличие содержимого в скобках
func CheckBrackets(tokens []string) error {
	resultError := ""
	brackets := [][]int{}
	emptyBrackets := []string{}
	// Создаем map-ы для индексов открывающихся и закрывающихся скобок
	mapFirst := make(map[int]int)
	mapSecond := make(map[int]int)
	// Переменная для нумерации скобок 
	cb := 0
	// Перебираем token-ы выражения
	for i, token := range tokens {
		// Если token - открывающаяся скобка, то добавляем его индекс в соответствующую map-у под соответствующим номером (cb)
		if token == "(" {
			cb++
			if _, ok := mapFirst[cb]; !ok {
				mapFirst[cb] = i
			}
		} else if token == ")" {
			// Если token - закрывающаяся скобка, то добавляем его индекс в соответствующую map-у под соответствующим номером (cb)
			if _, ok := mapSecond[cb]; !ok {
				mapSecond[cb] = i
			} else {
				for {
					cb--
					if _, ok := mapSecond[cb]; !ok {
						mapSecond[cb] = i
						break
					}
				}
			}
		}
	}

	// Перебираем map-ы и добавляем индексы token-ов с одинаковыми номерами в массив brackets
	for count := range mapFirst {
		brackets = append(brackets, []int{mapFirst[count], mapSecond[count]})
	}

	// Проверяем скобки на пустоту, и если скобка пустая то добавляем её в список emptyBrackets
	for _, bracket := range brackets {
		if len(tokens[bracket[0]+1:bracket[1]]) == 0 {
			emptyBrackets = append(emptyBrackets, fmt.Sprintf("[%d:%d]", bracket[0], bracket[1]))
		}
	}
	// Добавляем в resultError ошибки о пустых скобках
	if len(emptyBrackets) > 0 {
		resultError += fmt.Sprintf("%v, expression[%v]\n", errors.ErrEmptyBracket, strings.Join(emptyBrackets, ","))
	}
	// Если resultError не пустая то возвращаем ошибку, иначе nil
	if resultError != "" {
		return fmt.Errorf("%v", resultError)
	}

	return nil
}

// Функция проверяет является ли token числом
func IsInteger(token string) bool {
	_, err := strconv.ParseFloat(token, 64)
	return err == nil
}

// Функция проверяет является ли token допустимой математической операцией
func IsOperation(token string) bool {
	if token == "+" || token == "/" || token == "-" || token == "*" {
		return true
	}
	return false
}

// Функция проверяет является ли token скобкой
func IsBracket(token string) bool {
	if token == "(" || token == ")" {
		return true
	}
	return false
}