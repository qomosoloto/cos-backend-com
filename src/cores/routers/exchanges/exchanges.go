package exchanges

import (
	"cos-backend-com/src/common/flake"
	"cos-backend-com/src/common/validate"
	"cos-backend-com/src/cores/routers"
	"cos-backend-com/src/libs/apierror"
	"cos-backend-com/src/libs/models/exchangemodels"
	"cos-backend-com/src/libs/models/startupmodels"
	"cos-backend-com/src/libs/sdk/cores"
	"github.com/wujiu2020/strip/utils/apires"
	"net/http"
)

type ExchangesHandler struct {
	routers.Base
}

func (h *ExchangesHandler) CreateExchange(startupId flake.ID) (res interface{}) {
	var startup cores.StartUpResult
	if err := startupmodels.Startups.Get(h.Ctx, startupId, &startup); err != nil {
		h.Log.Warn(err)
		res = apierror.HandleError(err)
		return
	}

	var input cores.CreateExchangeInput
	if err := h.Params.BindJsonBody(&input); err != nil {
		h.Log.Warn(err)
		res = apierror.HandleError(err)
		return
	}
	input.StartupId = startup.Id
	input.TokenName1 = startup.Setting.TokenName
	input.TokenSymbol1 = startup.Setting.TokenSymbol
	input.TokenAddress1 = *startup.Setting.TokenAddr
	input.TokenName2 = "ETH"
	input.TokenSymbol2 = "ETH"
	input.Status = cores.ExchangeStatusPending

	if err := validate.Default.Struct(input); err != nil {
		h.Log.Warn(err)
		res = apierror.HandleError(err)
		return
	}

	var output cores.CreateExchangeResult
	if err := exchangemodels.Exchanges.CreateExchange(h.Ctx, &input, &output); err != nil {
		h.Log.Warn(err)
		res = apierror.HandleError(err)
		return
	}

	res = apires.With(&output, http.StatusOK)
	return
}

func (h *ExchangesHandler) GetExchange(id flake.ID) (res interface{}) {
	var input cores.GetExchangeInput
	input.Id = id
	var output cores.ExchangeResult
	if err := exchangemodels.Exchanges.GetExchange(h.Ctx, &input, &output); err != nil {
		h.Log.Warn(err)
		res = apierror.HandleError(err)
		return
	}

	res = apires.With(&output, http.StatusOK)
	return
}

func (h *ExchangesHandler) GetExchangeByStartup(id flake.ID) (res interface{}) {
	var input cores.GetExchangeInput
	input.StartupId = id
	var output cores.ExchangeResult
	if err := exchangemodels.Exchanges.GetExchange(h.Ctx, &input, &output); err != nil {
		h.Log.Warn(err)
		res = apierror.HandleError(err)
		return
	}

	res = apires.With(&output, http.StatusOK)
	return
}

func (h *ExchangesHandler) CreateExchangeTx(exchangeId flake.ID) (res interface{}) {
	var input cores.CreateExchangeTxInput
	if err := h.Params.BindJsonBody(&input); err != nil {
		h.Log.Warn(err)
		res = apierror.HandleError(err)
		return
	}
	input.Status = cores.ExchangeTxStatusPending

	if err := validate.Default.Struct(input); err != nil {
		h.Log.Warn(err)
		res = apierror.HandleError(err)
		return
	}

	var output cores.CreateExchangeTxResult
	if err := exchangemodels.Exchanges.CreateExchangeTx(h.Ctx, &input, &output); err != nil {
		h.Log.Warn(err)
		res = apierror.HandleError(err)
		return
	}

	res = apires.With(&output, http.StatusOK)
	return
}

func (h *ExchangesHandler) GetExchangeTx(id flake.ID) (res interface{}) {
	var input cores.GetExchangeTxInput
	input.Id = id
	var output cores.ExchangeTxResult
	if err := exchangemodels.Exchanges.GetExchangeTx(h.Ctx, &input, &output); err != nil {
		h.Log.Warn(err)
		res = apierror.HandleError(err)
		return
	}

	res = apires.With(&output, http.StatusOK)
	return
}
