package ext

type VecM2[T0, T1 any] T2[Vec[T0], Vec[T1]]

type VecM3[T0, T1, T2 any] T3[Vec[T0], Vec[T1], Vec[T2]]

func VecM2_[T0, T1 any](vec0 Vec[T0], vec1 Vec[T1]) VecM2[T0, T1] {
	return VecM2[T0, T1]{vec0, vec1}
}

func (m VecM2[T0, T1]) ForEach(fn func(T2[T0, T1])) {
	for i := range m.Len() {
		fn(T2_(m.V0[i], m.V1[i]))
	}
}

func (m VecM2[T0, T1]) ForEachWhile(fn func(T2[T0, T1]) bool) {
	for i := range m.Len() {
		if !fn(T2_(m.V0[i], m.V1[i])) {
			return
		}
	}
}

func (m VecM2[T0, T1]) Len() int {
	return min(m.V0.Len(), m.V1.Len())
}

func (m VecM2[T0, T1]) Empty() bool {
	return m.Len() == 0
}

func (m VecM2[T0, T1]) ToVec() Vec[T2[T0, T1]] {
	rs := Vec_[T2[T0, T1]](m.Len())
	m.ForEach(rs.Append)
	return rs
}

func VecM3_[T0, T1, T2 any](vec0 Vec[T0], vec1 Vec[T1], vec2 Vec[T2]) VecM3[T0, T1, T2] {
	return VecM3[T0, T1, T2]{vec0, vec1, vec2}
}

func (m VecM3[T0, T1, T2]) ForEach(fn func(T3[T0, T1, T2])) {
	for i := range m.Len() {
		fn(T3_(m.V0[i], m.V1[i], m.V2[i]))
	}
}

func (m VecM3[T0, T1, T2]) ForEachWhile(fn func(T3[T0, T1, T2]) bool) {
	for i := range m.Len() {
		if !fn(T3_(m.V0[i], m.V1[i], m.V2[i])) {
			break
		}
	}
}

func (m VecM3[T0, T1, T2]) Len() int {
	return min(m.V0.Len(), m.V1.Len(), m.V2.Len())
}

func (m VecM3[T0, T1, T2]) Empty() bool {
	return m.Len() == 0
}

func (m VecM3[T0, T1, T2]) ToVec() Vec[T3[T0, T1, T2]] {
	rs := Vec_[T3[T0, T1, T2]](m.Len())
	m.ForEach(rs.Append)
	return rs
}
