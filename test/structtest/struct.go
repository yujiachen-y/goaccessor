// structtest contains only cases we support. Other cases like the following are
// related to type inference which is beyond the scope of this package.
//
//    type equal = struct_type_from_other_package_or_current_package
//    type alias struct_type_from_other_package_or_current_package
//
package structtest

type S struct{}

//go:generate go run ../../. -t Normal -a -i a,B,c,D,e,f
type Normal struct {
	a, B int
	c    *int64
	D    chan float64
	e    map[string]complex128
	f    []*S
	s    S
}

// SetE should not be generated
func (n *Normal) SetE(e map[string]complex128) {
	n.e = e
}

//go:generate go run ../../. -t Generic -a -e a
type Generic[T, U any] struct {
	a, B int
	c    *int64
	D    chan float64
	e    map[string]complex128
	f    []*S
	s    S
	t    T
	u    U
}

// GetB should not be generated
func (g *Generic[_, _]) GetB() int {
	return g.B
}

// SetU should not be generated
func (g *Generic[T, U]) SetU(u U) {
	g.u = u
}
