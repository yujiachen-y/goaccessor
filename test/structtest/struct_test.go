package structtest

import (
	"testing"

	"github.com/yjc567/goaccessor/test/utils"
)

func TestStruct(t *testing.T) {
	n := Normal{}
	testInt64 := int64(5)
	ch := make(chan float64)
	testMap := map[string]complex128{"test": complex(6, 0)}
	testSlices := []*S{{}, {}}
	for _, verifier := range []utils.Verifier{
		utils.NewSetterVerifier(&n.a, n.SetA, 2),
		utils.NewGetterVerifier(n.GetA, 2),
		utils.NewSetterVerifier(&n.B, n.SetB, 3),
		utils.NewGetterVerifier(n.GetB, 3),
		utils.NewSetterVerifier(&n.c, n.SetC, &testInt64),
		utils.NewGetterVerifier(n.GetC, &testInt64),
		utils.NewSetterVerifier(&n.D, n.SetD, ch),
		utils.NewGetterVerifier(n.GetD, ch),
		utils.NewMapSetterVerifier(&n.e, n.SetE, testMap),
		utils.NewMapGetterVerifier(n.GetE, testMap),
		utils.NewSliceSetterVerifier(&n.f, n.SetF, testSlices),
		utils.NewSliceGetterVerifier(n.GetF, testSlices),
	} {
		if err := verifier(); err != nil {
			t.Errorf("got error: %s", err.Error())
		}
	}
}

func TestGeneric(t *testing.T) {
	g := Generic[string, int]{}
	testInt64 := int64(5)
	ch := make(chan float64)
	testMap := map[string]complex128{"test": complex(6, 0)}
	testSlices := []*S{{}, {}}
	testS := S{}
	testT := "test string"
	testU := 9
	for _, verifier := range []utils.Verifier{
		utils.NewSetterVerifier(&g.B, g.SetB, 3),
		utils.NewGetterVerifier(g.GetB, 3),
		utils.NewSetterVerifier(&g.c, g.SetC, &testInt64),
		utils.NewGetterVerifier(g.GetC, &testInt64),
		utils.NewSetterVerifier(&g.D, g.SetD, ch),
		utils.NewGetterVerifier(g.GetD, ch),
		utils.NewMapSetterVerifier(&g.e, g.SetE, testMap),
		utils.NewMapGetterVerifier(g.GetE, testMap),
		utils.NewSliceSetterVerifier(&g.f, g.SetF, testSlices),
		utils.NewSliceGetterVerifier(g.GetF, testSlices),
		utils.NewSetterVerifier(&g.s, g.SetS, testS),
		utils.NewGetterVerifier(g.GetS, testS),
		utils.NewSetterVerifier(&g.t, g.SetT, testT),
		utils.NewGetterVerifier(g.GetT, testT),
		utils.NewSetterVerifier(&g.u, g.SetU, testU),
		utils.NewGetterVerifier(g.GetU, testU),
	} {
		if err := verifier(); err != nil {
			t.Errorf("got error: %s", err.Error())
		}
	}
}
