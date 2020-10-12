package bounties

import (
	"cos-backend-com/src/common/flake"
	"cos-backend-com/src/common/validate"
	"cos-backend-com/src/cores/routers"
	"cos-backend-com/src/libs/apierror"
	"cos-backend-com/src/libs/models/bountymodels"
	"cos-backend-com/src/libs/sdk/cores"
	"net/http"

	"github.com/wujiu2020/strip/utils/apires"
)

type BountiesHandler struct {
	routers.Base
}

func (h *BountiesHandler) ListBounties() (res interface{}) {
	var params cores.ListBountiesInput
	h.Params.BindValuesToStruct(&params)

	if err := validate.Default.Struct(params); err != nil {
		h.Log.Warn(err)
		res = apierror.HandleError(err)
		return
	}

	var output cores.ListBountiesResult
	total, err := bountymodels.Bounties.ListBounties(h.Ctx, 0, 0, false, &params, &output.Result)
	if err != nil {
		h.Log.Warn(err)
		res = apierror.HandleError(err)
		return
	}
	output.Total = total

	res = apires.With(&output, http.StatusOK)
	return
}

func (h *BountiesHandler) ListStartupBounties(startupId flake.ID) (res interface{}) {
	var params cores.ListBountiesInput
	h.Params.BindValuesToStruct(&params)

	if err := validate.Default.Struct(params); err != nil {
		h.Log.Warn(err)
		res = apierror.HandleError(err)
		return
	}

	var output cores.ListBountiesResult
	total, err := bountymodels.Bounties.ListBounties(h.Ctx, startupId, 0, false, &params, &output.Result)
	if err != nil {
		h.Log.Warn(err)
		res = apierror.HandleError(err)
		return
	}
	output.Total = total

	res = apires.With(&output, http.StatusOK)
	return
}

func (h *BountiesHandler) ListStartupBountiesMe(startupId flake.ID) (res interface{}) {
	var params cores.ListBountiesInput
	h.Params.BindValuesToStruct(&params)

	if err := validate.Default.Struct(params); err != nil {
		h.Log.Warn(err)
		res = apierror.HandleError(err)
		return
	}

	var uid flake.ID
	h.Ctx.Find(&uid, "uid")
	var output cores.ListBountiesResult
	total, err := bountymodels.Bounties.ListBounties(h.Ctx, startupId, 0, true, &params, &output.Result)
	if err != nil {
		h.Log.Warn(err)
		res = apierror.HandleError(err)
		return
	}
	output.Total = total

	res = apires.With(&output, http.StatusOK)
	return
}

func (h *BountiesHandler) ListBountiesMe() (res interface{}) {
	var params cores.ListBountiesInput
	h.Params.BindValuesToStruct(&params)

	if err := validate.Default.Struct(params); err != nil {
		h.Log.Warn(err)
		res = apierror.HandleError(err)
		return
	}

	var uid flake.ID
	h.Ctx.Find(&uid, "uid")
	var output cores.ListBountiesResult
	total, err := bountymodels.Bounties.ListBounties(h.Ctx, 0, uid, false, &params, &output.Result)
	if err != nil {
		h.Log.Warn(err)
		res = apierror.HandleError(err)
		return
	}
	output.Total = total

	res = apires.With(&output, http.StatusOK)
	return
}

func (h *BountiesHandler) ListUserBounties(uid flake.ID) (res interface{}) {
	var params cores.ListBountiesInput
	h.Params.BindValuesToStruct(&params)

	if err := validate.Default.Struct(params); err != nil {
		h.Log.Warn(err)
		res = apierror.HandleError(err)
		return
	}

	var output cores.ListBountiesResult
	total, err := bountymodels.Bounties.ListBounties(h.Ctx, 0, uid, false, &params, &output.Result)
	if err != nil {
		h.Log.Warn(err)
		res = apierror.HandleError(err)
		return
	}
	output.Total = total

	res = apires.With(&output, http.StatusOK)
	return
}

func (h *BountiesHandler) GetBounty(id flake.ID) (res interface{}) {
	var output cores.BountyOutput
	if err := bountymodels.Bounties.GetBounty(h.Ctx, id, false, &output); err != nil {
		h.Log.Warn(err)
		res = apierror.HandleError(err)
		return
	}

	res = apires.With(&output, http.StatusOK)
	return
}

func (h *BountiesHandler) GetBountyMe(id flake.ID) (res interface{}) {
	var uid flake.ID
	h.Ctx.Find(&uid, "uid")
	var output cores.BountyOutput
	if err := bountymodels.Bounties.GetBounty(h.Ctx, id, true, &output); err != nil {
		h.Log.Warn(err)
		res = apierror.HandleError(err)
		return
	}

	res = apires.With(&output, http.StatusOK)
	return
}

func (h *BountiesHandler) Create(startupId flake.ID) (res interface{}) {
	var input cores.CreateBountyInput
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
	var uid flake.ID
	h.Ctx.Find(&uid, "uid")
	if err := bountymodels.Bounties.CreateBounty(h.Ctx, startupId, uid, &input); err != nil {
		h.Log.Warn(err)
		res = apierror.HandleError(err)
		return
	}

	res = apires.With(&cores.PrepareIdOutput{
		input.Id,
	}, http.StatusOK)
	return
}

func (h *BountiesHandler) Closed(bountyId flake.ID) (res interface{}) {
	var uid flake.ID
	h.Ctx.Find(&uid, "uid")
	if err := bountymodels.Bounties.ClosedBounty(h.Ctx, bountyId, uid); err != nil {
		h.Log.Warn(err)
		res = apierror.HandleError(err)
		return
	}

	res = apires.Ret(http.StatusOK)
	return
}

func (h *BountiesHandler) StartWork(bountyId flake.ID) (res interface{}) {
	var input cores.CreateUndertakeBountyInput
	if err := h.Params.BindJsonBody(&input); err != nil {
		h.Log.Warn(err)
		res = apierror.HandleError(err)
		return
	}
	h.Ctx.Find(&input.UserId, "uid")
	input.BountyId = bountyId
	if err := validate.Default.Struct(input); err != nil {
		h.Log.Warn(err)
		res = apierror.HandleError(err)
		return
	}

	var output cores.UndertakeBountyResult
	if err := bountymodels.Bounties.CreateUndertakeBounty(h.Ctx, &input, &output); err != nil {
		h.Log.Warn(err)
		res = apierror.HandleError(err)
		return
	}

	res = apires.With(&output, http.StatusOK)
	return
}

func (h *BountiesHandler) Submitted(bountyId flake.ID) (res interface{}) {
	var input cores.UpdateUndertakeBountyInput
	h.Ctx.Find(&input.UserId, "uid")
	input.BountyId = bountyId
	if err := validate.Default.Struct(input); err != nil {
		h.Log.Warn(err)
		res = apierror.HandleError(err)
		return
	}

	var output cores.UndertakeBountyResult
	if err := bountymodels.Bounties.SubmittedUndertakeBounty(h.Ctx, &input, &output); err != nil {
		h.Log.Warn(err)
		res = apierror.HandleError(err)
		return
	}

	res = apires.With(&output, http.StatusOK)
	return
}

func (h *BountiesHandler) Quited(bountyId flake.ID) (res interface{}) {
	var input cores.UpdateUndertakeBountyInput
	h.Ctx.Find(&input.UserId, "uid")
	input.BountyId = bountyId
	if err := validate.Default.Struct(input); err != nil {
		h.Log.Warn(err)
		res = apierror.HandleError(err)
		return
	}

	var output cores.UndertakeBountyResult
	if err := bountymodels.Bounties.QuitedUndertakeBounty(h.Ctx, &input, &output); err != nil {
		h.Log.Warn(err)
		res = apierror.HandleError(err)
		return
	}

	res = apires.With(&output, http.StatusOK)
	return
}

func (h *BountiesHandler) Paid(bountyId flake.ID) (res interface{}) {
	var input cores.UpdateUndertakeBountyInput
	if err := h.Params.BindJsonBody(&input); err != nil {
		h.Log.Warn(err)
		res = apierror.HandleError(err)
		return
	}

	input.BountyId = bountyId
	if err := validate.Default.Struct(input); err != nil {
		h.Log.Warn(err)
		res = apierror.HandleError(err)
		return
	}

	var output cores.UndertakeBountyResult
	if err := bountymodels.Bounties.PaidUndertakeBounty(h.Ctx, &input, &output); err != nil {
		h.Log.Warn(err)
		res = apierror.HandleError(err)
		return
	}

	res = apires.With(&output, http.StatusOK)
	return
}

func (h *BountiesHandler) Rejected(bountyId flake.ID) (res interface{}) {
	var input cores.UpdateUndertakeBountyInput
	if err := h.Params.BindJsonBody(&input); err != nil {
		h.Log.Warn(err)
		res = apierror.HandleError(err)
		return
	}

	input.BountyId = bountyId
	if err := validate.Default.Struct(input); err != nil {
		h.Log.Warn(err)
		res = apierror.HandleError(err)
		return
	}

	var output cores.UndertakeBountyResult
	if err := bountymodels.Bounties.RejectedUndertakeBounty(h.Ctx, &input, &output); err != nil {
		h.Log.Warn(err)
		res = apierror.HandleError(err)
		return
	}

	res = apires.With(&output, http.StatusOK)
	return
}
