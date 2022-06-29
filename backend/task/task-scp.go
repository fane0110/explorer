package main

import (
	"github.com/iost-official/explorer/backend/config"
	"os"
	"sync"
	//"github.com/iost-official/explorer/backend/task/cron"
	//"github.com/iost-official/explorer/backend/model/db"
	//"time"
	"fmt"
	"io/ioutil"
	"log"
	"os/exec"
	"strconv"
	"io"
	//"strings"
	//"github.com/globalsign/mgo"
)

var ws2 = new(sync.WaitGroup)

func main() {
	config.ReadConfig("")
	//avg := sigar.LoadAverage{}

	const interval = int64(1000)
	/* collections := []interface{}{
		[]interface{}{db.CollectionBlocks, []string{"number"}},
		[]interface{}{db.CollectionTxs, []string{"txhash"}},
		[]interface{}{db.CollectionAccount, []string{"name"}},
		[]interface{}{db.CollectionAccountTx, []string{"txhash", "name"}},
		[]interface{}{db.CollectionAccountPubkey, []string{"name", "pubkey"}},
		[]interface{}{db.CollectionContract, []string{"id"}},
		[]interface{}{db.CollectionContractTx, []string{"id", "time", "txhash"}},
		[]interface{}{db.CollectionVoteTx, []string{"action", "blockNumber"}}} */


	fromblock, _ := strconv.ParseInt(os.Args[1], 10, 64)
	toblock, _ := strconv.ParseInt(os.Args[2], 10, 64)

	inner_from_block := fromblock
	inner_to_block := inner_from_block + interval - 1

	for inner_from_block < toblock {
		log.Print(inner_from_block)
		from_to_name := strconv.FormatInt(inner_from_block, 10) + "_" + strconv.FormatInt(inner_to_block, 10)
		
		

		ws2.Add(1)
		go func(from_to_name string) {
			
			
			cmd := exec.Command("scp",
				from_to_name+".tar.gz",
				"es_admin@"+"45.76.52.95"+":/home/es_admin/.go/src/github.com/iost-official/explorer/backend/task",
				)
			stdErrorPipe, err := cmd.StderrPipe()
			stdin, _ := cmd.StdinPipe()
			io.WriteString(stdin, "yes")
			stdin.Close()
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
		}(from_to_name)
		



		ws2.Wait()

		inner_from_block += interval
		inner_to_block += interval
	}

	// time.Sleep(10 * time.Second)
	// start tasks

}
