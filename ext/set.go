package ext

import (
	"encoding/json"
	"fmt"
	"maps"
)

// Set 定义一个泛型集合类型 Set，底层使用 map 实现
// description: 线程不安全
type Set[E comparable] map[E]Unit

// Set_ 创建一个指定容量的集合
func Set_[E comparable](cap int) Set[E] {
	return make(map[E]Unit, cap)
}

// SetOf 创建一个包含指定元素的集合
func SetOf[E comparable](es ...E) Set[E] {
	s := make(map[E]Unit, len(es)) // 根据元素数量分配容量
	for _, e := range es {
		s[e] = Unit{} // 将每个元素添加到集合中
	}
	return s
}

// SetOfUnmarshalJSON 通过JSON 反序列化方法 创建一个集合
func SetOfUnmarshalJSON[E comparable](data []byte) (Set[E], error) {
	vec := make(Vec[E], 0) // 临时存储 JSON 数据
	err := json.Unmarshal(data, &vec)
	if err != nil {
		return nil, err
	}
	return SetOf[E](vec...), nil // 将切片中的元素存入集合
}

// ForEach 遍历集合，对每个元素执行给定的函数
func (s Set[E]) ForEach(fn func(E)) {
	for e := range s {
		fn(e)
	}
}

// Len 获取集合的元素数量
func (s Set[E]) Len() int {
	return len(s)
}

// Empty 判断集合是否为空
func (s Set[E]) Empty() bool {
	return len(s) == 0
}

// Add 向集合中插入元素并返回自身
func (s Set[E]) Add(es ...E) Set[E] {
	for _, e := range es {
		s[e] = Unit{}
	}
	return s
}

// Remove 从集合中移除一个元素
func (s Set[E]) Remove(element E) {
	delete(s, element)
}

// Contains 判断集合是否包含某个元素
func (s Set[E]) Contains(element E) bool {
	_, ok := s[element]
	return ok
}

// Or 求两个集合的并集
func (s Set[E]) Or(other Set[E]) Set[E] {
	s0, s1 := s, other
	// 优化：选择较大的集合作为基准
	if s0.Len() < s1.Len() {
		s0, s1 = s1, s0
	}
	s_ := maps.Clone(s0) // 克隆较大的集合
	for e := range s1 {
		s_[e] = Unit{} // 将另一个集合的元素添加到结果集合
	}
	return s_
}

// And 求两个集合的交集
func (s Set[E]) And(other Set[E]) Set[E] {
	s0, s1 := s, other
	// 优化：选择较小的集合作为迭代基准
	if s0.Len() < s1.Len() {
		s0, s1 = s1, s0
	}
	s_ := Set_[E]((s1.Len() + 1) / 2) // 分配较小集合的一半容量
	for e := range s1 {
		if s0.Contains(e) {
			s_[e] = Unit{} // 只保留同时存在的元素
		}
	}
	return s_
}

// Sub 求两个集合的差集
func (s Set[E]) Sub(other Set[E]) Set[E] {
	if s.Len() < other.Len()*2 {
		// 优化：当当前集合较小时，逐一检查
		s_ := Set_[E]((s.Len() + 1) / 2)
		for e := range s {
			if !other.Contains(e) {
				s_[e] = Unit{}
			}
		}
		return s_
	} else {
		// 克隆当前集合并移除另一个集合中的元素
		s_ := maps.Clone(s)
		for e := range other {
			s_.Remove(e)
		}
		return s_
	}
}

// Xor 求两个集合的对称差集
func (s Set[E]) Xor(other Set[E]) Set[E] {
	s0, s1 := s, other
	// 优化：选择较大的集合作为基准
	if s0.Len() < s1.Len() {
		s0, s1 = s1, s0
	}
	s_ := maps.Clone(s0) // 克隆较大的集合
	for e := range s1 {
		if s_.Contains(e) {
			s_.Remove(e) // 如果两个集合都包含该元素，则移除
		} else {
			s_[e] = Unit{} // 如果仅在另一个集合中，则添加
		}
	}
	return s_
}

// ToVec 将集合转换为切片
func (s Set[E]) ToVec() Vec[E] {
	v := Vec_[E](len(s)) // 分配与集合大小相同的容量
	for e := range s {
		v.Append(e)
	}
	return v
}

// Clear 清空集合
func (s Set[E]) Clear() {
	clear(s)
}

// ToString 将集合转换为字符串表示
func (s Set[E]) ToString() string {
	return fmt.Sprintf("set%v", s.ToVec())
}

// MarshalJSON 实现 JSON 序列化方法
func (s Set[E]) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.ToVec()) // 将集合转换为切片再序列化
}
