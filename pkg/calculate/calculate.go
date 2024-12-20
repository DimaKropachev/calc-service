package calculate

import (
	"strconv"

	"github.com/DimaKropachev/calculate-web-server/pkg/errors"
)

var priority map[string]int = map[string]int{
	"*": 2,
	"/": 2,
	"+": 1,
	"-": 1,
}

// Создаем структуру со стеком для математических операций
type stackOperations struct {
	stackOper []string
}

func (stOp *stackOperations) Pop() string {
	last := stOp.stackOper[len(stOp.stackOper)-1]
	stOp.stackOper = stOp.stackOper[:len(stOp.stackOper)-1]
	return last
}

func (stOp *stackOperations) Top() string {
	return stOp.stackOper[len(stOp.stackOper)-1]
}

func (stOp *stackOperations) Append(element string) {
	stOp.stackOper = append(stOp.stackOper, element)
}

// Создаем структуру со стеком для чисел
type stackNumbers struct {
	stackNum []string
}

func (stNum *stackNumbers) Pop() (string, string) {
	last1 := stNum.stackNum[len(stNum.stackNum)-1]
	stNum.stackNum = stNum.stackNum[:len(stNum.stackNum)-1]
	last2 := stNum.stackNum[len(stNum.stackNum)-1]
	stNum.stackNum = stNum.stackNum[:len(stNum.stackNum)-1]
	return last1, last2
}

func (stNum *stackNumbers) Top() string {
	return stNum.stackNum[len(stNum.stackNum)-1]
}

func (stNum *stackNumbers) Append(element string) {
	stNum.stackNum = append(stNum.stackNum, element)
}

func Calc(expression string) (float64, error) {
	// Создаем 2 стека, один для чисел, другой для операций
	Numbers := &stackNumbers{}
	Operations := &stackOperations{}

	// Разбиваем выражение на токены
	tokens := GetTokens(expression)

	// Отправляем выражение на проверку
	err := CheckExpression(tokens)
	if err != nil {
		return 0, err
	}

	// Перебираем токены в цикле
	for _, token := range tokens {
		// Проверяем является ли токен числовой операцией или скобками
		if token == "+" || token == "-" || token == "*" || token == "/" || token == "(" || token == ")" {
			// Проверяем пустой ли стек с операциями
			if len(Operations.stackOper) == 0 {
				Operations.Append(token)
			} else {

				// Если токен является открывающейся скобкой то мы должны его просто добавить в стек с операциями, чтобы после он был для нас сигналом
				if token == "(" {
					Operations.Append(token)
					continue
				}

				// Если токен является закрывающейся скобкой то нам нужно сосчитать всё что находится до открывающей скобки
				if token == ")" {
					for {
						topOperation := Operations.Top()
						// Если в стеке дошли до открывающейся скобки, то мы должны закончить цикл
						if topOperation == "(" {
							_ = Operations.Pop()
							break
						}
						// Берем 2 числа из стека чисел и знак операции из стека операций и считаем
						num1, num2 := Numbers.Pop()
						operation := Operations.Pop()
						res, err := Calculate(operation, num1, num2)
						if err != nil {
							return 0, err
						}
						// Результат вычислений записываем в стек с числами
						Numbers.Append(res)
					}
					continue
				}

				topOperation := Operations.Top()

				// Проверяем является ли приоритет операции, которую мы сейчас обрабатываем больше операции, которая находится на вершине стека операций
				if priority[token] > priority[topOperation] {
					Operations.Append(token)
					continue
				} else {
					// Если нет, то мы должны сосчитать всё до того, пока либо стек операций не станет пустым, либо приоритет верхней операции стека не будет меньше приоритета операции, которую мы сейчас обрабатываем
					for {
						if len(Operations.stackOper) == 0 || priority[Operations.Top()] < priority[token] {
							break
						}
						// Берем 2 числа из стека чисел и знак операции из стека операций и считаем
						num1, num2 := Numbers.Pop()
						operation := Operations.Pop()
						res, err := Calculate(operation, num1, num2)
						if err != nil {
							return 0, err
						}
						// Результат вычислений записываем в стек с числами
						Numbers.Append(res)
					}
					// После окончания вычислений записываем операцию, которую мы обрабатываем в стек
					Operations.Append(token)
					continue
				}
			}
		} else {
			// Если токен не является математической операцией, то мы записываем его в стек чисел
			Numbers.Append(token)
		}
	}

	// Если после всех вычислений в стеке операций еще есть значения, то мы должны их досчитать
	for {
		if len(Operations.stackOper) == 0 {
			break
		}
		num1, num2 := Numbers.Pop()
		operation := Operations.Pop()
		res, err := Calculate(operation, num1, num2)
		if err != nil {
			return 0, err
		}
		Numbers.Append(res)
	}

	// Преобразуем получившийся результат в вещественное число
	result, err := strconv.ParseFloat(Numbers.stackNum[0], 64)
	if err != nil {
		return 0, errors.ErrConvertString
	}

	return result, nil
}

// Функция возвращает разбитую на числа и математические операции строку в виде массива
func GetTokens(expression string) []string {
	result := []string{}

	var num string
	// Перебираем строку по символам в цикле
	for _, r := range expression {
		sym := string(r)
		if sym == " " {
			// Если дошли до пробела, то если в переменной num есть значения, то добавляем его в массив
			if num != "" {
				result = append(result, num)
				// После добавления обновляем переменную num
				num = ""
			}
			continue
		}
		// Проверяем является ли символ строки знаком математического выражения
		if sym == "+" || sym == "-" || sym == "*" || sym == "/" || sym == "(" || sym == ")" {
			// Если да, то проверяем есть ли значение в переменной num. Если есть, то добавляем его в массив
			if num != "" {
				result = append(result, num)
				// После добавления обновляем переменную num
				num = ""
			}
			// Добавляем символ операции в массив
			result = append(result, sym)
			continue
		} else {
			// Если символ не является знаком математической операции, то добавляем его к переменной num (если этот символ не пробел)
			if sym != " " {
				num += sym
			}
		}
	}
	// После прохода по всем символам проверяем есть ли значения в переменной num
	if num != "" {
		// Если есть, то добавляем их в массив
		result = append(result, num)
	}

	return result
}

// Функция считает простейшие математические операции
func Calculate(operation, num1, num2 string) (string, error) {
	num1Float, err := strconv.ParseFloat(num1, 64)
	if err != nil {
		return "", errors.ErrConvertString
	}
	num2Float, err := strconv.ParseFloat(num2, 64)
	if err != nil {
		return "", errors.ErrConvertString
	}
	var result string

	if operation == "*" {
		result = strconv.FormatFloat(num1Float*num2Float, 'f', -1, 64)
	} else if operation == "/" {
		if num1Float == 0 {
			return "", errors.ErrDivisionByZero
		}
		result = strconv.FormatFloat(num2Float/num1Float, 'f', -1, 64)
	} else if operation == "+" {
		result = strconv.FormatFloat(num1Float+num2Float, 'f', -1, 64)
	} else if operation == "-" {
		result = strconv.FormatFloat(num2Float-num1Float, 'f', -1, 64)
	}

	return result, nil
}
