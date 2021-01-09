package discos

import (
	"cos-backend-com/src/common/flake"
	"cos-backend-com/src/common/validate"
	"cos-backend-com/src/cores/routers"
	"cos-backend-com/src/libs/apierror"
	"cos-backend-com/src/libs/models/discomodels"
	"cos-backend-com/src/libs/sdk/cores"
	"net/http"

	"github.com/wujiu2020/strip/utils/apires"
)

type DiscosInvestorsHandler struct {
	routers.Base
}

func (h *DiscosInvestorsHandler) ListStartupDiscoInvestor(startupId flake.ID) (res interface{}) {
	var input cores.ListDiscoInvestorsInput
	h.Params.BindValuesToStruct(&input)

	if err := validate.Default.Struct(input); err != nil {
		h.Log.Warn(err)
		res = apierror.HandleError(err)
		return
	}

	var output cores.ListDiscoInvestorsResult
	totalEth, total, err := discomodels.DiscoInvestors.ListDiscoInvestor(h.Ctx, startupId, &input, &output.Result)
	if err != nil {
		h.Log.Warn(err)
		res = apierror.HandleError(err)
		return
	}

	if total == 0 {
		output.Result = make([]cores.DiscoInvestorOutput, 0, 0)
	}

	output.TotalEth = totalEth
	output.Total = total
	res = apires.RawWith(&output, http.StatusOK)
	return
}

func (h *DiscosInvestorsHandler) CreateStartupDiscoInvestor(startupId flake.ID) (res interface{}) {
	var input cores.CreateDiscoInvestorInput
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
	if err := discomodels.DiscoInvestors.CreateDiscoInvestor(h.Ctx, startupId, uid, &input); err != nil {
		h.Log.Warn(err)
		res = apierror.HandleError(err)
		return
	}

	res = apires.With(http.StatusOK)
	return
}
