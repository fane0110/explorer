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
	"strings"
	"github.com/globalsign/mgo"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"

)




var ws2 = new(sync.WaitGroup)


func main() {
	config.ReadConfig("")
	//avg := sigar.LoadAverage{}
	var topHeightInMongo int64
	var retryTime int
	var d *mgo.Database
	var err error
	const interval =int64(1000)
	var collections []string = []string{
		db.CollectionBlocks                    , 
		db.CollectionTxs                       , 
		db.CollectionAccount                   , 
		db.CollectionAccountTx                 , 
		db.CollectionAccountPubkey             , 
		db.CollectionContract                  , 
		db.CollectionContractTx                , 
		db.CollectionVoteTx                    ,}

	maxSessions :=100
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
	inner_from_block :=fromblock
	inner_to_block := inner_from_block +interval -1

	// sessionの作成
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		Config: aws.Config{Region: aws.String("ap-southeast-1")},
		Profile:           "default",
		SharedConfigState: session.SharedConfigEnable,
	}))


	
	for;inner_from_block<toblock;{
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
				if getnum > inner_to_block{
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
			if topHeightInMongo >= inner_to_block{
				break
			}
	
		}
		
		
		from_to_name:=strconv.FormatInt(inner_from_block,10)+"_"+strconv.FormatInt(inner_to_block,10)
		exec.Command("mkdir","src/"+from_to_name).Run()
	
		for _, value := range collections {
			ws2.Add(1)
			go func(collectionname string){

				cmdstr := strings.Join([]string{
					"mongoexport",
				"-d=explorer4",//explorer
				"-c="+collectionname,
				"--type=json",
				"|","sed",
				"'/"+strconv.Quote("_id")+":/s/"+strconv.Quote("_id")+":[^,]*,//'",
				">",
				"src/"+from_to_name+"/"+collectionname+".json",
				}, " ") 
				log.Print(cmdstr)
				cmd:=exec.Command("sh", "-c", cmdstr)

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

		ws2.Add(1)
		go func(from_to_name string) {
			
			
			cmd := exec.Command("tar",
				"-zcvf", //explorer
				"src/"+from_to_name+".tar.gz",
				"src/"+from_to_name,
				)
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
		}(from_to_name)
		ws2.Wait()

		exec.Command("rm","-rf",
		"src/"+from_to_name,
		 ).Run()

		//S3へアップロード
		ws2.Add(1)
		go func(from_to_name string) {
				

			// ファイルを開く
			targetFilePath := "src/"+from_to_name+".tar.gz"
			file, err := os.Open(targetFilePath)
			if err != nil {
				log.Fatal(err)
			}
			defer file.Close()

			bucketName := "stg.stir-hosyu"
			objectKey := "iost-explorer2/"+from_to_name+".tar.gz"

			// Uploaderを作成し、ローカルファイルをアップロード
			uploader := s3manager.NewUploader(sess)
			_, err = uploader.Upload(&s3manager.UploadInput{
				Bucket: aws.String(bucketName),
				Key:    aws.String(objectKey),
				Body:   file,
			})
			if err != nil {
				log.Fatal(err)
			}
			log.Println("done") 
			ws2.Done()
		}(from_to_name)

		ws2.Wait()


		inner_from_block += interval
		inner_to_block += interval 
	}
	
	// time.Sleep(10 * time.Second)
	// start tasks
	
	
	
}
func Printhoge(ws *sync.WaitGroup,getnum int64) {
	defer ws.Done()
	fmt.Print(getnum,"\n")
}

