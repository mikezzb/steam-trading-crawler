package handler

type BaseHandler struct {
	onResult   func(result interface{})
	onError    func(err error)
	onComplete func(result interface{})
}

func (b *BaseHandler) OnResult(result interface{}) {
	b.onResult(result)
}

func (b *BaseHandler) OnError(err error) {
	b.onError(err)
}

func (b *BaseHandler) OnComplete(result interface{}) {
	b.onComplete(result)
}

func NewBaseHandler(onResult func(result interface{}), onError func(err error), onComplete func(result interface{})) *BaseHandler {
	return &BaseHandler{
		onResult:   onResult,
		onError:    onError,
		onComplete: onComplete,
	}
}
