package importtest

import (
	p3 "github.com/yujiachen-y/goaccessor/test/importtest/p1"
	"github.com/yujiachen-y/goaccessor/test/importtest/p2"
)

func (i *Import) SetOption1(v p3.Option1) {
	i.Option1 = v
}

func (i *Import) GetStruct() p2.Struct {
	return i.Struct
}
