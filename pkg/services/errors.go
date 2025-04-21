package services

import "errors"

// ErrCityNotFound повертається, коли вказане місто не знайдено
var ErrCityNotFound = errors.New("city not found")

// ErrDuplicateSubscription повертається, коли підписка вже існує
var ErrDuplicateSubscription = errors.New("duplicate subscription")

// ErrTokenNotFound повертається, коли токен підтвердження не знайдено
var ErrTokenNotFound = errors.New("token not found")

// ErrTokenExpired повертається, коли токен підтвердження прострочено
var ErrTokenExpired = errors.New("token expired")
