package proposals

import (
	"cos-backend-com/src/common/flake"
	"cos-backend-com/src/cores/routers"
	"cos-backend-com/src/libs/apierror"
	"cos-backend-com/src/libs/models/proposalmodels"
	"cos-backend-com/src/libs/sdk/cores"
	"github.com/wujiu2020/strip/utils/apires"
	"net/http"
)

type ProposalsHandler struct {
	routers.Base
}

func (h *ProposalsHandler) GetProposal(id flake.ID) (res interface{}) {
	var input cores.GetProposalInput
	input.Id = id
	var output cores.ProposalResult
	if err := proposalmodels.Proposals.GetProposal(h.Ctx, &input, &output); err != nil {
		h.Log.Warn(err)
		res = apierror.HandleError(err)
		return
	}

	res = apires.With(&output, http.StatusOK)
	return
}
