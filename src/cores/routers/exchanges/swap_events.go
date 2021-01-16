package exchanges

import (
	"cos-backend-com/src/common/validate"
	"cos-backend-com/src/cores/routers"
	"cos-backend-com/src/libs/apierror"
	"cos-backend-com/src/libs/models/exchangemodels"
	"cos-backend-com/src/libs/sdk/cores"
	"github.com/wujiu2020/strip/utils/apires"
	"net/http"
)

type SwapEventsHandler struct {
	routers.Base
}

func (h *ExchangesHandler) CreatePair() (res interface{}) {
	var pairinput cores.CreateSwapPairInput
	if err := h.Params.BindJsonBody(&pairinput); err != nil {
		h.Log.Warn(err)
		res = apierror.HandleError(err)
		return
	}
	if err := validate.Default.Struct(pairinput); err != nil {
		h.Log.Warn(err)
		res = apierror.HandleError(err)
		return
	}

	var input cores.CreateExchangeInput
	input.TxId = pairinput.TxId
	input.StartupId = pairinput.StartupId
	input.PairAddress = pairinput.PairAddress
	input.TokenName1 = pairinput.Token0.Name
	input.TokenSymbol1 = pairinput.Token0.Symbol
	input.TokenAddress1 = pairinput.Token0.Address
	input.TokenDivider1 = power(10, pairinput.Token0.Decimals)
	input.TokenName2 = pairinput.Token1.Name
	input.TokenSymbol2 = pairinput.Token1.Symbol
	input.TokenAddress2 = pairinput.Token1.Address
	input.TokenDivider2 = power(10, pairinput.Token1.Decimals)
	input.PairName = input.TokenName1 + "-" + input.TokenName2
	input.Status = cores.ExchangeStatusCompleted
	if err := validate.Default.Struct(input); err != nil {
		h.Log.Warn(err)
		res = apierror.HandleError(err)
		return
	}

	var output cores.CreateExchangeResult
	var exchangeinput cores.GetExchangeInput
	exchangeinput.StartupId = pairinput.StartupId
	var exchangeresult cores.ExchangeResult
	if err := exchangemodels.Exchanges.GetExchange(h.Ctx, &exchangeinput, &exchangeresult); err == nil {
		if err := exchangemodels.Exchanges.UpdateExchange(h.Ctx, &input, &output); err != nil {
			h.Log.Warn(err)
			res = apierror.HandleError(err)
			return
		}
	} else {
		if err := exchangemodels.Exchanges.CreateExchange(h.Ctx, &input, &output); err != nil {
			h.Log.Warn(err)
			res = apierror.HandleError(err)
			return
		}
	}

	res = apires.With(&output, http.StatusOK)
	return

}

func power(x int, n int) int {
	if n == 0 {
		return 1
	} else {
		return x * power(x, n-1)
	}
}
