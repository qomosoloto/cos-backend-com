package cores

import (
	"cos-backend-com/src/common/validate"
	"testing"
)

func TestUpdateProposalStatusInput_ValidateProposalStatus(t *testing.T) {
	input := UpdateProposalStatusInput{Id: 1, Status: ProposalStatus(0)}
	assertion := func(t testing.TB, input UpdateProposalStatusInput, erHandler func(er error)) {
		t.Helper()
		err := validate.Default.Struct(&input)
		erHandler((err))
	}
	t.Run("Validation should fail if Status not in (4, 5, 6)", func(t *testing.T) {
		assertion(t, input, func(er error) {
			if er == nil {
				t.Fatal("Validation should not pass!")
			}
		})
	})

	t.Run("Validation should pass if Status in (4, 5, 6)", func(t *testing.T) {
		input.Status = 5
		assertion(t, input, func(er error) {
			if er != nil {
				t.Fatal("Validation should pass!")
			}
		})
	})

}
