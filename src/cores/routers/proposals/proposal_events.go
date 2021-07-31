package proposals

import (
	"cos-backend-com/src/common/flake"
	"cos-backend-com/src/common/validate"
	"cos-backend-com/src/cores/routers"
	"cos-backend-com/src/libs/apierror"
	"cos-backend-com/src/libs/models/proposalmodels"
	"cos-backend-com/src/libs/models/usermodels"
	"cos-backend-com/src/libs/sdk/account"
	"cos-backend-com/src/libs/sdk/cores"
	"github.com/wujiu2020/strip/utils/apires"
	"net/http"
)

type ProposalEventsHandler struct {
	routers.Base
}

func (h *ProposalEventsHandler) CreateProposal() (res interface{}) {
	var input cores.CreateProposalInput
	if err := h.Params.BindJsonBody(&input); err != nil {
		h.Log.Warn(err)
		res = apierror.HandleError(err)
		return
	}
	if err := validate.Default.Struct(input); err != nil {
		h.Log.Warn(err)
		res = apierror.HandleError(err)
		return
	}

	var user account.UsersModel
	if err := usermodels.Users.GetBypublicKey(h.Ctx, input.WalletAddr, &user); err != nil {
		h.Log.Warn(err)
		res = apierror.HandleError(err)
		return
	}
	input.UserId = user.Id

	var output cores.CreateProposalResult
	if err := proposalmodels.Proposals.CreateProposalWithTerms(h.Ctx, &input, &output); err != nil {
		h.Log.Warn(err)
		res = apierror.HandleError(err)
		return
	}

	res = apires.With(&output, http.StatusOK)
	return
}

func (h *ProposalEventsHandler) UpdateProposalOver() (res interface{}) {
	var output cores.ProposalOverResult
	if err := proposalmodels.Proposals.UpdateProposalOver(h.Ctx, &output.Done); err != nil {
		h.panicIf(err)
	}
	output.Done = true
	res = apires.With(&output, http.StatusOK)
	return
}

func (h *ProposalEventsHandler) UpdateProposalStatus(id flake.ID) (res interface{}) {
	var input cores.UpdateProposalStatusInput
	input.Id = id
	if err := h.Params.BindJsonBody(&input); err != nil {
		h.panicIf(err)
	}

	if err := validate.Default.Struct(input); err != nil {
		h.panicIf(err)
	}
	var output cores.UpdateProposalStatusResult
	if err := proposalmodels.Proposals.UpdateProposalStatus(h.Ctx, &input, &output); err != nil {
		h.panicIf(err)
	}
	res = apires.With(&output, http.StatusOK)
	return
}

func (h *ProposalEventsHandler) VoteProposal(id flake.ID) (res interface{}) {
	var input cores.VoteProposalInput
	input.Id = id
	if err := h.Params.BindJsonBody(&input); err != nil {
		h.panicIf(err)
	}
	if err := validate.Default.Struct(input); err != nil {
		h.panicIf(err)
	}
	var output cores.VoteProposalResult
	if err := proposalmodels.Proposals.VoteProposal(h.Ctx, &input, &output); err != nil {
		h.panicIf(err)
	}
	res = apires.With(&output, http.StatusOK)
	return
}

func (h *ProposalEventsHandler) panicIf(err error) (res interface{}) {
	h.Log.Warn(err)
	return apierror.HandleError(err)
}
