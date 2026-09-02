package ext

func Try(err error) {
	if err != nil {
		panic(err)
	}
}

func Try1[T0 any](t0 T0, e error) T0 {
	Try(e)
	return t0
}

func Try2[T0, T1 any](t0 T0, t1 T1, e error) (T0, T1) {
	Try(e)
	return t0, t1
}

func Try3[T0, T1, T2 any](t0 T0, t1 T1, t2 T2, e error) (T0, T1, T2) {
	Try(e)
	return t0, t1, t2
}

func Try4[T0, T1, T2, T3 any](t0 T0, t1 T1, t2 T2, t3 T3, e error) (T0, T1, T2, T3) {
	Try(e)
	return t0, t1, t2, t3
}

func Try5[T0, T1, T2, T3, T4 any](t0 T0, t1 T1, t2 T2, t3 T3, t4 T4, e error) (T0, T1, T2, T3, T4) {
	Try(e)
	return t0, t1, t2, t3, t4
}

func Try6[T0, T1, T2, T3, T4, T5 any](
	t0 T0, t1 T1, t2 T2, t3 T3, t4 T4, t5 T5, e error) (T0, T1, T2, T3, T4, T5) {
	Try(e)
	return t0, t1, t2, t3, t4, t5
}

func Try7[T0, T1, T2, T3, T4, T5, T6 any](
	t0 T0, t1 T1, t2 T2, t3 T3, t4 T4, t5 T5, t6 T6, e error) (T0, T1, T2, T3, T4, T5, T6) {
	Try(e)
	return t0, t1, t2, t3, t4, t5, t6
}

func Try8[T0, T1, T2, T3, T4, T5, T6, T7 any](
	t0 T0, t1 T1, t2 T2, t3 T3, t4 T4, t5 T5, t6 T6, t7 T7, e error) (T0, T1, T2, T3, T4, T5, T6, T7) {
	Try(e)
	return t0, t1, t2, t3, t4, t5, t6, t7
}

func Try9[T0, T1, T2, T3, T4, T5, T6, T7, T8 any](
	t0 T0, t1 T1, t2 T2, t3 T3, t4 T4, t5 T5, t6 T6, t7 T7, t8 T8, e error) (T0, T1, T2, T3, T4, T5, T6, T7, T8) {
	Try(e)
	return t0, t1, t2, t3, t4, t5, t6, t7, t8
}
