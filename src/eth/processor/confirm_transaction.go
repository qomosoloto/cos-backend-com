package processor

import (
	"context"
	"cos-backend-com/src/eth/proto"
	ethSdk "cos-backend-com/src/libs/sdk/eth"

	"github.com/ethereum/go-ethereum/common"
	"github.com/qiniu/x/log"
)

type Confirmer struct {
	TransactionInput  <-chan *proto.TransactionsOutput
	TransactionOutput chan<- *proto.TransactionsOutput
}

func (c *Confirmer) Process() {
	for transactionInput := range c.TransactionInput {
		txHash := common.HexToHash(transactionInput.TxId)
		receipt, err := EthClient.TransactionReceipt(context.Background(), txHash)
		if err != nil {
			if transactionInput.RetryTime < 5 {
				transactionInput.State = ethSdk.TransactionStateWaitConfirm
			} else {
				transactionInput.State = ethSdk.TransactionStateFailed
			}
			log.Warn(err, txHash)
		} else {
			transactionInput.BlockAddr = receipt.BlockHash.Hex()
			if receipt.Status == 1 {
				transactionInput.State = ethSdk.TransactionStateSuccess
			} else {
				transactionInput.State = ethSdk.TransactionStateFailed
				log.Warn("transaction failed", txHash)
			}
		}
		c.TransactionOutput <- transactionInput
	}
}
