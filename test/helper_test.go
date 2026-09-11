package test

// panics reports whether fn panics.
func panics(fn func()) (p bool) {
	defer func() {
		if r := recover(); r != nil {
			p = true
		}
	}()
	fn()
	return
}
