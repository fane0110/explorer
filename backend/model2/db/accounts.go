package db

import (
	"github.com/globalsign/mgo/bson"
	"log"
	"time"
)

type Account struct {
	Address string  `json:"address"`
	Balance float64 `json:"balance"`
	Percent float64 `json:"percent"`
	TxCount int     `json:"txCount"`
}

type ApplyTestIOST struct {
	Address   string  `json:"address"`
	Amount    float64 `json:"amount"`
	Email     string  `json:"email"`
	Mobile    string  `json:"mobile"`
	ApplyTime int64   `json:"apply_time"`
}

type AddressNonce struct {
	Address string `json:"address"`
	Nonce   int64  `json:"nonce"`
}

func GetAccounts(start, limit int) ([]*Account, error) {
	accountC, err := GetCollection("accounts")
	if err != nil {
		return nil, err
	}
	query := bson.M{
		"balance": bson.M{"$ne": 0},
	}
	var accountList []*Account
	err = accountC.Find(query).Sort("-balance").Skip(start).Limit(limit).All(&accountList)
	if err != nil {
		return nil, err
	}

	return accountList, nil
}

func GetAccountByAddress(address string) (*Account, error) {
	accountC, err := GetCollection("accounts")
	if err != nil {
		return nil, err
	}

	query := bson.M{
		"address": address,
	}
	var account *Account
	err = accountC.Find(query).One(&account)

	if err != nil {
		return nil, err
	}

	return account, nil
}

func GetAccountsTotalLen() (int, error) {
	accountC, err := GetCollection("accounts")
	if err != nil {
		return 0, err
	}
	query := bson.M{
		"balance": bson.M{"$ne": 0},
	}
	return accountC.Find(query).Count()
}

func GetAccountLastPage(eachPage int64) (int64, error) {
	accountC, err := GetCollection("accounts")
	if err != nil {
		log.Println("GetAccounts get collection error:", err)
		return 0, err
	}

	query := bson.M{
		"balance": bson.M{"$ne": 0},
	}
	totalLen, _ := accountC.Find(query).Count()
	totalLenInt64 := int64(totalLen)

	var pageLast int64
	if totalLenInt64%eachPage == 0 {
		pageLast = totalLenInt64 / eachPage
	} else {
		pageLast = totalLenInt64/eachPage + 1
	}

	if pageLast == 0 {
		pageLast = 1
	}

	return pageLast, nil
}

func GetAccountTxnLastPage(address string, eachPage int64) (int64, error) {
	txnLen, err := GetFlatTxnLenByAccount(address)
	if err != nil {
		return 0, err
	}
	txnLenInt64 := int64(txnLen)

	var pageLast int64
	if txnLenInt64%eachPage == 0 {
		pageLast = txnLenInt64 / eachPage
	} else {
		pageLast = txnLenInt64/eachPage + 1
	}

	if pageLast == 0 {
		pageLast = 1
	}

	return pageLast, nil
}

func SaveApplyTestIOST(at *ApplyTestIOST) error {
	applyC, err := GetCollection("applyTestIOST")
	if err != nil {
		log.Println("SaveApplyTestIost get collection error:", err)
		return err
	}

	return applyC.Insert(at)
}

func GetApplyNumTodayByMobile(mobile string) (int, error) {
	applyC, err := GetCollection("applyTestIOST")
	if err != nil {
		log.Println("SaveApplyTestIost get collection error:", err)
		return 0, err
	}

	t := time.Now()
	dayBegin := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location()).Unix()
	dayEnd := dayBegin + 24*3600

	query := bson.M{
		"mobile": mobile,
		"applytime": bson.M{
			"$gte": dayBegin,
			"$lt":  dayEnd,
		},
	}
	return applyC.Find(query).Count()
}

func IncAddressNonce(address string) error {
	anc, err := GetCollection("addressNonce")
	if err != nil {
		log.Println("IncAddressNonce get collection error:", err)
		return err
	}

	query := bson.M{"address": address}
	inc := bson.M{"$inc": bson.M{"nonce": 1}}

	_, err = anc.Upsert(query, inc)
	return err
}

func GetAddressNonce(address string) (int64, error) {
	anc, err := GetCollection("addressNonce")
	if err != nil {
		log.Println("GetAddressNonce get collection error:", err)
		return 0, err
	}

	query := bson.M{"address": address}

	var addressNonce *AddressNonce
	err = anc.Find(query).One(&addressNonce)
	if err != nil {
		return 0, err
	}

	return addressNonce.Nonce, nil
}

func GetFlatTxnLenByAccount(account string) (int, error) {
	txnDC, err := GetCollection(CollectionFlatTx)
	if err != nil {
		log.Println("GetFlatTxnLenByAccount CollectionFlatTx collection error:", err)
		return 0, err
	}

	fromLen, err := txnDC.Find(bson.M{"from": account}).Count()
	if err != nil {
		log.Println("GetFlatTxnLenByAccount get from len error:", err)

		return 0, err
	}
	toLen, err := txnDC.Find(bson.M{"to": account}).Count()
	if err != nil {
		log.Println("GetFlatTxnLenByAccount get to len error:", err)

		return 0, err
	}
	return fromLen + toLen, nil
}

func GetAccountTxCount(address string) (int, error) {
	ftxCol, err := GetCollection(CollectionFlatTx)
	if err != nil {
		return 0, err
	}
	num, err := ftxCol.Find(bson.M{"$or": []bson.M{
		bson.M{"from": address},
		bson.M{"to": address},
	}}).Count()
	return num, err
}

func GetTxnListByAccount(account string, start, limit int) ([]*FlatTx, error) {
	txnDC, err := GetCollection(CollectionFlatTx)
	if err != nil {
		return nil, err
	}
	query := bson.M{
		"$or": []bson.M{
			bson.M{"from": account},
			bson.M{"to": account},
		},
	}
	var txnList []*FlatTx
	err = txnDC.Find(query).Sort("-blockNumber").Skip(start).Limit(limit).All(&txnList)
	if err != nil {
		return nil, err
	}
	return txnList, nil
}
