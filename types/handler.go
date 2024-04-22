package types

type IHandler interface {
	OnResult(result interface{})
	OnError(err error)
	OnComplete(result interface{})
}
