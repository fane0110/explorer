package main

import (
	
	"sync"
	"github.com/iost-official/explorer/backend/config"
	"github.com/iost-official/explorer/backend/task/cron"
	"github.com/iost-official/explorer/backend/model/db"
	"time"
	"fmt"
	"log"
	
)


var ws2 = new(sync.WaitGroup)

func main() {
	config.ReadConfig("")
	
	var topHeightInMongo int64
	maxSessions :=150
	ticker := time.NewTicker(time.Second)
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
	
	for j := 1;;j++{
		
		
		
		for i := 1; i <= maxSessions; i++ {
			ws2.Add(1)
			getnum :=topHeightInMongo + int64(i)

			// download block
			go cron.GetBlock(ws2,getnum)
			/*go func(i int64) {
				fmt.Println(i)
				ws2.Done()
			
			}(getnum) */
		}
		
		//fmt.Println(j)
		ws2.Wait()
		
		topHeightInMongo += int64(maxSessions)
	}
	// time.Sleep(10 * time.Second)
	// start tasks

	
	
}

func Printhoge(ws *sync.WaitGroup,getnum int64) {
	defer ws.Done()
	fmt.Print(getnum,"\n")
}

