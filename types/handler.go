package types

type Handler interface {
	OnResult(result interface{})
	OnError(err error)
	OnComplete(result interface{})
}
