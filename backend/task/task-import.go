package main

import (
	"github.com/iost-official/explorer/backend/config"
	"os"
	"sync"
	//"github.com/iost-official/explorer/backend/task/cron"
	"github.com/iost-official/explorer/backend/model/db"
	//"time"
	"fmt"
	"io/ioutil"
	"log"
	"os/exec"
	"strconv"
	"strings"
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
	mergecollections := []interface{}{
			//[]interface{}{db.CollectionBlocks, []string{"number"}},
			//[]interface{}{db.CollectionTxs, []string{"txhash"}},
			[]interface{}{db.CollectionAccount, []string{"name"}},
			//[]interface{}{db.CollectionAccountTx, []string{"txhash", "name"}},
			[]interface{}{db.CollectionAccountPubkey, []string{"name", "pubkey"}},
			[]interface{}{db.CollectionContract, []string{"id"}},
			//[]interface{}{db.CollectionContractTx, []string{"id", "time", "txhash"}},
			//[]interface{}{db.CollectionVoteTx, []string{"action", "blockNumber"}},
			} 
	
	
	var insertcollections []string = []string{
			db.CollectionBlocks                    , 
			db.CollectionTxs                       , 
			//db.CollectionAccount                   , 
			db.CollectionAccountTx                 , 
			//db.CollectionAccountPubkey             , 
			//db.CollectionContract                  , 
			db.CollectionContractTx                , 
			db.CollectionVoteTx                    ,}

	dname := "explorer"

	fromblock, _ := strconv.ParseInt(os.Args[1], 10, 64)
	toblock, _ := strconv.ParseInt(os.Args[2], 10, 64)

	inner_from_block := fromblock
	inner_to_block := inner_from_block + interval - 1

	for inner_from_block < toblock {
		log.Print(inner_from_block)
		from_to_name := strconv.FormatInt(inner_from_block, 10) + "_" + strconv.FormatInt(inner_to_block, 10)
		exec.Command("tar","-zxvf",
		 "src/"+from_to_name+".tar.gz",
		 ).Run()

		for _, value := range insertcollections {
			ws2.Add(1)
			go func(Collection string) {
				collectionname := Collection
				
				cmd := exec.Command("mongoimport",
					"-d="+dname, //explorer
					"-c="+collectionname,
					"--type=json",
					"--file="+"src/"+from_to_name+"/"+collectionname+".json")
				stdErrorPipe, err := cmd.StderrPipe()
				if err != nil {
					log.Fatal(err)
				}

				if err := cmd.Start(); err != nil {
					log.Fatal(err)
				}

				slurp, _ := ioutil.ReadAll(stdErrorPipe)
				fmt.Printf("stderr: %s\n", collectionname)
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

		for _, value := range mergecollections {
			ws2.Add(1)
			go func(eCollection []interface{}) {
				collectionname := eCollection[0].(string)
				fields := eCollection[1].([]string)
				var m2 string

				m2 = strings.Join(fields, ",")
				log.Print(m2)
				cmd := exec.Command("mongoimport",
					"-d="+dname, //explorer
					"-c="+collectionname,
					"--type=json",
					//"--mode=merge",
					//"--upsertFields="+m2,
					"--file="+"src/"+from_to_name+"/"+collectionname+".json")
				stdErrorPipe, err := cmd.StderrPipe()
				if err != nil {
					log.Fatal(err)
				}

				if err := cmd.Start(); err != nil {
					log.Fatal(err)
				}

				slurp, _ := ioutil.ReadAll(stdErrorPipe)
				fmt.Printf("stderr: %s\n", collectionname)
				fmt.Printf("stderr: %s\n", slurp)

				if err := cmd.Wait(); err != nil {
					log.Fatal(err)
				}
				if err != nil {
					log.Println("command err:", err)
				}
				ws2.Done()
			}(value.([]interface{}))
		
			
		}

		ws2.Wait()

		exec.Command("rm","-rf",
		"src/"+from_to_name,
		 ).Run()

		inner_from_block += interval
		inner_to_block += interval
	}

	// time.Sleep(10 * time.Second)
	// start tasks

}
