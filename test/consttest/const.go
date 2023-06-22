package consttest

//go:generate go run ../../. -t Pi -g
const Pi float64 = 3.14159265358979323846

//go:generate go run ../../. -t zero -g
const zero = 0.0 // untyped floating-point constant

//go:generate go run ../../. -t a,b,c -g
const a, b, c = 3, 4, "foo" // a = 3, b = 4, c = "foo", untyped integer and string constants

//go:generate go run ../../. -t u,v -g
const u, v float32 = 0, 3 // u = 0.0, v = 3.0

//go:generate go run ../../. -t size,eof -g
const (
	size int64 = 1024
	eof        = -1 // untyped integer constant
)

//go:generate go run ../../. -t Sunday,Monday,Tuesday,Wednesday,Thursday,Friday,Partyday,numberOfDays -g
const (
	Sunday = iota
	Monday
	Tuesday
	Wednesday
	Thursday
	Friday
	Partyday
	numberOfDays // this constant is not exported
)

/*
Maybe we should support the (constant expression)[https://go.dev/ref/spec#Constant_expressions]
 but not now.

const a = 2 + 3.0          // a == 5.0   (untyped floating-point constant)
const b = 15 / 4           // b == 3     (untyped integer constant)
const c = 15 / 4.0         // c == 3.75  (untyped floating-point constant)
const Θ float64 = 3/2      // Θ == 1.0   (type float64, 3/2 is integer division)
const Π float64 = 3/2.     // Π == 1.5   (type float64, 3/2. is float division)
const d = 1 << 3.0         // d == 8     (untyped integer constant)
const e = 1.0 << 3         // e == 8     (untyped integer constant)
const f = int32(1) << 33   // illegal    (constant 8589934592 overflows int32)
const g = float64(2) >> 1  // illegal    (float64(2) is a typed floating-point constant)
const h = "foo" > "bar"    // h == true  (untyped boolean constant)
const j = true             // j == true  (untyped boolean constant)
const k = 'w' + 1          // k == 'x'   (untyped rune constant)
const l = "hi"             // l == "hi"  (untyped string constant)
const m = string(k)        // m == "x"   (type string)
const Σ = 1 - 0.707i       //            (untyped complex constant)
const Δ = Σ + 2.0e-4       //            (untyped complex constant)
const Φ = iota*1i - 1/1i   //            (untyped complex constant)

const ic = complex(0, c)   // ic == 3.75i  (untyped complex constant)
const iΘ = complex(0, Θ)   // iΘ == 1i     (type complex128)

const Huge = 1 << 100         // Huge == 1267650600228229401496703205376  (untyped integer constant)
const Four int8 = Huge >> 98  // Four == 4                                (type int8)
*/
