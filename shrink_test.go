package shrink

import (
	"math/rand"
	"reflect"
	"testing"
	"testing/quick"
)

func intlistshrinker(v reflect.Value, tactic int) (rv reflect.Value, err error) {

	var nv []int
	iv := v.Interface().([]int)

	switch tactic {
	case 0:
		if len(iv) < 2 {
			return reflect.Value{}, ErrDeadEnd
		}

		nv = make([]int, len(iv)/2)
		copy(nv, iv)
		return reflect.ValueOf(nv), nil
	case 1:

		if len(iv) < 2 {
			return reflect.Value{}, ErrDeadEnd
		}

		nv = make([]int, len(iv)/2)
		copy(nv, iv[len(nv):])
		return reflect.ValueOf(nv), nil
	case 2:

		if len(iv) < 2 {
			return reflect.Value{}, ErrDeadEnd
		}

		nv = make([]int, len(iv)-1)
		copy(nv, iv)
		return reflect.ValueOf(nv), nil
	case 3:

		if len(iv) < 2 {
			return reflect.Value{}, ErrDeadEnd
		}

		nv = make([]int, len(iv)-1)
		copy(nv, iv[1:])
		return reflect.ValueOf(nv), nil
	case 4:

		if len(iv) < 2 {
			return reflect.Value{}, ErrDeadEnd
		}

		var allzeros bool

		nv = make([]int, len(iv))
		copy(nv, iv)
		for i, nn := range nv {
			allzeros = allzeros || nn == 0
			nv[i] = nn / 2
		}
		if allzeros {
			return reflect.Value{}, ErrDeadEnd
		}
		return reflect.ValueOf(nv), nil
	}

	return reflect.Value{}, ErrNoMoreTactics
}

type smallints []int

func (s smallints) Generate(rand *rand.Rand, size int) reflect.Value {
	nv := make([]int, size)

	for i := range nv {
		nv[i] = rand.Intn(1000)
	}

	return reflect.ValueOf(nv)
}

func buggyReverse(l []int) []int {

	for i := 0; i < len(l); i++ {
		l[i], l[len(l)-i-1] = l[len(l)-i-1], l[i]
	}

	return l
}

func TestQuick(t *testing.T) {

	checkReverse := func(l smallints) bool {
		l1 := append([]int(nil), l...)
		buggyReverse(l1)
		for i := range l {
			if l1[i] != l[len(l)-i-1] {
				return false
			}
		}
		return true
	}

	if err := quick.Check(checkReverse, nil); err != nil {
		t.Error(err)
		check := err.(*quick.CheckError)
		s, serr := Shrink(checkReverse, check.In, []Shrinker{intlistshrinker})
		t.Logf("shrunk: %v serr=%+v", s[0].Interface(), serr)
	}
}
