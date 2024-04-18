package handler

type BaseHandler struct {
	onResult   func(result interface{})
	onError    func(err error)
	onComplete func()
}

func (b *BaseHandler) OnResult(result interface{}) {
	b.onResult(result)
}

func (b *BaseHandler) OnError(err error) {
	b.onError(err)
}

func (b *BaseHandler) OnComplete() {
	b.onComplete()
}

func NewBaseHandler(onResult func(result interface{}), onError func(err error), onComplete func()) *BaseHandler {
	return &BaseHandler{
		onResult:   onResult,
		onError:    onError,
		onComplete: onComplete,
	}
}
