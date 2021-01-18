package exchanges

import (
	"cos-backend-com/src/common/validate"
	"cos-backend-com/src/cores/routers"
	"cos-backend-com/src/libs/apierror"
	"cos-backend-com/src/libs/models/exchangemodels"
	"cos-backend-com/src/libs/sdk/cores"
	"fmt"
	"github.com/shopspring/decimal"
	"github.com/wujiu2020/strip/utils/apires"
	"math"
	"net/http"
	"strconv"
)

type SwapEventsHandler struct {
	routers.Base
}

func (h *SwapEventsHandler) CreatePair() (res interface{}) {
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
		if exchangeresult.Status != cores.ExchangeStatusCompleted {
			if err := exchangemodels.Exchanges.UpdateExchange(h.Ctx, &input, &output); err != nil {
				h.Log.Warn(err)
				res = apierror.HandleError(err)
				return
			}
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

func (h *SwapEventsHandler) Mint() (res interface{}) {
	var mintinput cores.CreateSwapMintInput
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

	var exchangeinput cores.GetExchangeInput
	exchangeinput.StartupId = mintinput.StartupId
	var exchangeresult cores.ExchangeResult
	if err := exchangemodels.Exchanges.GetExchange(h.Ctx, &exchangeinput, &exchangeresult); err != nil {
		h.Log.Warn(err)
		res = apierror.HandleError(err)
		return
	}

	var input cores.CreateExchangeTxInput
	input.TxId = mintinput.TxId
	input.ExchangeId = exchangeresult.Id
	input.Sender = mintinput.Sender
	input.Type = cores.ExchangeTxTypeAddLiquidity
	input.Name = "Add " + exchangeresult.TokenSymbol1 + " and " + exchangeresult.TokenSymbol2
	input.TotalValue = 0
	input.TokenAmount1 = 10124504517.37366977 / math.Pow10(8)
	//fmt.Printf("pi=%s\n", strconv.FormatUint(mintinput.Amount0,10))
	//fmt.Printf("%.10f\n",mintinput.Amount0)
	//fmt.Printf("input.TokenAmount1=%.20f\n",input.TokenAmount1)
	//fmt.Println(mintinput.Amount0)
	//fmt.Println(float64(mintinput.Amount0))
	input.TokenAmount2 = float64(mintinput.Amount1) / float64(exchangeresult.TokenDivider2)
	fmt.Printf("input.TokenAmount2=%.10f\n", input.TokenAmount2)
	input.Fee = 0

	amount1Str := strconv.FormatUint(mintinput.Amount0, 10)
	amount1, _ := decimal.NewFromString(amount1Str)
	amount2Str := strconv.FormatUint(mintinput.Amount1, 10)
	amount2, _ := decimal.NewFromString(amount2Str)
	divider1Str := strconv.Itoa(exchangeresult.TokenDivider1)
	divider1, _ := decimal.NewFromString(divider1Str)
	a1 := amount1.Div(divider1)
	input.TokenAmount1, _ = a1.Float64()
	fmt.Println("a1=", a1)
	fmt.Println("input.TokenAmount1=", input.TokenAmount1)
	fmt.Println("mint.amount0=", mintinput.Amount0)
	fmt.Println("amount1=", amount1)
	fmt.Println("amount2=", amount2)

	//input.PricePerToken1 = amount2 / amount1
	input.PricePerToken2 = input.TokenAmount1 / input.TokenAmount2
	input.Status = cores.ExchangeTxStatusCompleted

	var output cores.CreateExchangeTxResult
	var exchangetxinput cores.GetExchangeTxInput
	exchangetxinput.TxId = mintinput.TxId
	var exchangetxoutput cores.ExchangeTxResult
	if err := exchangemodels.Exchanges.GetExchangeTx(h.Ctx, &exchangetxinput, &exchangetxoutput); err == nil {
		fmt.Println("tx existed")
		if exchangetxoutput.Status != cores.ExchangeTxStatusCompleted {
			fmt.Println("tx status not completed update")
			if err := exchangemodels.Exchanges.UpdateExchangeTx(h.Ctx, &input, &output); err != nil {
				h.Log.Warn(err)
				res = apierror.HandleError(err)
				return
			}
		}
	} else {
		fmt.Println("tx new create")
		if err := exchangemodels.Exchanges.CreateExchangeTx(h.Ctx, &input, &output); err != nil {
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
