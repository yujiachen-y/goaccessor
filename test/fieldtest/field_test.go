package fieldtest

import (
	"testing"

	"github.com/yjc567/goaccessor/test/utils"
)

func TestNormal1(t *testing.T) {
	normal1 = &Normal{}
	testInt2 := int64(2)
	testMap4 := map[string]complex128{"test": complex(4, 0)}
	for _, verifier := range []utils.Verifier{
		utils.NewSetterVerifier(&normal1.a, SetNormal1A, 2),
		utils.NewGetterVerifier(GetNormal1A, 2),
		utils.NewSetterVerifier(&normal1.B, SetNormal1B, 3),
		utils.NewGetterVerifier(GetNormal1B, 3),
		utils.NewSetterVerifier(&normal1.c, SetNormal1C, &testInt2),
		utils.NewGetterVerifier(GetNormal1C, &testInt2),
		utils.NewMapSetterVerifier(&normal1.d, SetNormal1D, testMap4),
		utils.NewMapGetterVerifier(GetNormal1D, testMap4),
	} {
		if err := verifier(); err != nil {
			t.Errorf("got error: %s", err.Error())
		}
	}
}

func TestNormal2(t *testing.T) {
	testInt5 := int64(5)
	testMap7 := map[string]complex128{"test": complex(7, 0)}
	for _, verifier := range []utils.Verifier{
		utils.NewSetterVerifier(&normal2.a, SetNormal2A, 5),
		utils.NewGetterVerifier(GetNormal2A, 5),
		utils.NewSetterVerifier(&normal2.B, SetNormal2B, 6),
		utils.NewGetterVerifier(GetNormal2B, 6),
		utils.NewSetterVerifier(&normal2.c, SetNormal2C, &testInt5),
		utils.NewGetterVerifier(GetNormal2C, &testInt5),
		utils.NewMapSetterVerifier(&normal2.d, SetNormal2D, testMap7),
		utils.NewMapGetterVerifier(GetNormal2D, testMap7),
	} {
		if err := verifier(); err != nil {
			t.Errorf("got error: %s", err.Error())
		}
	}
}

func TestGeneric1(t *testing.T) {
	testInt := 2
	testFloat64 := 3.0
	testString := "test4"
	testMap := map[int]string{5: "test6"}
	testSlices := []float64{7.0, 8.0}
	for _, verifier := range []utils.Verifier{
		utils.NewSetterVerifier(&generic1.t, SetGeneric1T, testInt),
		utils.NewGetterVerifier(GetGeneric1T, testInt),
		utils.NewSetterVerifier(&generic1.u, SetGeneric1U, &testFloat64),
		utils.NewGetterVerifier(GetGeneric1U, &testFloat64),
		utils.NewSetterVerifier(&generic1.v, SetGeneric1V, testString),
		utils.NewGetterVerifier(GetGeneric1V, testString),
		utils.NewMapSetterVerifier(&generic1.m, SetGeneric1M, testMap),
		utils.NewMapGetterVerifier(GetGeneric1M, testMap),
		utils.NewSliceSetterVerifier(&generic1.s, SetGeneric1S, testSlices),
		utils.NewSliceGetterVerifier(GetGeneric1S, testSlices),
	} {
		if err := verifier(); err != nil {
			t.Errorf("got error: %s", err.Error())
		}
	}
}

func TestGeneric2(t *testing.T) {
	generic2 = &Generic[float64, string, int]{}
	testFloat64 := 9.0
	testString := "test10"
	testInt := 11
	testMap := map[float64]int{12.0: 13}
	testSlices := []string{"test14", "test15"}
	for _, verifier := range []utils.Verifier{
		utils.NewSetterVerifier(&generic2.t, SetGeneric2T, testFloat64),
		utils.NewGetterVerifier(GetGeneric2T, testFloat64),
		utils.NewSetterVerifier(&generic2.u, SetGeneric2U, &testString),
		utils.NewGetterVerifier(GetGeneric2U, &testString),
		utils.NewSetterVerifier(&generic2.v, SetGeneric2V, testInt),
		utils.NewGetterVerifier(GetGeneric2V, testInt),
		utils.NewMapSetterVerifier(&generic2.m, SetGeneric2M, testMap),
		utils.NewMapGetterVerifier(GetGeneric2M, testMap),
		utils.NewSliceSetterVerifier(&generic2.s, SetGeneric2S, testSlices),
		utils.NewSliceGetterVerifier(GetGeneric2S, testSlices),
	} {
		if err := verifier(); err != nil {
			t.Errorf("got error: %s", err.Error())
		}
	}
}

func TestGeneric3(t *testing.T) {
	testInt := 2
	testFloat := 3.0
	testMap := map[string]float64{"test": 4.0}
	testSlice := []int{5, 6}
	testString := "test7"
	for _, verifier := range []utils.Verifier{
		utils.NewSetterVerifier(&generic3.t, SetGeneric3T, testString),
		utils.NewGetterVerifier(GetGeneric3T, testString),
		utils.NewSetterVerifier(&generic3.u, SetGeneric3U, &testInt),
		utils.NewGetterVerifier(GetGeneric3U, &testInt),
		utils.NewSetterVerifier(&generic3.v, SetGeneric3V, testFloat),
		utils.NewGetterVerifier(GetGeneric3V, testFloat),
		utils.NewMapSetterVerifier(&generic3.m, SetGeneric3M, testMap),
		utils.NewMapGetterVerifier(GetGeneric3M, testMap),
		utils.NewSliceSetterVerifier(&generic3.s, SetGeneric3S, testSlice),
		utils.NewSliceGetterVerifier(GetGeneric3S, testSlice),
	} {
		if err := verifier(); err != nil {
			t.Errorf("got error: %s", err.Error())
		}
	}
}

func TestGeneric4(t *testing.T) {
	testInt := 8
	testFloat := 9.0
	testMap := map[string]float64{"test": 10.0}
	testSlice := []int{11, 12}
	testString := "test13"
	for _, verifier := range []utils.Verifier{
		utils.NewSetterVerifier(&generic4.t, SetGeneric4T, testString),
		utils.NewGetterVerifier(GetGeneric4T, testString),
		utils.NewSetterVerifier(&generic4.u, SetGeneric4U, &testInt),
		utils.NewGetterVerifier(GetGeneric4U, &testInt),
		utils.NewSetterVerifier(&generic4.v, SetGeneric4V, testFloat),
		utils.NewGetterVerifier(GetGeneric4V, testFloat),
		utils.NewMapSetterVerifier(&generic4.m, SetGeneric4M, testMap),
		utils.NewMapGetterVerifier(GetGeneric4M, testMap),
		utils.NewSliceSetterVerifier(&generic4.s, SetGeneric4S, testSlice),
		utils.NewSliceGetterVerifier(GetGeneric4S, testSlice),
	} {
		if err := verifier(); err != nil {
			t.Errorf("got error: %s", err.Error())
		}
	}
}

func TestAnonymous1(t *testing.T) {
	testQuite := "2"
	testDifficult := "3"

	for _, verifier := range []utils.Verifier{
		utils.NewSetterVerifier(&anonymous1.quite, SetAnonymous1Quite, testQuite),
		utils.NewGetterVerifier(GetAnonymous1Quite, testQuite),
		utils.NewSetterVerifier(&anonymous1.difficult, SetAnonymous1Difficult, testDifficult),
		utils.NewGetterVerifier(GetAnonymous1Difficult, testDifficult),
	} {
		if err := verifier(); err != nil {
			t.Errorf("got error: %s", err.Error())
		}
	}
}

func TestAnonymous2(t *testing.T) {
	testQuite := "4"
	testDifficult := "5"

	for _, verifier := range []utils.Verifier{
		utils.NewSetterVerifier(&anonymous2.quite, SetAnonymous2Quite, testQuite),
		utils.NewGetterVerifier(GetAnonymous2Quite, testQuite),
		utils.NewSetterVerifier(&anonymous2.difficult, SetAnonymous2Difficult, testDifficult),
		utils.NewGetterVerifier(GetAnonymous2Difficult, testDifficult),
	} {
		if err := verifier(); err != nil {
			t.Errorf("got error: %s", err.Error())
		}
	}
}

func TestAnonymous3(t *testing.T) {
	quite3Value := "Quite Value 2"
	difficult3Value := "Difficult Value 3"
	for _, verifier := range []utils.Verifier{
		utils.NewSetterVerifier(&anonymous3.quite, SetAnonymous3Quite, quite3Value),
		utils.NewGetterVerifier(GetAnonymous3Quite, quite3Value),
		utils.NewSetterVerifier(&anonymous3.difficult, SetAnonymous3Difficult, difficult3Value),
		utils.NewGetterVerifier(GetAnonymous3Difficult, difficult3Value),
	} {
		if err := verifier(); err != nil {
			t.Errorf("got error: %s", err.Error())
		}
	}
}

func TestAnonymous4(t *testing.T) {
	quite4Value := "Quite Value 4"
	difficult4Value := "Difficult Value 5"
	for _, verifier := range []utils.Verifier{
		utils.NewSetterVerifier(&anonymous4.quite, SetAnonymous4Quite, quite4Value),
		utils.NewGetterVerifier(GetAnonymous4Quite, quite4Value),
		utils.NewSetterVerifier(&anonymous4.difficult, SetAnonymous4Difficult, difficult4Value),
		utils.NewGetterVerifier(GetAnonymous4Difficult, difficult4Value),
	} {
		if err := verifier(); err != nil {
			t.Errorf("got error: %s", err.Error())
		}
	}
}

func TestAnonymous5(t *testing.T) {
	// Initialize the anonymous5 struct
	anonymous5 = &struct {
		quite, difficult string
	}{}

	// Set test values
	testQuite := "2"
	testDifficult := "3"

	for _, verifier := range []utils.Verifier{
		utils.NewSetterVerifier(&anonymous5.quite, SetAnonymous5Quite, testQuite),
		utils.NewGetterVerifier(GetAnonymous5Quite, testQuite),
		utils.NewSetterVerifier(&anonymous5.difficult, SetAnonymous5Difficult, testDifficult),
		utils.NewGetterVerifier(GetAnonymous5Difficult, testDifficult),
	} {
		if err := verifier(); err != nil {
			t.Errorf("got error: %s", err.Error())
		}
	}
}
