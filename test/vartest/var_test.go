package vartest

import (
	"testing"

	"github.com/yjc567/goaccessor/test/utils"
)

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
