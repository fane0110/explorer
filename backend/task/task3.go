package main

import (
	"fmt"
	"sync"
	"os"
	"strconv"
	"github.com/iost-official/explorer/backend/config"
	"github.com/iost-official/explorer/backend/task/cron"
	
)

var ws1 = new(sync.WaitGroup)

func main() {
	config.ReadConfig("")
	s1 := os.Args[1]
	getnum, _ := strconv.ParseInt(s1, 10, 64)
	
	fmt.Println(getnum)           // 98
	fmt.Printf("%T\n", getnum)  
	// start tasks
	ws1.Add(1)
	// download block
	go cron.GetBlock(ws1,getnum)
	ws1.Wait()
	
}
