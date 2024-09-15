package myErrors

import "errors"

var (
	ErrUserNotFound   = errors.New("пользователь не существует")
	ErrBadRequest     = errors.New("некорректный запрос")
	ErrForbidden      = errors.New("недостаточно прав для выполнения запроса")
	ErrTenderNotFound = errors.New("тендер не найден")
	ErrBidNotFound    = errors.New("предложение не найдено")
	ErrInternal       = errors.New("внутренняя ошибка сервера")
)
