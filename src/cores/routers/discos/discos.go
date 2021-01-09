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

type DiscosHandler struct {
	routers.Base
}

func (h *DiscosHandler) CreateStartupDisco(startupId flake.ID) (res interface{}) {
	var input cores.CreateDiscosInput
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

	if err := discomodels.Discos.CreateDisco(h.Ctx, startupId, &input); err != nil {
		h.Log.Warn(err)
		res = apierror.HandleError(err)
		return
	}

	res = apires.With(http.StatusOK)
	return
}

func (h *DiscosHandler) GetStartupDisco(startupId flake.ID) (res interface{}) {
	var output cores.StartupDiscosResult
	if err := discomodels.Discos.GetDisco(h.Ctx, startupId, &output); err != nil {
		h.Log.Warn(err)
		res = apierror.HandleError(err)
		return
	}

	res = apires.With(&output, http.StatusOK)
	return
}

func (h *DiscosHandler) ListDisco() (res interface{}) {
	var input cores.ListDiscosInput
	h.Params.BindValuesToStruct(&input)

	if err := validate.Default.Struct(input); err != nil {
		h.Log.Warn(err)
		res = apierror.HandleError(err)
		return
	}

	var output cores.ListDiscosResult
	total, err := discomodels.Discos.ListDisco(h.Ctx, &input, &output.Result)
	if err != nil {
		h.Log.Warn(err)
		res = apierror.HandleError(err)
		return
	}

	if total == 0 {
		output.Result = make([]cores.DiscoOutput, 0, 0)
	}

	output.Total = total

	res = apires.With(&output, http.StatusOK)
	return
}
