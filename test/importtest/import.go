package importtest

import (
	"github.com/robfig/cron/v3"
	"github.com/yujiachen-y/goaccessor/test/importtest/p1"
	p3 "github.com/yujiachen-y/goaccessor/test/importtest/p2"
)

//go:generate go run ../../. -t Import -a
type Import struct {
	Option1   p1.Option1
	Option2   *p1.Option2
	Struct    p3.Struct
	StructPtr *p3.Struct
}

//go:generate go run ../../. -t sche -a
var sche cron.Schedule

//go:generate go run ../../. -t anonymous1 -f -a
var anonymous1 struct{ Option1 p1.Option1 }

//go:generate go run ../../. -t anonymous2 -f -a
var anonymous2 = struct{ StructPtr *p3.Struct }{}

type Generic[T any] struct{ T T }

//go:generate go run ../../. -t generic -f -a
var generic = Generic[p1.Option1]{T: "T"}
