package errors

import "errors"

var (
	ErrInvalidToken         error = errors.New("в выражении не должы присутствовать символы кроме цифр и допустимых знаков математических операций")
	ErrBracketsInExpression error = errors.New("в выражении неверно расставлены скобки")
	ErrBeforeFirstBracket   error = errors.New("перед открывающейся скобкой ожидалась математическая операция")
	ErrAfterFirstBracket    error = errors.New("после открывающейся скобки ожидалось число")
	ErrBeforeSecondBracket  error = errors.New("перед закрывающей скобкой одидалось число")
	ErrAfterSecondBracket   error = errors.New("после закрываюшейся скобки ожидалась математическая операция")
	ErrOperationsContract   error = errors.New("не допускается использование 2 и более математических операций подряд")
	ErrIntegersContract     error = errors.New("не допускается использование 2 и более чисел подряд")
	ErrEmptyBracket         error = errors.New("пустая скобка")
	ErrFirstOperation       error = errors.New("выражение не должно начинаться с математической операции")
	ErrLastOperation        error = errors.New("выражение не должно заканчиваться математической операцией")
	ErrDivisionByZero       error = errors.New("деление на 0")
	ErrEmptyExpression      error = errors.New("expression is empty")
	ErrConvertString        error = errors.New("не удалось преобразовать строку в число")
)
