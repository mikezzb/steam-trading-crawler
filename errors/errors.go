package errors

import "errors"

var ErrItemNotFound = errors.New("item not found")
var ErrMarketNotFound = errors.New("market not found")
var ErrTaskNotFound = errors.New("task not found")
var ErrCrawlerManuallyStopped = errors.New("crawler manually stopped")
var ErrNoCookies = errors.New("no cookies")

const (
	SafeInvalidPrice = "invalid price"
)
