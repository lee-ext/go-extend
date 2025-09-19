package test

import (
	"fmt"
	"maps"
	"slices"
	"testing"

	. "github.com/lee-ext/go-extend/ext"
)

func TestFunc(t *testing.T) {
	vec0 := VecOf[int64](5, 4, 3, 2, 1)
	fmt.Printf("vec0: %v\n", vec0)

	vec1 := Map(vec0, func(t int64) float64 {
		return float64(t) * 10
	})
	fmt.Printf("vec1: %v\n", vec1)

	if v := VecOf(50.0, 40.0, 30.0, 20.0, 10.0); !slices.Equal(vec1, v) {
		t.Errorf("result is %v, expected %v", vec1, v)
	}

	set0 := IntactTo(vec0, Set_[int64])
	fmt.Printf("set0: %v\n", set0)

	if !set0.Contains(5) {
		t.Errorf("%v should contain 5", set0)
	}

	dict0 := IntactTo(vec0.ToIndexed(), Dict_[int, int64])
	fmt.Printf("dict0: %v\n", dict0)

	if d := map[int]int64{0: 5, 1: 4, 2: 3, 3: 2, 4: 1}; !maps.Equal(dict0, d) {
		t.Errorf("result is %v, expected %v", dict0, d)
	}

	vec3 := IntactTo(vec0.ToReverse(), Vec_[int64])
	fmt.Printf("vec3: %v\n", vec3)

	if v := VecOf[int64](1, 2, 3, 4, 5); !slices.Equal(vec3, v) {
		t.Errorf("result is %v, expected %v", vec3, v)
	}

	dict1 := MapTo(vec0, func(t int64) KV[int64, string] {
		return KV_(t, fmt.Sprintf("num:%d", t))
	}, Dict_[int64, string])
	fmt.Printf("dict1: %v\n", dict1)

	if d := map[int64]string{5: "num:5", 4: "num:4", 3: "num:3", 2: "num:2", 1: "num:1"}; //
	!maps.Equal(dict1, d) {
		t.Errorf("result is %v, expected %v", dict1, d)
	}

	group := VecOf(VecOf(1, 2, 3), VecOf(4, 5, 6), VecOf(7, 8, 9))
	vec4 := Flatten[int](group)
	fmt.Printf("vec4: %v\n", vec4)

	if v := VecOf(1, 2, 3, 4, 5, 6, 7, 8, 9); !slices.Equal(vec4, v) {
		t.Errorf("result is %v, expected %v", vec3, v)
	}

	deque0 := FlattenTo(group, Deque_[int])
	fmt.Printf("deque0: %v\n", deque0)

	if v := VecOf(1, 2, 3, 4, 5, 6, 7, 8, 9); !slices.Equal(deque0.ToVec(), v) {
		t.Errorf("result is %v, expected %v", deque0, v)
	}

	group1 := VecOf("hello", " ", "world")
	vec5 := FlatMap(group1, StrCastBytes)
	fmt.Printf("vec5 to str: %v\n", string(vec5))

	if v := VecOf[byte]('h', 'e', 'l', 'l', 'o', ' ', 'w', 'o', 'r', 'l', 'd'); //
	!slices.Equal(vec5, v) {
		t.Errorf("result is %v, expected %v", vec5, v)
	}

	set1 := FlatMapTo(group1, StrCastBytes, Set_[byte])
	fmt.Printf("set1: %v\n", set1)

	if s := SetOf[byte]('h', 'e', 'l', 'l', 'o', ' ', 'w', 'o', 'r', 'l', 'd'); //
	!maps.Equal(set1, s) {
		t.Errorf("result is %v, expected %v", set1, s)
	}

	vec6 := Filter(vec0, func(t int64) bool {
		return t > 2
	})
	fmt.Printf("vec6: %v\n", vec6)

	if v := VecOf[int64](5, 4, 3); !slices.Equal(vec6, v) {
		t.Errorf("result is %v, expected %v", vec6, v)
	}

	deque1 := FilterTo(vec0, func(t int64) bool {
		return t > 2
	}, Deque_[int64])
	fmt.Printf("deque1: %v\n", deque1)

	if v := VecOf[int64](5, 4, 3); !slices.Equal(deque1.ToVec(), v) {
		t.Errorf("result is %v, expected %v", deque1, v)
	}

	ptrVec := Map(vec0,
		func(t int64) *int64 {
			if t == 3 {
				return nil
			}
			return &t
		})

	vec7 := FilterMap(ptrVec, PtrToOpt)
	fmt.Printf("vec7: %v\n", vec7)

	if v := VecOf[int64](5, 4, 2, 1); !slices.Equal(vec7, v) {
		t.Errorf("result is %v, expected %v", vec7, v)
	}

	set2 := FilterMapTo(ptrVec, PtrToOpt, Set_[int64])
	fmt.Printf("set2: %v\n", set2)

	if s := SetOf[int64](5, 4, 2, 1); !maps.Equal(set2, s) {
		t.Errorf("result is %v, expected %v", set2, s)
	}

	sum := Reduce(vec0, 0, func(l int64, r int64) int64 {
		return l + r
	})
	fmt.Printf("sum: %v\n", sum)

	if s := int64(1 + 2 + 3 + 4 + 5); sum != s {
		t.Errorf("result is %v, expected %v", sum, s)
	}

	type Student struct {
		Id   int
		Name string
		Age  int
	}

	stuVec := VecOf(Student{
		Id:   1,
		Name: "Tom",
		Age:  15,
	}, Student{
		Id:   2,
		Name: "Jerry",
		Age:  15,
	}, Student{
		Id:   3,
		Name: "Spike",
		Age:  20,
	})

	stuDict0 := ToDict(stuVec, func(t Student) int {
		return t.Id
	})
	fmt.Printf("stuDict0: %v\n", stuDict0)

	if d := map[int]Student{1: {
		Id:   1,
		Name: "Tom",
		Age:  15,
	}, 2: {
		Id:   2,
		Name: "Jerry",
		Age:  15,
	}, 3: {
		Id:   3,
		Name: "Spike",
		Age:  20,
	}}; //
	!maps.Equal(stuDict0, d) {
		t.Errorf("result is %v, expected %v", stuDict0, d)
	}

	stuDict1 := VToDict(stuVec, func(t Student) (int, string) {
		return t.Id, t.Name
	})
	fmt.Printf("stuDict1: %v\n", stuDict1)

	if d := map[int]string{1: "Tom", 2: "Jerry", 3: "Spike"}; //
	!maps.Equal(stuDict1, d) {
		t.Errorf("result is %v, expected %v", stuDict1, d)
	}

	stuGroup0 := VGroupBy(stuVec, func(t Student) (int, string) {
		return t.Age, t.Name
	})
	fmt.Printf("stuGroup0: %v\n", stuGroup0)

	if g := map[int]Vec[string]{15: VecOf("Tom", "Jerry"), 20: VecOf("Spike")}; //
	!maps.EqualFunc(stuGroup0, g, slices.Equal[Vec[string]]) {
		t.Errorf("result is %v, expected %v", stuGroup0, g)
	}

}
