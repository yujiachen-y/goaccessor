package vartest

import (
	"testing"

	"github.com/yjc567/goaccessor/test/utils"
)

//go:generate goaccessor -t i -a
var i int

//go:generate goaccessor -t A,B,C -a
var A, B, C float64

//go:generate goaccessor -t k -a
var k = 0

//go:generate goaccessor -t x,y -a
var x, y float32 = -1, -2

//go:generate goaccessor -t j,u,v,s -a
var (
	j       int
	u, v, s = 2.0, 3.0, "bar"
)

type Struct struct{}

//go:generate goaccessor -t s1,s2 -a
var (
	s1 = []int{1, 2, 3}
	s2 []Struct
)

//go:generate goaccessor -t m1,m2 -a
var (
	m1 = map[string]struct{}{}
	m2 map[int]*Struct
)

//go:generate goaccessor -t hello,world -a
var hello = &struct{ hello string }{}
var world *Struct

/*
Supporting type inference is beyond the scope.

var re, im = complexSqrt(-1)
var _, found = entries[name]  // map lookup; only interested in "found"
var d = math.Sin(0.5)  // d is float64
var t, ok = x.(T)      // t is T, ok is bool

var a = pkg.Struct.Field
var a = 1 + 2
*/

func TestGetVar(t *testing.T) {
	for _, verifier := range []utils.Verifier{
		utils.NewGetterVerifier(GetI, i),
		utils.NewGetterVerifier(GetA, A),
		utils.NewGetterVerifier(GetB, B),
		utils.NewGetterVerifier(GetC, C),
		utils.NewGetterVerifier(GetK, k),
		utils.NewGetterVerifier(GetX, x),
		utils.NewGetterVerifier(GetY, y),
		utils.NewGetterVerifier(GetJ, j),
		utils.NewGetterVerifier(GetU, u),
		utils.NewGetterVerifier(GetV, v),
		utils.NewGetterVerifier(GetS, s),
		utils.NewSliceGetterVerifier(GetS1, s1),
		utils.NewSliceGetterVerifier(GetS2, s2),
		utils.NewMapGetterVerifier(GetM1, m1),
		utils.NewMapGetterVerifier(GetM2, m2),
		utils.NewGetterVerifier(GetHello, hello),
		utils.NewGetterVerifier(GetWorld, world),
	} {
		if err := verifier(); err != nil {
			t.Errorf("got err: %s", err.Error())
		}
	}
}

func TestSetVar(t *testing.T) {
	for _, verifier := range []utils.Verifier{
		utils.NewSetterVerifier(&i, SetI, 5),
		utils.NewSetterVerifier(&A, SetA, 6.7),
		utils.NewSetterVerifier(&B, SetB, 7.8),
		utils.NewSetterVerifier(&C, SetC, 8.9),
		utils.NewSetterVerifier(&k, SetK, 9),
		utils.NewSetterVerifier(&x, SetX, 10.1),
		utils.NewSetterVerifier(&y, SetY, 11.12),
		utils.NewSetterVerifier(&j, SetJ, 12),
		utils.NewSetterVerifier(&u, SetU, 13.14),
		utils.NewSetterVerifier(&v, SetV, 14.15),
		utils.NewSetterVerifier(&s, SetS, "15.16"),
		utils.NewSliceSetterVerifier(&s1, SetS1, []int{1, 6, 1, 7}),
		utils.NewSliceSetterVerifier(&s2, SetS2, []Struct{{}}),
		utils.NewMapSetterVerifier(&m1, SetM1, map[string]struct{}{"m1": {}}),
		utils.NewMapSetterVerifier(&m2, SetM2, map[int]*Struct{17: {}, 18: {}}),
		utils.NewPointSetterVerifier(&hello, SetHello, struct{ hello string }{"world"}),
		utils.NewPointSetterVerifier(&world, SetWorld, Struct{}),
	} {
		if err := verifier(); err != nil {
			t.Errorf("got err: %s", err.Error())
		}
	}
}
