package main

import (
	"os"
	"sync"
	"github.com/iost-official/explorer/backend/config"
	"github.com/iost-official/explorer/backend/task/cron"
	"github.com/iost-official/explorer/backend/model/db"
	"time"
	"fmt"
	"log"
	"strconv"
	"os/exec"
	"io/ioutil"
	"github.com/globalsign/mgo"

)




var ws2 = new(sync.WaitGroup)

func main() {
	config.ReadConfig("")
	//avg := sigar.LoadAverage{}
	var topHeightInMongo int64
	var retryTime int
	var d *mgo.Database
	var err error

	var collections []string = []string{
		db.CollectionBlocks                    , 
		db.CollectionBP                        , 
		db.CollectionTxs                       , 
		db.CollectionFlatTx                    , 
		db.CollectionAccount                   , 
		db.CollectionAccountTx                 , 
		db.CollectionAccountPubkey             , 
		db.CollectionContract                  , 
		db.CollectionContractTx                , 
		db.CollectionTaskCursor                , 
		db.CollectionBlockPay                  , 
		db.CollectionApplyIOST                 , 
		db.CollectionVoteTx                    ,
		db.CollectionProducerAward             , 
		db.CollectionUserAward                 , 
		db.CollectionProducerContributionAward , 
		db.CollectionUserContributionAward     , 
		db.CollectionFailedAward               , 
		db.CollectionAwardInfo                 , 
		db.CollectionProducerLevelInfo          }

	maxSessions :=20
	ticker := time.NewTicker(time.Second)
	
	fromblock, _ := strconv.ParseInt(os.Args[1], 10, 64)
	toblock, _ := strconv.ParseInt(os.Args[2], 10, 64)


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

	topHeightInMongo = fromblock - 1


	
	
	for {
		log.Print("delete取得 ")
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
	d.DropDatabase()//いらないデータを一度消す
	
	
	for j := 1;;j++{
		
		for i := 1; i <= maxSessions; i++ {
			ws2.Add(1)
			getnum :=topHeightInMongo + int64(i)
			log.Print(getnum)
			if getnum > toblock{
				log.Print("stop")
				ws2.Done()
				continue
			}
			// download block
			go cron.GetBlock(ws2,getnum)
			/* go func(i int64) {
				fmt.Println(i)
				
				ws2.Done()
			
			}(getnum)  */
		}
		
		//fmt.Println(j)
		ws2.Wait()

		topHeightInMongo += int64(maxSessions)
		if topHeightInMongo >= toblock{
			break
		}

	}
	
	
	from_to_name:=strconv.FormatInt(fromblock,10)+"_"+strconv.FormatInt(toblock,10)
	exec.Command("mkdir",from_to_name).Run()

	for _, value := range collections {
		ws2.Add(1)
		go func(collectionname string){
			cmd:=exec.Command("mongoexport",
			"-d=explorer",//explorer
			"-c="+collectionname,
			"--type=json",
			"--out="+from_to_name+"/"+collectionname+"_"+from_to_name+".json")
			stdErrorPipe, err := cmd.StderrPipe()
			if err != nil {
				log.Fatal(err)
			}

			if err := cmd.Start(); err != nil {
				log.Fatal(err)
			}

			slurp, _ := ioutil.ReadAll(stdErrorPipe)
			fmt.Printf("stderr: %s\n", slurp)

			if err := cmd.Wait(); err != nil {
				log.Fatal(err)
			}
			if err != nil {
				log.Println("command err:", err)
			}
			ws2.Done()
		}(value)
	}
	ws2.Wait()
	// time.Sleep(10 * time.Second)
	// start tasks
	
	
	
}

func Printhoge(ws *sync.WaitGroup,getnum int64) {
	defer ws.Done()
	fmt.Print(getnum,"\n")
}

