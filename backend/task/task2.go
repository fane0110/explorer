package main

import (
	"sync"

	"github.com/iost-official/explorer/backend/config"
	"github.com/iost-official/explorer/backend/task/cron"
	
)

var ws = new(sync.WaitGroup)

func main() {
	config.ReadConfig("")

	// start tasks
	ws.Add(1)
	// download block
	go cron.LimitUpdateBlocks(ws)
	ws.Wait()
}
