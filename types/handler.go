package types

type Handler struct {
	OnResult   func(result interface{})
	OnError    func(err error)
	OnComplete func()
}
