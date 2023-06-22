package vartest

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
