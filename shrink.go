package shrink

import (
	"errors"
	"math"
	"reflect"
)

type Shrinker func(reflect.Value, int) (reflect.Value, error)

var (
	ErrPassingTest = errors.New("shrink: can't shrink passing test")
	ErrBadFunction = errors.New("shrink: bad function type")

	ErrDeadEnd       = errors.New("shrink: dead end")
	ErrNoMoreTactics = errors.New("shrink: no more tactics")
)

func Shrink(f interface{}, in []interface{}, shrinkers []Shrinker) ([]reflect.Value, error) {

	args := make([]reflect.Value, len(in))

	for i, a := range in {
		args[i] = reflect.ValueOf(a)
	}

	rval := reflect.ValueOf(f)
	if rval.Kind() != reflect.Func || rval.Type().NumOut() != 1 || rval.Type().Out(0).Kind() != reflect.Bool {
		return nil, ErrBadFunction
	}

	out := rval.Call(args)
	if out[0].Bool() {
		return nil, ErrPassingTest
	}

	progress := true
	for progress {
		progress = false
		for i := range args {
			ok, err := shrinkOne(rval, args, i, shrinkers[i])
			if err != nil {
				return nil, err
			}

			progress = ok
		}
	}

	return args, nil
}

const MaxTactics = math.MaxUint32

func shrinkOne(f reflect.Value, args []reflect.Value, ai int, shrinker Shrinker) (bool, error) {

	if shrinker == nil {
		return false, nil
	}

	for tactic := 0; tactic < MaxTactics; tactic++ {
		cur := args[ai]
		nv, err := shrinker(cur, tactic)
		if err != nil {
			switch err {
			case ErrNoMoreTactics:
				return false, nil
			case ErrDeadEnd:
				continue
			default:
				return false, err
			}
		}

		args[ai] = nv

		fret := f.Call(args)[0].Bool()
		if !fret {
			return true, nil
		}
		args[ai] = cur
	}

	return false, nil
}
