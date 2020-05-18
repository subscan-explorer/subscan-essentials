package result

// Wrap a function with recover and do nothing.
func Wrap(fn func()) {
	defer func() {
		if r := recover(); r != nil {
			return
		}
	}()
	fn()
}

func WrapCallback(fn func(), handle func(x interface{})) {
	defer func() {
		if x := recover(); x != nil {
			handle(x)
		}
	}()
	fn()
}
