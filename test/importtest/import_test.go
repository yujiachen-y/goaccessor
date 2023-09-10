package importtest

import (
	"testing"
	"time"

	"github.com/yujiachen-y/goaccessor/test/importtest/p1"
	"github.com/yujiachen-y/goaccessor/test/importtest/p2"
	"github.com/yujiachen-y/goaccessor/test/utils"
)

func TestAnonymous1(t *testing.T) {
	option1 := p1.Option1("1")
	for _, verifier := range []utils.Verifier{
		utils.NewGetterVerifier(GetOption1, ""),
		utils.NewSetterVerifier(&anonymous1.Option1, SetOption1, option1),
	} {
		if err := verifier(); err != nil {
			t.Errorf("got error: %s", err.Error())
		}
	}
}

func TestAnonymous2(t *testing.T) {
	s := &p2.Struct{
		Option1: "2",
		Option2: 3,
	}
	for _, verifier := range []utils.Verifier{
		utils.NewGetterVerifier(GetStructPtr, nil),
		utils.NewSetterVerifier(&anonymous2.StructPtr, SetStructPtr, s),
	} {
		if err := verifier(); err != nil {
			t.Errorf("got error: %s", err.Error())
		}
	}
}

func TestImport(t *testing.T) {
	option1 := p1.Option1("4")
	ipt := &Import{
		Option1: option1,
	}
	option2 := new(p1.Option2)
	stc := p2.Struct{
		Option1: "5",
		Option2: 6,
	}
	stcPTR := &p2.Struct{
		Option1: "7",
		Option2: 8,
	}
	for _, verifier := range []utils.Verifier{
		utils.NewGetterVerifier(ipt.GetOption1, option1),
		utils.NewGetterVerifier(ipt.GetOption2, nil),
		utils.NewSetterVerifier(&ipt.Option2, ipt.SetOption2, option2),
		utils.NewSetterVerifier(&ipt.Struct, ipt.SetStruct, stc),
		utils.NewGetterVerifier(ipt.GetStructPtr, nil),
		utils.NewSetterVerifier(&ipt.StructPtr, ipt.SetStructPtr, stcPTR),
	} {
		if err := verifier(); err != nil {
			t.Errorf("got error: %s", err.Error())
		}
	}
}

type testSchedule struct {
	t time.Time
}

func (t *testSchedule) Next(_ time.Time) time.Time {
	return t.t
}

func TestSche(t *testing.T) {
	now := time.Now()
	testSche := &testSchedule{t: now}
	SetSche(testSche)
	if next := GetSche().Next(now); next != now {
		t.Errorf("expected %v got %v", now, next)
	}
}

func TestGeneric(t *testing.T) {
	option1 := p1.Option1("9")
	for _, verifier := range []utils.Verifier{
		utils.NewGetterVerifier(GetT, "T"),
		utils.NewSetterVerifier(&generic.T, SetT, option1),
	} {
		if err := verifier(); err != nil {
			t.Errorf("got error: %s", err.Error())
		}
	}
}
