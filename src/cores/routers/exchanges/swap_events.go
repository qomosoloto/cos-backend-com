package exchanges

import (
	"cos-backend-com/src/common/validate"
	"cos-backend-com/src/cores/routers"
	"cos-backend-com/src/libs/apierror"
	"cos-backend-com/src/libs/models/exchangemodels"
	"cos-backend-com/src/libs/models/startupmodels"
	"cos-backend-com/src/libs/sdk/cores"
	"github.com/shopspring/decimal"
	"github.com/wujiu2020/strip/utils/apires"
	"net/http"
	"strconv"
)

type SwapEventsHandler struct {
	routers.Base
}

func (h *SwapEventsHandler) CreatePair() (res interface{}) {
	var pairinput cores.SwapPairCreatedInput
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

	/* 2021-06-14 paul
	var exchangeinput cores.GetExchangeInput
	exchangeinput.TxId = pairinput.TxId
	var exchangeresult cores.ExchangeResult
	if err := exchangemodels.Exchanges.GetExchange(h.Ctx, &exchangeinput, &exchangeresult); err != nil {
		h.Log.Warn(err)
		res = apierror.HandleError(err)
		return
	}*/
	/* 2021-06-14 paul */
	startupId, err := startupmodels.Startups.GetId(h.Ctx, pairinput.Token0.Address)
	if err != nil {
		h.Log.Warn(err)
		res = apierror.HandleError(err)
		return
	}

	var input cores.CreateExchangeInput
	input.TxId = pairinput.TxId
	input.StartupId = startupId
	input.PairAddress = pairinput.PairAddress
	input.TokenName1 = pairinput.Token0.Name
	input.TokenSymbol1 = pairinput.Token0.Symbol
	input.TokenAddress1 = pairinput.Token0.Address
	input.TokenDivider1 = power(10, pairinput.Token0.Decimals)
	input.TokenName2 = pairinput.Token1.Name
	input.TokenSymbol2 = pairinput.Token1.Symbol
	input.TokenAddress2 = pairinput.Token1.Address
	input.TokenDivider2 = power(10, pairinput.Token1.Decimals)
	input.PairName = input.TokenSymbol1 + "-" + input.TokenSymbol2
	input.Status = cores.ExchangeStatusCompleted
	if err := validate.Default.Struct(input); err != nil {
		h.Log.Warn(err)
		res = apierror.HandleError(err)
		return
	}

	/* 2021-06-14 paul
	var output cores.CreateExchangeResult
	if exchangeresult.Status != cores.ExchangeStatusCompleted {
		if err := exchangemodels.Exchanges.UpdateExchange(h.Ctx, &input, &output); err != nil {
			h.Log.Warn(err)
			res = apierror.HandleError(err)
			return
		}
	}*/
	/* 2021-06-14 paul */
	var output cores.CreateExchangeResult
	if err := exchangemodels.Exchanges.CreateExchange(h.Ctx, &input, &output); err != nil {
		h.Log.Warn(err)
		res = apierror.HandleError(err)
		return
	}

	res = apires.With(&output, http.StatusOK)
	return
}

func (h *SwapEventsHandler) Sync() (res interface{}) {
	var syncinput cores.SwapSyncInput
	if err := h.Params.BindJsonBody(&syncinput); err != nil {
		h.Log.Warn(err)
		res = apierror.HandleError(err)
		return
	}
	if err := validate.Default.Struct(syncinput); err != nil {
		h.Log.Warn(err)
		res = apierror.HandleError(err)
		return
	}

	var exchangeinput cores.GetExchangeInput
	exchangeinput.PairAddress = syncinput.PairAddress
	var balanceresult cores.ExchangeBalanceResult
	if err := exchangemodels.Exchanges.GetBalance(h.Ctx, &exchangeinput, &balanceresult); err != nil {
		h.Log.Warn(err)
		res = apierror.HandleError(err)
		return
	}

	var input cores.ExchangeBalanceInput
	input.StartupId = balanceresult.StartupId
	input.Reserve0 = syncinput.Reserve0
	input.Reserve1 = syncinput.Reserve1

	amount1, _ := decimal.NewFromString(input.Reserve0)
	divider1Str := strconv.Itoa(balanceresult.TokenDivider1)
	divider1, _ := decimal.NewFromString(divider1Str)
	input.NewestPooledTokens1, _ = amount1.Div(divider1).Float64()
	amount2, _ := decimal.NewFromString(input.Reserve1)
	divider2Str := strconv.Itoa(balanceresult.TokenDivider2)
	divider2, _ := decimal.NewFromString(divider2Str)
	input.NewestPooledTokens2, _ = amount2.Div(divider2).Float64()

	occuredday := syncinput.OccuredAt[0:10]
	input.NewestDay = occuredday
	if balanceresult.NewestDay != occuredday && balanceresult.NewestDay != "" {
		input.LastDay = balanceresult.NewestDay
		input.LastPooledTokens1 = balanceresult.NewestPooledTokens1
		input.LastPooledTokens2 = balanceresult.NewestPooledTokens2
	}

	if err := validate.Default.Struct(input); err != nil {
		h.Log.Warn(err)
		res = apierror.HandleError(err)
		return
	}

	var output cores.CreateExchangeResult
	if err := exchangemodels.Exchanges.UpdateBalance(h.Ctx, &input, &output); err != nil {
		h.Log.Warn(err)
		res = apierror.HandleError(err)
		return
	}

	res = apires.With(&output, http.StatusOK)
	return
}

func (h *SwapEventsHandler) Mint() (res interface{}) {
	var mintinput cores.SwapMintInput
	if err := h.Params.BindJsonBody(&mintinput); err != nil {
		h.Log.Warn(err)
		res = apierror.HandleError(err)
		return
	}
	if err := validate.Default.Struct(mintinput); err != nil {
		h.Log.Warn(err)
		res = apierror.HandleError(err)
		return
	}
	var input cores.CreateExchangeTxInput
	input.TxId = mintinput.TxId
	input.PairAddress = mintinput.PairAddress
	input.Sender = mintinput.Sender
	input.Amount0 = mintinput.Amount0
	input.Amount1 = mintinput.Amount1
	input.OccuredAt = mintinput.OccuredAt
	input.Type = cores.ExchangeTxTypeAddLiquidity
	input.Status = cores.ExchangeTxStatusCompleted

	var output cores.CreateExchangeTxResult
	if err := InputExchangeTx(h, &input, &output); err != nil {
		h.Log.Warn(err)
		res = apierror.HandleError(err)
		return
	}

	res = apires.With(&output, http.StatusOK)
	return
}

func (h *SwapEventsHandler) Burn() (res interface{}) {
	var burninput cores.SwapBurnInput
	if err := h.Params.BindJsonBody(&burninput); err != nil {
		h.Log.Warn(err)
		res = apierror.HandleError(err)
		return
	}
	if err := validate.Default.Struct(burninput); err != nil {
		h.Log.Warn(err)
		res = apierror.HandleError(err)
		return
	}
	var input cores.CreateExchangeTxInput
	input.TxId = burninput.TxId
	input.PairAddress = burninput.PairAddress
	input.Sender = burninput.Sender
	input.Amount0 = burninput.Amount0
	input.Amount1 = burninput.Amount1
	input.To = burninput.To
	input.OccuredAt = burninput.OccuredAt
	input.Type = cores.ExchangeTxTypeRemoveLiquidity
	input.Status = cores.ExchangeTxStatusCompleted

	var output cores.CreateExchangeTxResult
	if err := InputExchangeTx(h, &input, &output); err != nil {
		h.Log.Warn(err)
		res = apierror.HandleError(err)
		return
	}

	res = apires.With(&output, http.StatusOK)
	return
}

func (h *SwapEventsHandler) Swap() (res interface{}) {
	var swapinput cores.SwapSwapInput
	if err := h.Params.BindJsonBody(&swapinput); err != nil {
		h.Log.Warn(err)
		res = apierror.HandleError(err)
		return
	}
	if err := validate.Default.Struct(swapinput); err != nil {
		h.Log.Warn(err)
		res = apierror.HandleError(err)
		return
	}
	var input cores.CreateExchangeTxInput
	input.TxId = swapinput.TxId
	input.PairAddress = swapinput.PairAddress
	input.Sender = swapinput.Sender
	if swapinput.Amount0In != "0" {
		input.Amount0 = swapinput.Amount0In
		input.Amount1 = swapinput.Amount1Out
		input.Type = cores.ExchangeTxTypeSwap1for2
	} else if swapinput.Amount1In != "0" {
		input.Amount0 = swapinput.Amount0Out
		input.Amount1 = swapinput.Amount1In
		input.Type = cores.ExchangeTxTypeSwap2for1
	}
	input.To = swapinput.To
	input.OccuredAt = swapinput.OccuredAt
	input.Status = cores.ExchangeTxStatusCompleted

	var output cores.CreateExchangeTxResult
	if err := InputExchangeTx(h, &input, &output); err != nil {
		h.Log.Warn(err)
		res = apierror.HandleError(err)
		return
	}

	res = apires.With(&output, http.StatusOK)
	return
}

func InputExchangeTx(h *SwapEventsHandler, input *cores.CreateExchangeTxInput, output *cores.CreateExchangeTxResult) (err error) {
	var exchangeinput cores.GetExchangeInput
	exchangeinput.PairAddress = input.PairAddress
	var exchangeresult cores.ExchangeResult
	if err = exchangemodels.Exchanges.GetExchange(h.Ctx, &exchangeinput, &exchangeresult); err != nil {
		return
	}
	input.ExchangeId = exchangeresult.Id

	switch input.Type {
	case cores.ExchangeTxTypeAddLiquidity:
		input.Name = "Add " + exchangeresult.TokenSymbol1 + " and " + exchangeresult.TokenSymbol2
	case cores.ExchangeTxTypeRemoveLiquidity:
		input.Name = "Remove " + exchangeresult.TokenSymbol1 + " and " + exchangeresult.TokenSymbol2
	case cores.ExchangeTxTypeSwap1for2:
		input.Name = "Swap " + exchangeresult.TokenSymbol1 + " for " + exchangeresult.TokenSymbol2
	case cores.ExchangeTxTypeSwap2for1:
		input.Name = "Swap " + exchangeresult.TokenSymbol2 + " for " + exchangeresult.TokenSymbol1
	}

	input.TotalValue = 0
	amount1, _ := decimal.NewFromString(input.Amount0)
	divider1Str := strconv.Itoa(exchangeresult.TokenDivider1)
	divider1, _ := decimal.NewFromString(divider1Str)
	input.TokenAmount1, _ = amount1.Div(divider1).Float64()
	amount2, _ := decimal.NewFromString(input.Amount1)
	divider2Str := strconv.Itoa(exchangeresult.TokenDivider2)
	divider2, _ := decimal.NewFromString(divider2Str)
	input.TokenAmount2, _ = amount2.Div(divider2).Float64()
	input.Fee = 0
	input.PricePerToken1 = input.TokenAmount2 / input.TokenAmount1
	input.PricePerToken2 = input.TokenAmount1 / input.TokenAmount2

	var exchangetxinput cores.GetExchangeTxInput
	exchangetxinput.TxId = input.TxId
	var exchangetxoutput cores.ExchangeTxResult
	if err = exchangemodels.Exchanges.GetExchangeTx(h.Ctx, &exchangetxinput, &exchangetxoutput); err == nil {
		if exchangetxoutput.Status == cores.ExchangeTxStatusCompleted {
			output.Id = exchangetxoutput.Id
			output.Status = exchangetxoutput.Status
		} else {
			err = exchangemodels.Exchanges.UpdateExchangeTx(h.Ctx, input, output)
		}
	} else {
		err = exchangemodels.Exchanges.CreateExchangeTx(h.Ctx, input, output)
	}
	return
}

func power(x int, n int) int {
	if n == 0 {
		return 1
	} else {
		return x * power(x, n-1)
	}
}
