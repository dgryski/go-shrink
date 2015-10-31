package shrink

import (
	"math/rand"
	"reflect"
	"testing"
	"testing/quick"
)

func listshrinker(v reflect.Value, tactic int) (rv reflect.Value, err error) {

	switch tactic {
	case 0:
		if v.Len() < 2 {
			return reflect.Value{}, ErrDeadEnd
		}

		nv := reflect.MakeSlice(v.Type(), v.Len()/2, v.Len()/2)
		reflect.Copy(nv, v)
		return nv, nil
	case 1:
		if v.Len() < 2 {
			return reflect.Value{}, ErrDeadEnd
		}

		nv := reflect.MakeSlice(v.Type(), v.Len()/2, v.Len()/2)
		reflect.Copy(nv, v.Slice(v.Len()/2, v.Len()))
		return nv, nil
	case 2:
		if v.Len() < 2 {
			return reflect.Value{}, ErrDeadEnd
		}

		nv := reflect.MakeSlice(v.Type(), v.Len()-1, v.Len()-1)
		reflect.Copy(nv, v)
		return nv, nil
	case 3:
		if v.Len() < 2 {
			return reflect.Value{}, ErrDeadEnd
		}

		nv := reflect.MakeSlice(v.Type(), v.Len()-1, v.Len()-1)
		reflect.Copy(nv, v.Slice(1, v.Len()))
		return nv, nil
	case 4:
		if v.Len() < 2 {
			return reflect.Value{}, ErrDeadEnd
		}

		if !isNumeric(v.Type().Elem()) {
			return reflect.Value{}, ErrDeadEnd
		}

		var allzeros bool

		nv := reflect.MakeSlice(v.Type(), v.Len(), v.Len())
		reflect.Copy(nv, v)

		for i := 0; i < nv.Len(); i++ {
			nn := nv.Index(i)
			allzeros = allzeros || isZero(nn)
			div2(nn)
		}
		if allzeros {
			return reflect.Value{}, ErrDeadEnd
		}
		return nv, nil
	}

	return reflect.Value{}, ErrNoMoreTactics
}

func isNumeric(v reflect.Type) bool {
	switch v.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Float32, reflect.Float64:
		return true
	}

	return false
}

func isZero(v reflect.Value) bool {

	switch v.Type().Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return v.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return v.Float() == 0
	}

	panic("isZero: not a numeric type")
}

func div2(v reflect.Value) {

	switch v.Type().Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		v.SetInt(v.Int() / 2)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		v.SetUint(v.Uint() / 2)
	case reflect.Float32, reflect.Float64:
		v.SetFloat(v.Float() / 2)
	default:
		panic("div2: not a numeric type")
	}
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
		s, serr := Shrink(checkReverse, check.In, []Shrinker{listshrinker})
		t.Logf("shrunk: %v serr=%+v", s[0].Interface(), serr)
	}
}
