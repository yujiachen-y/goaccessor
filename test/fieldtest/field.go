// fieldtest contains only cases we support. Other cases like the following
// which refers to other package is beyond the scope of this package.
//
//	var a struct_type_from_other_package
//	var b = struct_type_from_other_package{}
//
// Please note, field accessors are only supported for struct types. And we
// don't support nil or concurrent access checking.
package fieldtest

type Normal struct {
	a, B int
	c    *int64
	d    map[string]complex128
}

//go:generate go run ../../. -t normal1 -a -f -p normal1
var normal1 *Normal

//go:generate go run ../../. -t normal2 -a -f -p normal2
var normal2 = &Normal{}

type Generic[T comparable, U any, V any] struct {
	t T
	u *U
	v V
	m map[T]V
	s []U
}

//go:generate go run ../../. -t generic1 -a -f -p generic1
var generic1 Generic[int, float64, string]

//go:generate go run ../../. -t generic2 -a -f -p generic2
var generic2 *Generic[float64, string, int]

//go:generate go run ../../. -t generic3 -a -f -p generic3
var generic3 = &Generic[string, int, float64]{}

//go:generate go run ../../. -t generic4 -a -f -p generic4
var generic4 = Generic[string, int, float64]{}

//go:generate go run ../../. -t anonymous1 -a -f -p anonymous1
var anonymous1 = &struct {
	quite, difficult string
}{}

//go:generate go run ../../. -t anonymous2 -a -f -p anonymous2
var anonymous2 = struct {
	quite, difficult string
}{}

//go:generate go run ../../. -t anonymous3 -a -f -p anonymous3
//go:generate go run ../../. -t anonymous4 -a -f -p anonymous4
var anonymous3, anonymous4 struct {
	quite, difficult string
}

//go:generate go run ../../. -t anonymous5 -a -f -p anonymous5
var anonymous5 *struct {
	quite, difficult string
}

type Pure struct {
	pure int
}

//go:generate go run ../../. -t pure -a -pg -f
var pure Pure
