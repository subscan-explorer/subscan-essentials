package result

type Result struct {
	Ok  interface{}
	Err error
}

func (r *Result) IsOk() bool {
	return r.Err == nil
}

func (r *Result) IsErr() bool { return !r.IsOk() }

func (r *Result) Map(fn func(interface{}) interface{}) *Result {
	if r.IsOk() {
		return &Result{Ok: fn(r.Ok)}
	}
	return &Result{Err: r.Err}
}

func (r *Result) MapErr(fn func(err error) error) *Result {
	if r.IsOk() {
		return &Result{Ok: r.Ok}
	}
	return &Result{Err: fn(r.Err)}
}

func (r *Result) And(res *Result) *Result {
	if r.IsOk() {
		return res
	}
	return &Result{Err: r.Err}
}

func (r *Result) AndThen(fn func(interface{}) *Result) *Result {
	if r.IsOk() {
		return fn(r.Ok)
	}
	return &Result{Err: r.Err}
}

func (r *Result) Or(res *Result) *Result {
	if r.IsOk() {
		return &Result{Ok: r.Ok}
	}
	return res
}

func (r *Result) OrElse(fn func(err error) *Result) *Result {
	if r.IsOk() {
		return &Result{Ok: r.Ok}
	}
	return fn(r.Err)
}

func (r *Result) Unwrap() interface{} {
	if r.IsOk() {
		return r.Ok
	}
	panic("err: try unwrap Result(Err) into Ok")
}

func (r *Result) UnwrapOr(v interface{}) interface{} {
	if r.IsOk() {
		return r.Ok
	}
	return v
}

func (r *Result) UnwrapOrElse(fn func(err error) interface{}) interface{} {
	if r.IsOk() {
		return r.Ok
	}
	return fn(r.Err)
}
