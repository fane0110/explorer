package model2

import (
	"github.com/globalsign/mgo/bson"
	"github.com/iost-official/explorer/backend/model2/db"
	"log"
)

/// this struct is used as json to return
type TxnDetail struct {
	Hash        string  `json:"txHash"`
	BlockNumber int64   `json:"blockHeight"`
	From        string  `json:"from"`
	To          string  `json:"to"`
	Amount      float64 `json:"amount"`
	GasLimit    int64   `json:"gasLimit"`
	GasPrice    int64   `json:"price"`
	Age         string  `json:"age"`
	UTCTime     string  `json:"utcTime"`
	Code        string  `json:"code"`
}

func GetDetailTxn(txHash string) (TxnDetail, error) {
	txnC, err := db.GetCollection(db.CollectionFlatTx)

	if err != nil {
		log.Println("failed To open collection collectionTxs")
		return TxnDetail{}, err
	}

	var tx db.FlatTx

	err = txnC.Find(bson.M{"hash": txHash}).One(&tx)

	if err != nil {
		log.Println("transaction not found")
		return TxnDetail{}, err
	}

	txnOut := ConvertFlatTx2TxnDetail(&tx)

	return txnOut, nil
}

/// convert FlatTx to TxnDetail
func ConvertFlatTx2TxnDetail (tx *db.FlatTx) TxnDetail {
	txnOut := TxnDetail{
		Hash:        tx.Hash,
		BlockNumber: tx.BlockNumber,
		From:        tx.From,
		To:          tx.To,
		Amount:      tx.Amount,
		GasLimit:    tx.GasLimit,
		GasPrice:    tx.GasPrice,
	}

	if tx.ActionName == "SetCode" {
		txnOut.Code = tx.Action.Data
	}

	txnOut.Age = modifyIntToTimeStr(tx.Time / (1000 * 1000 * 1000))
	txnOut.UTCTime = formatUTCTime(tx.Time / (1000 * 1000 * 1000))

	return txnOut
}

/// get a list of transactions for a specific page using account and block
func GetFlatTxnSlicePage(page, eachPageNum, block int64, address string) ([]*TxnDetail, error) {
	lastPageNum, err := db.GetFlatTxTotalPageCnt(eachPageNum, address, block)

	if lastPageNum == 0 {
		return []*TxnDetail{}, nil
	}

	if err != nil {
		return nil, err
	}

	if page > lastPageNum {
		page = lastPageNum
	}

	start := int((page - 1) * eachPageNum)
	txnsFlat, err := db.GetFlatTxnSlice(start, int(eachPageNum), int(block), address)

	if err != nil {
		return nil, err
	}

	var txnDetailList []*TxnDetail

	for _, v := range txnsFlat {
		td := ConvertFlatTx2TxnDetail(v)
		txnDetailList = append(txnDetailList, &td)
	}

	return txnDetailList, nil
}