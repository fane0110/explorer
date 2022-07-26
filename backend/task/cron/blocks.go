package cron

import (
	"log"
	"sync"
	"time"
	"github.com/globalsign/mgo"
	"github.com/iost-official/explorer/backend/model/blockchain"
	"github.com/iost-official/explorer/backend/model/blockchain/rpcpb"
	"github.com/iost-official/explorer/backend/model/db"
)

func UpdateBlocks(ws *sync.WaitGroup) {
	defer ws.Done()

	blockChannel := make(chan *rpcpb.Block, 10)
	go insertBlock(blockChannel)

	ticker := time.NewTicker(time.Second)

	var topHeightInMongo int64
	for range ticker.C {
		topBlkInMongo, err := db.GetTopBlock()
		if err != nil {
			log.Println("updateBlock get topBlk in mongo error:", err)
			if err.Error() != "not found" {
				continue
			} else {
				topHeightInMongo = 0
				break
			}
		}
		topHeightInMongo = topBlkInMongo.Number + 1
		log.Println("Got Top Block From Mongo With Number: ", topHeightInMongo)
		break
	}

	for {
		blockRspn, err := blockchain.GetBlockByNum(topHeightInMongo, true)
		if err != nil {
			log.Println("Download block", topHeightInMongo, "error:", err)
			time.Sleep(time.Second)
			continue
		}
		if blockRspn.Status == rpcpb.BlockResponse_PENDING {
			log.Println("Download block", topHeightInMongo, "Pending")
			time.Sleep(time.Second)
			continue
		}
		blockChannel <- blockRspn.Block
		topHeightInMongo++
		log.Println("Download block", topHeightInMongo, " Succ!")


	}

}

func LimitUpdateBlocks(ws *sync.WaitGroup) {
	limit:= int64(210461000)
	defer ws.Done()

	blockChannel := make(chan *rpcpb.Block, 10)
	go insertBlock(blockChannel)

	ticker := time.NewTicker(time.Second)

	var topHeightInMongo int64
	for range ticker.C {
		topBlkInMongo, err := db.GetTopBlock()
		if err != nil {
			log.Println("updateBlock get topBlk in mongo error:", err)
			if err.Error() != "not found" {
				continue
			} else {
				topHeightInMongo = 0
				break
			}
		}
		topHeightInMongo = topBlkInMongo.Number + 1
		log.Println("Got Top Block From Mongo With Number: ", topHeightInMongo)
		break
	}

	for {
		if topHeightInMongo == limit{
			time.Sleep(time.Second * 5)
			break
		}
		blockRspn, err := blockchain.GetBlockByNum(topHeightInMongo, true)
		if err != nil {
			log.Println("Download block", topHeightInMongo, "error:", err)
			time.Sleep(time.Second)
			continue
		}
		if blockRspn.Status == rpcpb.BlockResponse_PENDING {
			log.Println("Download block", topHeightInMongo, "Pending")
			time.Sleep(time.Second)
			continue
		}
		blockChannel <- blockRspn.Block
		topHeightInMongo++
		log.Println("Download block", topHeightInMongo, " Succ!")


	}

}



func GetBlock(ws *sync.WaitGroup,getnum int64) {
	defer ws.Done()
	var d *mgo.Database
	var err error
	var retryTime int

	//blockChannel := make(chan *rpcpb.Block, 10)
	//go insertBlock(blockChannel)

	var topHeightInMongo int64

	var blockRspn *rpcpb.BlockResponse

	topHeightInMongo = getnum
	for{
		blockRspn, err = blockchain.GetBlockByNum(topHeightInMongo, true)
		if err != nil {
			log.Println("Download block", topHeightInMongo, "error:", err)
			time.Sleep(time.Second )
			continue
	
		}
		break
	}

	if blockRspn.Status == rpcpb.BlockResponse_PENDING {
		log.Println("Download block", topHeightInMongo, "Pending")
		time.Sleep(time.Second)
	}
	//blockChannel <- blockRspn.Block

	


	for {
		d, err = db.GetDb()
		if err != nil {
			log.Println("fail to get db collection ", err)
			time.Sleep(time.Second)
			retryTime++
			if retryTime > 10 {
				log.Fatalln("fail to get db collection, retry time exceeds")
			}
			continue
		}
		break
	}
	
	
	hogeinsertBlock(d,blockRspn.Block)
	log.Println("Download block", topHeightInMongo, " Succ!")


	
	
	

}




func HogeGetBlock(sdb *mgo.Database,ws *sync.WaitGroup,getnum int64) {
	defer ws.Done()

	var err error
	var blockRspn *rpcpb.BlockResponse

	//blockChannel := make(chan *rpcpb.Block, 10)
	//go insertBlock(blockChannel)

	var topHeightInMongo int64

	topHeightInMongo = getnum
	for{
		blockRspn, err = blockchain.GetBlockByNum(topHeightInMongo, true)
		if err != nil {
			log.Println("Download block", topHeightInMongo, "error:", err)
			time.Sleep(time.Second )
			continue
	
		}
		break
	}

	if blockRspn.Status == rpcpb.BlockResponse_PENDING {
		log.Println("Download block", topHeightInMongo, "Pending")
		time.Sleep(time.Second)
	}
	//blockChannel <- blockRspn.Block

	


	
	
	
	hogeinsertBlock(sdb,blockRspn.Block)
	log.Println("Download block", topHeightInMongo, " Succ!")

	
	
	
	

}


func insertBlock(blockChannel chan *rpcpb.Block) {
	collection := db.GetCollection(db.CollectionBlocks)

	for {
		select {
		case b := <-blockChannel:
			txs := b.Transactions

			wg := new(sync.WaitGroup)
			wg.Add(2)
			go func() {
				db.ProcessTxs(txs, b.Number)
				wg.Done()
			}()
			go func() {
				db.ProcessTxsForAccount(txs, b.Time, b.Number)
				wg.Done()
			}()
			wg.Wait()

			b.Transactions = make([]*rpcpb.Transaction, 0)
			
			err := collection.Insert(*b)

			if err != nil {
				log.Println("updateBlock insert mongo error:", err)
			}
			
		default:

		}
	}
}

func hogeinsertBlock(d *mgo.Database,blockChannel  *rpcpb.Block) {



	collection := db.HogeGetCollection(d,db.CollectionBlocks)
	
	b := blockChannel
	txs := b.Transactions

	wg := new(sync.WaitGroup)
	wg.Add(2)
	go func() {
		db.HogeProcessTxs(d,txs, b.Number)
		wg.Done()
	}()
	go func() {
		db.HogeProcessTxsForAccount(d,txs, b.Time, b.Number)
		wg.Done()
	}()
	wg.Wait()

	b.Transactions = make([]*rpcpb.Transaction, 0)
	
	
	wg.Add(1)
	go func(){
		defer wg.Done()
		err2 := collection.Insert(*b)
	
		if err2 != nil {
			log.Println("updateBlock insert mongo error:", err2)
		}
		log.Println("Insert end : block ", b.Number)
		
	}()
	//log.Print("待機開始")
	wg.Wait()
	//log.Print("待機おわり")

	
}