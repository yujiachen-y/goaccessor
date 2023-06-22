// consttest goaccessor supports generating code for constants, however we
// don't recommend to do so.
package consttest

import (
	"testing"

	"github.com/yjc567/goaccessor/test/utils"
)

func TestConstGet(t *testing.T) {
	for _, verifier := range []utils.Verifier{
		utils.NewGetterVerifier(GetPi, Pi),
		utils.NewGetterVerifier(GetZero, zero),
		utils.NewGetterVerifier(GetA, a),
		utils.NewGetterVerifier(GetB, b),
		utils.NewGetterVerifier(GetC, c),
		utils.NewGetterVerifier(GetU, u),
		utils.NewGetterVerifier(GetV, v),
		utils.NewGetterVerifier(GetSize, size),
		utils.NewGetterVerifier(GetEof, eof),
		utils.NewGetterVerifier(GetSunday, Sunday),
		utils.NewGetterVerifier(GetMonday, Monday),
		utils.NewGetterVerifier(GetTuesday, Tuesday),
		utils.NewGetterVerifier(GetWednesday, Wednesday),
		utils.NewGetterVerifier(GetThursday, Thursday),
		utils.NewGetterVerifier(GetFriday, Friday),
		utils.NewGetterVerifier(GetPartyday, Partyday),
		utils.NewGetterVerifier(GetNumberOfDays, numberOfDays),
	} {
		if err := verifier(); err != nil {
			t.Errorf("got err: %s", err.Error())
		}
	}
}
